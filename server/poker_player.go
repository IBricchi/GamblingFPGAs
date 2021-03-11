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
*/
type player struct {
	Name                 string       `json:"name"`
	Hand                 []poker.Card `json:"hand"`
	MoneyAvailableAmount int          `json:"moneyAvailableAmount"`
	RelativeCardScore    int          `json:"relativeCardScore"`
	VerboseScore         string       `json:"verboseScore"`
	IsDealer             bool         `json:"dealer"`
	IsSmallBlind         bool         `json:"smallBlind"`
	IsBigBlind           bool         `json:"bigBlind"`
	HasFolded            bool         `json:"hadFolded"`
	LastBetAmount        int          `json:"lastBetAmount"`
	TotalMoneyBetAmount  int          `json:"totalMoneyBetAmount"`
	AllInCurrentRound    bool         `json:"allInCurrentRound"`
	ShowCardsMe          bool         `json:"showCardsMe"`
	ShowCardsEveryone    bool         `json:"showCardsEveryone"`
}

func getPlayerPointerFromName(playerName string) (*player, error) {
	for i := range pokerGame.players {
		if pokerGame.players[i].Name == playerName {
			return &pokerGame.players[i], nil
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
	p.AllInCurrentRound = true
	p.LastBetAmount = p.MoneyAvailableAmount
	p.MoneyAvailableAmount = 0
}

func (p *player) bet(amount int) error {
	if p.MoneyAvailableAmount < amount {
		p.allIn()
	} else if amount < p.getMinimumBetAmount() {
		return fmt.Errorf("server: poker: player %v's bet amount is smaller than the minimum bet amount", p.Name)
	} else {
		p.LastBetAmount = amount
		p.MoneyAvailableAmount -= p.LastBetAmount
	}

	return nil
}

func (p *player) call() {
	if p.MoneyAvailableAmount < pokerGame.lastBetAmountCurrentRound {
		p.allIn()
	} else {
		p.LastBetAmount = pokerGame.lastBetAmountCurrentRound
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
		p.MoneyAvailableAmount -= p.LastBetAmount
	}

	return nil
}
