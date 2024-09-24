package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/heroiclabs/nakama-common/runtime"
)

const TICK_RATE = 3 // number of ticks the server runs per second
const W = 30
const H = 20

const NUM_BOTS = 2

type SMatchLabel struct {
	Open  int `json:"open"`
	Snake int `json:"snake"`
}

type SMatchState struct {
	// match lifecycle related
	playing         bool
	label           *SMatchLabel
	joinsInProgress int

	// user maps
	presences map[string]*runtime.Presence

	snakeGame SnakeGame
}

type SMatch struct{}

func newMatch(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule) (m runtime.Match, err error) {
	return &SMatch{}, nil
}

func (m *SMatch) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	state := &SMatchState{
		playing:         false,
		label:           &SMatchLabel{Open: 1, Snake: 1},
		joinsInProgress: 0,

		presences: make(map[string]*runtime.Presence),

		snakeGame: *newSnakeGame(W, H, 0),
	}

	label := ""
	labelBytes, err := json.Marshal(state.label)
	if err == nil {
		label = string(labelBytes)
	}

	for botI := 0; botI < NUM_BOTS; botI++ {
		state.snakeGame.addSnake()
		state.snakeGame.SnakeIds = append(state.snakeGame.SnakeIds, "")
	}

	return state, TICK_RATE, label
}

func (m *SMatch) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	state := state_.(*SMatchState)

	state.joinsInProgress++
	return state, true, ""
}

func (m *SMatch) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, presences []runtime.Presence) interface{} {
	state := state_.(*SMatchState)

	for _, p := range presences {
		state.joinsInProgress--
		id := p.GetUserId()
		state.presences[id] = &p
		state.snakeGame.addSnake()
		state.snakeGame.SnakeIds = append(state.snakeGame.SnakeIds, id)
	}

	state.playing = true

	return state
}

func (m *SMatch) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, presences []runtime.Presence) interface{} {
	state := state_.(*SMatchState)

	for _, p := range presences {
		id := p.GetUserId()
		snI := state.snakeGame.getSnakeIndexFromId(id)
		bcFeedback(dispatcher, fmt.Sprintf("player %s left!", id), nil)
		state.snakeGame.removeSnake(snI)
		delete(state.presences, id)
	}

	bcUpdate(dispatcher, state, nil)

	if len(state.presences) == 0 {
		return nil
	}

	return state
}

func (m *SMatch) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, messages []runtime.MatchData) interface{} {
	state := state_.(*SMatchState)

	if !state.playing {
		return state
	}

	// move snakes
	for snI, snake := range state.snakeGame.Snakes {
		if !state.snakeGame.validateDirection(snake, snake.Direction) {
			state.playing = false
			feedbackContents := fmt.Sprintf("Player %s lost!", state.snakeGame.SnakeIds[snI])
			logger.Error(feedbackContents)
			bcFeedback(dispatcher, feedbackContents, nil)
			break
		}
		state.snakeGame.move(snake)
	}

	bcUpdate(dispatcher, state, nil)

	if !state.playing {
		return state
	}

	// accept changes of direction from bots
	for snI, snake := range state.snakeGame.Snakes {
		if state.snakeGame.SnakeIds[snI] != "" {
			continue
		}
		potMoves := state.snakeGame.getValidDirections(snake)
		if len(potMoves) > 0 {
			snake.Direction = potMoves[rand.Intn(len(potMoves))]
		}
	}

	// accept changes of direction from players
	for _, message := range messages {
		senderUserId := message.GetUserId()
		op := message.GetOpCode()
		data := message.GetData()
		snake := state.snakeGame.Snakes[state.snakeGame.getSnakeIndexFromId(senderUserId)]

		//logger.Debug("SENDER USER ID: %s | OPCODE: %d | DATA: %s", senderUserId, op, data)

		switch op {
		case OpMove:
			var dir MoveBody
			if err := json.Unmarshal(data, &dir); err != nil {
				logger.Error("error unmarshalling move body: %v", err)
				continue
			}
			if !state.snakeGame.validateDirection(snake, dir) {
				feedbackContents := "invalid direction received"
				logger.Error(feedbackContents)
				bcFeedback(dispatcher, feedbackContents, getJustSender(state, senderUserId))
			} else {
				snake.Direction = dir
			}

		default:
			feedbackContents := fmt.Sprintf("unsupported opcode received: (%d)", op)
			logger.Error(feedbackContents)
			bcFeedback(dispatcher, feedbackContents, getJustSender(state, senderUserId))
		}

		if !state.playing {
			return nil
		}
	}

	return state
}

func (m *SMatch) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, data string) (interface{}, string) {
	state := state_.(*SMatchState)

	if data == "kill" {
		return nil, "killing match due to rpc signal"
	}

	return state, ""
}

func (m *SMatch) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, graceSeconds int) interface{} {
	state := state_.(*SMatchState)

	return state
}
