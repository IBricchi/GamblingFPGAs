package server

import (
	"fmt"

	"github.com/chehsunliu/poker"
)

/*
	RelativeCardScore reveals how good the player's cards are compared to the other player's cards.
	The score is between 0 and 100 with 0 being the worst and 100 being the best.
 	The score takes all counts that will appear during the duration of the game into account,
	not just the player's hand.
	TotalMoneyBetAmount refers to the current game.
	ShowCardsToPlayerNumbers refers to the players that are able to see the player's cards.
	TryPeekPlayerNumbers are the players that this player tried to peek in the current round.
*/
type player struct {
	Name                          string       `json:"name"`
	Hand                          []poker.Card `json:"hand"`
	MoneyAvailableAmount          int          `json:"moneyAvailableAmount"`
	RelativeCardScore             int          `json:"relativeCardScore"`
	VerboseScore                  string       `json:"verboseScore"`
	IsDealer                      bool         `json:"dealer"`
	IsSmallBlind                  bool         `json:"smallBlind"`
	IsBigBlind                    bool         `json:"bigBlind"`
	HasFolded                     bool         `json:"hasFolded"`
	LastBetAmount                 int          `json:"lastBetAmount"`
	TotalBetAmountCurrentRound    int          `json:"totalBetAmountCurrentRound"`
	TotalMoneyBetAmount           int          `json:"totalMoneyBetAmount"`
	AllIn                         bool         `json:"allIn"`
	ShowCardsMe                   bool         `json:"showCardsMe"`
	ShowCardsIfPeek               bool         `json:"showCardsIfPeek"`
	TryPeekPlayerNumbers          []int        `json:"tryPeekPlayerNumbers"`
	ShowCardsToPlayerNumbers      []int        `json:"showCardsToPlayerNumbers"`
	FailedPeekAttemptsCurrentGame int          `json:"failedPeekAttemptsCurrentGame"`
}

type maskedPlayer struct {
	Name                          string       `json:"name"`
	Hand                          []poker.Card `json:"hand"`
	MoneyAvailableAmount          int          `json:"moneyAvailableAmount"`
	IsDealer                      bool         `json:"dealer"`
	IsSmallBlind                  bool         `json:"smallBlind"`
	IsBigBlind                    bool         `json:"bigBlind"`
	HasFolded                     bool         `json:"hasFolded"`
	LastBetAmount                 int          `json:"lastBetAmount"`
	TotalMoneyBetAmount           int          `json:"totalMoneyBetAmount"`
	AllIn                         bool         `json:"allIn"`
	ShowCardsMe                   bool         `json:"showCardsMe"`
	FailedPeekAttemptsCurrentGame int          `json:"failedPeekAttemptsCurrentGame"`
}

func getPlayerPointerFromName(players []player, playerName string) (*player, error) {
	for i := range players {
		if players[i].Name == playerName {
			return &players[i], nil
		}
	}
	return &player{}, fmt.Errorf("server: poker: player %v is not part of the active poker game", playerName)
}

func (p *player) getMinimumBetAmount() int {
	var minimumNextBetAmount int
	if pokerGame.currentRound == 1 && p.IsSmallBlind && !pokerGame.smallBlindPlayed {
		minimumNextBetAmount = pokerGame.smallBlindAmount
	} else if pokerGame.currentRound == 1 && p.IsBigBlind && !pokerGame.bigBlindPlayed {
		minimumNextBetAmount = pokerGame.smallBlindAmount*2 - 1 // Minus one due to raise logic
	} else if pokerGame.lastBetAmountCurrentRound != 0 {
		minimumNextBetAmount = pokerGame.lastBetAmountCurrentRound
	}

	return minimumNextBetAmount
}

// Does not set pokerGame attributes. Must be set by calling functions.
func (p *player) allIn() {
	p.AllIn = true
	p.LastBetAmount = p.MoneyAvailableAmount
	p.TotalMoneyBetAmount += p.LastBetAmount
	p.MoneyAvailableAmount = 0
}

