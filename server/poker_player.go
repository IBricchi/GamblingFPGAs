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
	ShowCardsToPlayerNumbers refers to the players that are able to see the player's cards.
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
	TotalMoneyBetAmount           int          `json:"totalMoneyBetAmount"`
	AllIn                         bool         `json:"allIn"`
	ShowCardsMe                   bool         `json:"showCardsMe"`
	ShowCardsIfPeek               bool         `json:"showCardsIfPeek"`
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
	if pokerGame.lastBetAmountCurrentRound != 0 {
		minimumNextBetAmount = pokerGame.lastBetAmountCurrentRound
	} else if p.IsSmallBlind {
		minimumNextBetAmount = pokerGame.smallBlindAmount
	} else if p.IsBigBlind {
		minimumNextBetAmount = pokerGame.smallBlindAmount * 2
	}

	return minimumNextBetAmount
}

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

	return nil
}

func (p *player) call() {
	if p.MoneyAvailableAmount < pokerGame.lastBetAmountCurrentRound {
		p.allIn()
	} else {
		p.LastBetAmount = pokerGame.lastBetAmountCurrentRound
		p.TotalMoneyBetAmount += p.LastBetAmount
		p.MoneyAvailableAmount -= p.LastBetAmount
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
