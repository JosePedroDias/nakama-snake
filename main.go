package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

const MODULE_NAME = "snake"
const RPC_JOIN_OR_CREATE_NAME = "snake_match"
const RPC_KILL_MATCH_NAME = "snake_kill"

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("setting up tic-tac-toe...")

	err := initializer.RegisterMatch(MODULE_NAME, newMatch)
	if err != nil {
		logger.Error("[RegisterMatch] error: ", err.Error())
		return err
	}

	if err := initializer.RegisterRpc(RPC_JOIN_OR_CREATE_NAME, SnakeMatchRPC); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc(RPC_KILL_MATCH_NAME, KillSnakeMatchesRPC); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	return nil
}

type MatchRpcBody struct {
	MatchIds []string `json:"matchIds"`
}

func SnakeMatchRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var min_size *int

	max_size := new(int)
	*max_size = 2

	var err error
	reply := MatchRpcBody{
		MatchIds: make([]string, 0),
	}

	// 1) check if a match already exists
	limit := 10
	var label string
	// https://heroiclabs.com/docs/nakama/server-framework/go-runtime/function-reference/#MatchList
	if matches, err := nk.MatchList(ctx, limit, true, label, min_size, max_size, "+label.open:1 +label.snake:1"); err != nil {
		logger.Error("[MatchList]: %s", err)
	} else {
		//logger.Warn("[MatchList]: matches: %#v", matches)
		if len(matches) > 0 {
			for _, match := range matches {
				reply.MatchIds = append(reply.MatchIds, match.MatchId)
			}
		}
	}

	if len(reply.MatchIds) == 0 {
		// 2) create a new match
		// https://heroiclabs.com/docs/nakama/server-framework/go-runtime/function-reference/#MatchCreate
		if matchId, err := nk.MatchCreate(ctx, MODULE_NAME, nil); err != nil {
			logger.Error("[MatchCreate]: %s", err)
			return "", err
		} else {
			reply.MatchIds = append(reply.MatchIds, matchId)
		}
	}

	reply2, err := json.Marshal(reply)
	if err != nil {
		return "", err
	} else {
		return string(reply2), nil
	}
}

func KillSnakeMatchesRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var minSize *int
	var maxSize *int
	var label string
	if matches, err := nk.MatchList(ctx, 20, true, label, minSize, maxSize, "+label.snake:1"); err != nil {
		logger.Error("[MatchList]: %s", err)
	} else {
		if len(matches) > 0 {
			for _, match := range matches {
				nk.MatchSignal(ctx, match.MatchId, "kill")
			}
		}
	}
	return "{}", nil
}
