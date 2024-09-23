package main

import (
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

type OpCode int

// game opcodes must be POSITIVE integers
const (
	// outgoing (starting in 100)
	OpUpdate OpCode = iota + 100
	OpFeedback

	// incoming (starting in 200)
	OpMove = iota + 198
)

// outgoing
type UpdateBody = SnakeGame
type FeedbackBody string

// incoming
type MoveBody = Point

////

func getJustSender(state *SMatchState, userId string) []runtime.Presence {
	destinations := make([]runtime.Presence, 0)
	destinations = append(destinations, *state.presences[userId])
	return destinations
}

//lint:ignore U1000 optional method
func getAllButSender(state *SMatchState, userId string) []runtime.Presence {
	destinations := make([]runtime.Presence, 0)

	for k, v := range state.presences {
		if k != userId {
			destinations = append(destinations, *v)
		}
	}

	return destinations
}

////

func bcUpdate(dispatcher runtime.MatchDispatcher, state *SMatchState, destinations []runtime.Presence) {
	data, err := json.Marshal(state.snakeGame)
	if err == nil {
		dispatcher.BroadcastMessage(int64(OpUpdate), data, destinations, nil, true)
	}
}

func bcFeedback(dispatcher runtime.MatchDispatcher, body string, destinations []runtime.Presence) {
	data, err := json.Marshal(body)
	if err == nil {
		dispatcher.BroadcastMessage(int64(OpFeedback), data, destinations, nil, true)
	}
}
