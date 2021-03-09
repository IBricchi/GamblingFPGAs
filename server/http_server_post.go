package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *HttpServer) handlePostDynamicTest(ctx context.Context) http.HandlerFunc {
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
	}
}

// Can be called by players once a game has been opened by handlePokerOpenGame().
// Expects a string: name
func (h *HttpServer) handlePokerJoinGame() http.HandlerFunc {
	type playerName struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if !pokerGameStart.open {
			http.Error(w, "Error: can only join once a game has been started", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		data := playerName{}
		if err := decoder.Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if valid data was send
		if data.Name == "" {
			http.Error(w, "Error: player name cannot be an empty string", http.StatusBadRequest)
			return
		}

		pokerGameStart.players = append(pokerGameStart.players, player{name: data.Name})
	}
}

func (h *HttpServer) handlePokerStartGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
