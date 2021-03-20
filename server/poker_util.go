package server

import (
	"math/rand"
	"time"

	"github.com/chehsunliu/poker"
)

func sortPlayersAccordingToRandomBlind(players []player) []player {
	// Randomly determine dealer player
	rand.Seed(time.Now().UnixNano())
	dealerPlayerIdx := rand.Intn(len(players))

	return sortPlayersAccordingToBlind(players, dealerPlayerIdx)
}

// Expects a slice of players of length >= 2.
func sortPlayersAccordingToBlind(players []player, dealerPlayerIdx int) []player {
	// Reset all remaining dealer/blind data
	for i := range players {
		players[i].IsDealer = false
		players[i].IsSmallBlind = false
		players[i].IsBigBlind = false
	}

	players[dealerPlayerIdx].IsDealer = true

	// Set small and big blind
	// Sort players so that index 0 refers to the first player
	sortedPlayers := make([]player, len(players))
	if len(players) == 2 {
		// Exception: dealer is small blind
		sortedPlayers[0] = players[dealerPlayerIdx]
		sortedPlayers[1] = players[(dealerPlayerIdx+1)%2]

		sortedPlayers[0].IsSmallBlind = true
		sortedPlayers[1].IsBigBlind = true
	} else {
		// Dealer -> Small blind -> Big blind
		for i := range sortedPlayers {
			sortedPlayers[i] = players[(dealerPlayerIdx+1+i)%len(players)]
		}

		sortedPlayers[0].IsSmallBlind = true
		sortedPlayers[1].IsBigBlind = true
	}

	return sortedPlayers
}

func getDealerPlayerIdx(players []player) int {
	for i := range players {
		if players[i].IsDealer {
			return i
		}
	}
	return 0
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

func getAvailableNextMoves() []string {
	if pokerGame.currentRound == 1 && pokerGame.lastBetAmountCurrentRound == 0 {
		return []string{"bet"}
	} else if pokerGame.lastBetAmountCurrentRound == 0 {
		return []string{"check", "bet"}
	}
	return []string{"fold", "call", "raise"}
}

func isMoveAnAvailableNextMove(move string) bool {
	for _, availableMove := range getAvailableNextMoves() {
		if move == availableMove {
			return true
		}
	}
	return false
}

func resetRoundSpecificPlayerData(players []player) {
	for i := range players {
		players[i].ShowCardsIfPeek = false
		players[i].TryPeekPlayerNumbers = []int{}
	}
}
