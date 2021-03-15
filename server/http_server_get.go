package server

import (
	"encoding/json"
	"net/http"

	"github.com/chehsunliu/poker"
	"golang.org/x/crypto/bcrypt"
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

func (h *HttpServer) handleGetIsAuthorised(creds map[string]string) http.HandlerFunc {
	type isAuthorised struct {
		Valid bool `json:"valid"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handleGetIsAuthorised called")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		data := isAuthorised{
			Valid: true,
		}

		user, pass, ok := r.BasicAuth()
		if !ok {
			data.Valid = false
		}

		credPassHash, credUserOk := creds[user]
		if !credUserOk {
			data.Valid = false
		}

		byteCredPassHash := []byte(credPassHash)
		bytePass := []byte(pass)
		if err := bcrypt.CompareHashAndPassword(byteCredPassHash, bytePass); err != nil {
			data.Valid = false
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Status of poker game in open phase.
func (h *HttpServer) handlePokerGetGameOpenStatus() http.HandlerFunc {
	// Open refers to the open phase while active refers to the active phase
	type gameOpenInfo struct {
		Open               bool     `json:"open"`
		Active             bool     `json:"active"`
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
			Active:             pokerGame.active,
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
		HasEnded        bool         `json:"hasEnded"`
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
			HasEnded:        pokerGame.hasEnded,
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

// Showdown data
func (h *HttpServer) handlePokerGetGameShowdownData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handlePokerGetGameShowdownData called")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if !pokerGame.active {
			http.Error(w, "Error: no active poker game exists", http.StatusBadRequest)
			return
		}
		if !pokerGame.hasEnded {
			http.Error(w, "Error: the active poker game is not in the showdown phase (has not yet ended)", http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(pokerGameShowdwon); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Start new game
		pokerGame.startNewGame()
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

		player, err := getPlayerPointerFromName(playerName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// playerAndGameData := getPlayerDataForFPGA()
		playerAndGameData := outgoingFPGAData{
			IsTurn:               pokerGame.players[pokerGame.currentPlayer].Name == player.Name,
			AvailableNextMoves:   getAvailableNextMoves(),
			MoneyAvailableAmount: player.MoneyAvailableAmount,
			MinimumNextBetAmount: player.getMinimumBetAmount(),
			RelativeCardScore:    player.RelativeCardScore,
		}

		if err := json.NewEncoder(w).Encode(playerAndGameData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
