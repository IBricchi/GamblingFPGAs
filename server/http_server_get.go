package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chehsunliu/poker"
)

type staticTestData struct {
	Info string `json:"info"`
	Data []int  `json:"data"`
}

func (h *HttpServer) handleGetStaticTest() http.HandlerFunc {
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

// Status of poker game in open phase.
func (h *HttpServer) handlePokerGetGameOpenStatus() http.HandlerFunc {
	type gameOpenInfo struct {
		Open               bool     `json:"open"`
		Players            []player `json:"players"`
		PlayerAmount       int      `json:"playerAmount"`
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
			PlayerAmount:       len(pokerGameStart.players),
			InitialPlayerMoney: pokerGameStart.initialPlayerMoney,
			SmallBlindValue:    pokerGameStart.smallBlindValue,
		}

		if err := json.NewEncoder(w).Encode(gameOpenInfo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Status of poker game in active phase.
func (h *HttpServer) handlePokerGetGameActiveStatus() http.HandlerFunc {
	type gameActiveInfo struct {
		Active          bool         `json:"active"`
		CommunityCards  []poker.Card `json:"communityCards"`
		Players         []player     `json:"players"`
		PlayerAmount    int          `json:"playerAmount"`
		CurrentRound    int          `json:"currentRound"`
		CurrentPlayer   int          `json:"currentPlayer"`
		SmallBlindValue int          `json:"smallBlindValue"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handlePokerGetGameActiveStatus called")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		gameActiveInfo := gameActiveInfo{
			Active:          pokerGame.active,
			CommunityCards:  pokerGame.communityCards,
			Players:         pokerGame.players,
			PlayerAmount:    len(pokerGame.players),
			CurrentRound:    pokerGame.currentRound,
			CurrentPlayer:   pokerGame.currentPlayer,
			SmallBlindValue: pokerGame.smallBlindAmount,
		}

		if err := json.NewEncoder(w).Encode(gameActiveInfo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Called by FPGA nodes
func (h *HttpServer) handlePokerGetFPGAData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handlePokerGetFPGAData called")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if !pokerGame.active {
			http.Error(w, "Error: no active poker game exists", http.StatusBadRequest)
			return
		}

		playerName, _, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Error: getting username from http basic auth failed", http.StatusInternalServerError)
			return
		}

		var player player
		var playerIsPartOfActiveGame bool
		for i := range pokerGame.players {
			if pokerGame.players[i].Name == playerName {
				playerIsPartOfActiveGame = true
				player = pokerGame.players[i]
			}
		}
		if !playerIsPartOfActiveGame {
			http.Error(w, fmt.Sprintf("Error: player %v not part of the active poker game", playerName), http.StatusBadRequest)
			return
		}

		// Determine available next moves
		var availableNextMoves []string
		if pokerGame.lastBetAmountCurrentRound == 0 {
			availableNextMoves = []string{"check", "bet"}
		} else {
			availableNextMoves = []string{"fold", "call", "raise"}
		}

		// Determine minimum next bet amount
		var minimumNextBetAmount int
		if pokerGame.lastBetAmountCurrentRound != 0 {
			minimumNextBetAmount = pokerGame.lastBetAmountCurrentRound
		} else if player.IsSmallBlind {
			minimumNextBetAmount = pokerGame.smallBlindAmount
		} else if player.IsBigBlind {
			minimumNextBetAmount = pokerGame.smallBlindAmount * 2
		}

		// playerAndGameData := getPlayerDataForFPGA()
		playerAndGameData := outgoingFPGAData{
			IsTurn:               pokerGame.players[pokerGame.currentPlayer].Name == player.Name,
			AvailableNextMoves:   availableNextMoves,
			MoneyAvailableAmount: player.MoneyAvailableAmount,
			MinimumNextBetAmount: minimumNextBetAmount,
			RelativeCardScore:    player.RelativeCardScore,
		}

		if err := json.NewEncoder(w).Encode(playerAndGameData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
