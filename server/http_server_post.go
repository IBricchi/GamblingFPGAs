package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *HttpServer) handlePostDynamicTest(ctx context.Context) http.HandlerFunc {
	type staticTestData struct {
		Info string `json:"info"`
		Data []int  `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		data := staticTestData{}
		if err := decoder.Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if correct data format was send
		if data.Info == "" && data.Data == nil {
			http.Error(w, "Error: Invalid data was send", http.StatusBadRequest)
			return
		}

		if err := h.db.insertTestData(ctx, data); err != nil {
			http.Error(w, "Error: Failed to insert data in DB", http.StatusInternalServerError)
			return
		}

		fmt.Println("Received data: ", data.Info, data.Data)
	}
}

// This endpoint is called to open a new game.
// Players can join by calling the handlePokerJoinGame() endpoint after a game is opened.
// Receives two ints: initialPlayerMoney, smallBlindValue.
func (h *HttpServer) handlePokerOpenGame() http.HandlerFunc {
	type pokerOpenGameData struct {
		InitialPlayerMoney int `json:"initialPlayerMoney"`
		SmallBlindValue    int `json:"smallBlindValue"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handlePokerOpenGame called")

		if pokerGameStart.open {
			http.Error(w, "Error: a poker game is already open", http.StatusBadRequest)
			return
		}
		if pokerGame.active {
			http.Error(w, "Error: a poker game is already active", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		data := pokerOpenGameData{}
		if err := decoder.Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if valid data was send
		if data.InitialPlayerMoney < 1 || data.SmallBlindValue < 1 {
			http.Error(w, "Error: initialPlayerMoney and smallBlindValue must be positive integers", http.StatusBadRequest)
			return
		}

		pokerGameStart.open = true
		pokerGameStart.initialPlayerMoney = data.InitialPlayerMoney
		pokerGameStart.smallBlindValue = data.SmallBlindValue

		h.logger.Info(fmt.Sprintf("poker game opened successfully: initialPlayerMoney=%v, smallBlindValue=%v", data.InitialPlayerMoney, data.SmallBlindValue))
	}
}

// Can be called by players once a game has been opened by handlePokerOpenGame().
// Uses the username from http basic auth for the player's name.
func (h *HttpServer) handlePokerJoinGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handlePokerJoinGame called")

		if !pokerGameStart.open {
			http.Error(w, "Error: can only join when a game is open", http.StatusBadRequest)
			return
		}

		playerName, _, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Error: getting username from http basic auth failed", http.StatusInternalServerError)
			return
		}

		pokerGameStart.players = append(pokerGameStart.players, player{name: playerName})

		h.logger.Info(fmt.Sprintf("%v joint poker game successfully", playerName))
	}
}

func (h *HttpServer) handlePokerStartGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handlePokerStartGame called")

		if !pokerGameStart.open {
			http.Error(w, "Error: can only start a game after a game has been opened", http.StatusBadRequest)
			return
		}

		if len(pokerGameStart.players) < 2 {
			http.Error(w, "Error: can only start a game after at least two players have joined", http.StatusBadRequest)
			return
		}

		game, err := initGame(pokerGameStart.players, pokerGameStart.initialPlayerMoney, pokerGameStart.smallBlindValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pokerGame = game
		pokerGameStart.open = false

		h.logger.Info("poker game started successfully")
	}
}

// Will force reset all game states.
// Should be called after a game is finished.
func (h *HttpServer) handlePokerTerminateGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handlePokerTerminateGame called")

		pokerGameStart = gameStart{}
		pokerGame = game{}

		h.logger.Info("poker game terminated successfully")
	}
}
