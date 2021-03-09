package server

import (
	"errors"
	"math/rand"
	"time"

	"github.com/chehsunliu/poker"
)

// The active game
var pokerGameStart gameStart
var pokerGame game

// Used when a game is open but not yet started.
type gameStart struct {
	open               bool
	players            []player
	initialPlayerMoney int
	smallBlindValue    int
}

// gameCards is a slice of cards that will appear in the game.
// currentRound is an integer between 1 and 4.
// currentPlayer refers to the player slice index.
// The player slice is sorted so that index 0 refers to the first player (small blind).
type game struct {
	active          bool
	deck            *poker.Deck
	communityCards  []poker.Card
	players         []player
	currentRound    int
	currentPlayer   int
	smallBlindValue int
}

// relativeCardScore reveals how good the player's cards are compared to the other player's cards.
// The score is between 0 and 100 with 0 being the worst and 100 being the best.
// The score takes all counts that will appear during the duration of the game into account,
// not just the player's hand.
type player struct {
	Name              string       `json:"name"`
	Hand              []poker.Card `json:"hand"`
	Money             int          `json:"money"`
	RelativeCardScore int          `json:"relativeCardScore"`
	VerboseScore      string       `json:"verboseScore"`
	Dealer            bool         `json:"dealer"`
	SmallBlind        bool         `json:"smallBlind"`
	BigBlind          bool         `json:"bigBlind"`
}

// Expects a slice of players that only have the name attribute initialised.
// All other attributes will be overriden.
func initGame(players []player, initialPlayerMoney int, smallBlindValue int) (game, error) {
	if len(players) == 2 {
		return game{}, errors.New("server: poker: Need at least 2 players to play a game")
	}

	deck := poker.NewDeck()

	// Give each player two cards and initial money
	for i := range players {
		players[i].Hand = deck.Draw(2)
		players[i].Money = initialPlayerMoney
	}

	// Determine which other cards will appear in game
	communityCards := deck.Draw(5)

	players = sortPlayersAccordingToRandomBlind(players)

	players = allocateRelativeCardScores(players, communityCards)

	return game{
		active:          true,
		deck:            deck,
		communityCards:  communityCards,
		players:         players,
		currentRound:    1,
		currentPlayer:   0,
		smallBlindValue: smallBlindValue,
	}, nil
}

// Expects a slice of players of length >= 2.
func sortPlayersAccordingToRandomBlind(players []player) []player {
	// Randomly determine dealer player
	rand.Seed(time.Now().UnixNano())
	dealerPlayerIdx := rand.Intn(len(players))
	players[dealerPlayerIdx].Dealer = true

	// Set small and big blind
	// Sort players so that index 0 refers to the first player
	sortedPlayers := make([]player, len(players))
	if len(players) == 2 {
		// Exception: dealer is small blind
		sortedPlayers[0] = players[dealerPlayerIdx]
		sortedPlayers[1] = players[(dealerPlayerIdx+1)%2]

		sortedPlayers[0].SmallBlind = true
		sortedPlayers[1].BigBlind = true
	} else {
		// Dealer -> Small blind -> Big blind
		for i := range sortedPlayers {
			sortedPlayers[i] = players[(dealerPlayerIdx+1+i)%len(players)]
		}

		sortedPlayers[0].SmallBlind = true
		sortedPlayers[1].BigBlind = true
	}

	return sortedPlayers
}

func allocateRelativeCardScores(players []player, communityCards []poker.Card) []player {
	type scoreMapping struct {
		score         int
		relativeScore int
	}

	// Combine player hands with community cards and calculate absolute cards score
	scoreMappings := make([]scoreMapping, len(players))
	for i := range scoreMappings {
		cards := append(communityCards, players[i].Hand...)

		scoreMappings[i].score = int(poker.Evaluate(cards))
	}

	// Relative score is calculated by
	// 1. Dividing all scores by the highest score (as float)
	// 2. Subtracting the result from one (as lowest absolute score is the best)
	// 3. Mulitply by 100 and cast to int
	highestAbsoluteScore := 0
	for _, scoreMapping := range scoreMappings {
		if scoreMapping.score > highestAbsoluteScore {
			highestAbsoluteScore = scoreMapping.score
		}
	}
	for i := range scoreMappings {
		relativeScore := float32(scoreMappings[i].score) / float32(highestAbsoluteScore)
		scoreMappings[i].relativeScore = int((1.0 - relativeScore) * 100)
	}

	// scoreMappings idx corresponds to players idx
	for i := range players {
		players[i].RelativeCardScore = scoreMappings[i].relativeScore

		players[i].VerboseScore = poker.RankString(int32(scoreMappings[i].score))
	}

	return players
}