func (p *player) bet(amount int) error {
	if p.MoneyAvailableAmount < amount {
		p.allIn()
	} else if amount < p.getMinimumBetAmount() {
		return fmt.Errorf("server: poker: player %v's bet amount is smaller than the minimum bet amount", p.Name)
	} else {
		p.LastBetAmount = amount
		p.TotalMoneyBetAmount += p.LastBetAmount
		p.MoneyAvailableAmount -= p.LastBetAmount
	}

	pokerGame.lastRaisePlayerNumber = pokerGame.getPlayerNumber(p.Name)
	pokerGame.lastBetAmountCurrentRound = p.LastBetAmount
	pokerGame.maxBetAmountCurrentRound = p.LastBetAmount

	if p.IsSmallBlind && !pokerGame.smallBlindPlayed {
		pokerGame.smallBlindPlayed = true
	}

	return nil
}

func (p *player) call() {
	if p.MoneyAvailableAmount < pokerGame.lastBetAmountCurrentRound {
		p.allIn()
	} else {
		p.LastBetAmount = pokerGame.maxBetAmountCurrentRound - p.TotalBetAmountCurrentRound
		p.TotalMoneyBetAmount += p.LastBetAmount
		p.MoneyAvailableAmount -= p.LastBetAmount

		pokerGame.lastBetAmountCurrentRound = p.LastBetAmount
	}
}

func (p *player) raise(amount int) error {
	if p.MoneyAvailableAmount < amount {
		p.allIn()
	} else if amount <= p.getMinimumBetAmount() {
		return fmt.Errorf("server: poker: player %v's raise amount is smaller than or equal to the minimum bet amount", p.Name)
	} else {
		p.LastBetAmount = amount
		p.TotalMoneyBetAmount += p.LastBetAmount
		p.MoneyAvailableAmount -= p.LastBetAmount
	}

	pokerGame.lastRaisePlayerNumber = pokerGame.getPlayerNumber(p.Name)
	pokerGame.lastBetAmountCurrentRound = p.LastBetAmount
	pokerGame.maxBetAmountCurrentRound = p.LastBetAmount

	if p.IsBigBlind && !pokerGame.bigBlindPlayed {
		pokerGame.bigBlindPlayed = true
	}

	return nil
}

// Data is meant to be shown to player that this method is called on (requesting player)
func (p *player) computeMaskedPlayers(players []player) []maskedPlayer {
	playerNumber := pokerGame.getPlayerNumber(p.Name)

	maskedPlayers := make([]maskedPlayer, len(players))
	for i := range players {
		// requesting player
		if players[i].Name == p.Name {
			maskedPlayers[i] = maskedPlayer{
				Name:                          p.Name,
				Hand:                          p.Hand,
				MoneyAvailableAmount:          p.MoneyAvailableAmount,
				IsDealer:                      p.IsDealer,
				IsSmallBlind:                  p.IsSmallBlind,
				IsBigBlind:                    p.IsBigBlind,
				HasFolded:                     p.HasFolded,
				LastBetAmount:                 p.LastBetAmount,
				TotalMoneyBetAmount:           p.TotalMoneyBetAmount,
				AllIn:                         p.AllIn,
				ShowCardsMe:                   p.ShowCardsMe,
				FailedPeekAttemptsCurrentGame: p.FailedPeekAttemptsCurrentGame,
			}

			if !p.ShowCardsMe {
				maskedPlayers[i].Hand = []poker.Card{}
			}

			continue
		}

		// default case
		maskedPlayers[i] = maskedPlayer{
			Name:                 players[i].Name,
			MoneyAvailableAmount: players[i].MoneyAvailableAmount,
			IsDealer:             players[i].IsDealer,
			IsSmallBlind:         players[i].IsSmallBlind,
			IsBigBlind:           players[i].IsBigBlind,
			HasFolded:            players[i].HasFolded,
			LastBetAmount:        players[i].LastBetAmount,
			TotalMoneyBetAmount:  players[i].TotalMoneyBetAmount,
			AllIn:                players[i].AllIn,
		}

		// can requesting player see other player's cards?
		for j := range players[i].ShowCardsToPlayerNumbers {
			if players[i].ShowCardsToPlayerNumbers[j] == playerNumber {
				maskedPlayers[i].Hand = players[i].Hand
				break
			}
		}
	}

	return maskedPlayers
}
