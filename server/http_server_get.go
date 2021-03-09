package server

import (
	"encoding/json"
	"net/http"
)

func (h *HttpServer) handleGetStaticTest() http.HandlerFunc {
	type staticTestData struct {
		Info string `json:"info"`
		Data []int  `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		data := staticTestData{
			Info: "Some static test data",
			Data: []int{
				1, 2, 3, 4, 5,
			},
		}
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *HttpServer) handlePokerGetGameOpenStatus() http.HandlerFunc {
	type gameOpenInfo struct {
		Open               bool     `json:"open"`
		Players            []player `json:"players"`
		PlayerNumber       int      `json:"playerNumber"`
		InitialPlayerMoney int      `json:"initialPlayerMoney"`
		SmallBlindValue    int      `json:"smallBlindValue"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handlePokerGetGameOpenStatus called")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		gameOpenInfo := gameOpenInfo{
			Open:               pokerGameStart.open,
			Players:            pokerGameStart.players,
			PlayerNumber:       len(pokerGameStart.players),
			InitialPlayerMoney: pokerGameStart.initialPlayerMoney,
			SmallBlindValue:    pokerGameStart.smallBlindValue,
		}

		if err := json.NewEncoder(w).Encode(gameOpenInfo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
