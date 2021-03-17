package server

import (
	"errors"
	"fmt"

	"github.com/chehsunliu/poker"
)

// The active game
var pokerGameStart gameStart
var pokerGame game
var pokerGameShowdwon gameShowdown

// Used when a game is open but not yet started.
type gameStart struct {
	open               bool
	players            []player
	initialPlayerMoney int
	smallBlindValue    int
}

/*
	gameCards is a slice of cards that will appear in the game.
 	currentRound is an integer between 1 and 4.
 	currentPlayer refers to the player slice index.
 	The player slice is sorted so that index 0 refers to the first player (small blind).
*/
type game struct {
	active                    bool
	hasEnded                  bool
	deck                      *poker.Deck
	communityCards            []poker.Card
	players                   []player
	currentRound              int
	currentPlayer             int
	lastBetAmountCurrentRound int
	smallBlindAmount          int
}

// Index in Winners corresponds with the index in WinningMoneyAmounts
type gameShowdown struct {
	CommunityCards      []poker.Card `json:"communityCards"`
	Players             []player     `json:"players"`
	Winners             []player     `json:"winners"`
	WinningReason       string       `json:"winningReason"`
	PotMoneyAmount      int          `json:"potMoneyAmount"`
	WinningMoneyAmounts []int        `json:"winningMoneyAmounts"`
}

/*
	Expects a slice of players that only have the name attribute initialised.
 	All other attributes will be overriden.
*/
func initGame(players []player, initialPlayerMoney int, smallBlindAmount int) (game, error) {
	if len(players) < 2 {
		return game{}, errors.New("server: poker: Need at least 2 players to play a game")
	}

	deck := poker.NewDeck()

	// Give each player two cards and initial money
	for i := range players {
		players[i].Hand = deck.Draw(2)
		players[i].MoneyAvailableAmount = initialPlayerMoney

		// Initialise field to avoid null JSON
		players[i].ShowCardsToPlayerNumbers = []int{}
	}

	// Determine which other cards will appear in game
	communityCards := deck.Draw(5)

	players = sortPlayersAccordingToRandomBlind(players)

	players = allocateRelativeCardScores(players, communityCards)

	return game{
		active:                    true,
		hasEnded:                  false,
		deck:                      deck,
		communityCards:            communityCards,
		players:                   players,
		currentRound:              1,
		currentPlayer:             0,
		lastBetAmountCurrentRound: 0,
		smallBlindAmount:          smallBlindAmount,
	}, nil
}

/*
	Go to next player.
	Go to next round if last player of this round.
*/
func (g *game) next() {
	if g.currentPlayer != len(g.players)-1 {
		g.currentPlayer++
	} else if g.currentRound < 4 {
		g.currentRound++
		g.currentPlayer = 0
		g.lastBetAmountCurrentRound = 0
		resetRoundSpecificPlayerData(g.players)
	} else {
		g.hasEnded = true
		g.computeShowdownData()

		// New game is started by handlePokerGetGameShowdownData()
	}
}

func (g *game) updateWithFPGAData(player *player, data incomingFPGAData) error {
	player.ShowCardsMe = data.ShowCardsMe

	// Will be set back to false at end of round
	if data.ShowCardsIfPeek {
		player.ShowCardsIfPeek = data.ShowCardsIfPeek
	}

	if data.NewTryPeek && data.NewTryPeekPlayerNumber == g.currentPlayer {
		peekSucceeded := g.tryPeek(data.NewTryPeekPlayerNumber, g.getPlayerNumber(player.Name))
		if !peekSucceeded {
			player.FailedPeekAttemptsCurrentGame++
		}
	}

	if !data.IsActiveData {
		return nil
	}

	// check if it is the player's turn
	if pokerGame.players[pokerGame.currentPlayer].Name != player.Name {
		return fmt.Errorf("server: poker: updateGameWithFPGAData: not player %v's turn, cannot process active data", player.Name)
	}

	// Player can't do anything
	if player.AllIn {
		g.next()
	}

	if !isMoveAnAvailableNextMove(data.NewMoveType) {
		return fmt.Errorf("server: poker: move %v is not one of the available moves", data.NewMoveType)
	}

	switch data.NewMoveType {
	case "fold":
		player.HasFolded = true
	case "check":
		// Do nothing
	case "bet":
		if err := player.bet(data.NewBetAmount); err != nil {
			return fmt.Errorf("server: poker: failed to place bet")
		}
	case "call":
		player.call()
	case "raise":
		if err := player.raise(data.NewBetAmount); err != nil {
			return fmt.Errorf("server: poker: failed to place raise")
		}
	}

	g.next()

	return nil
}

// Return true if peek succeeded
func (g *game) tryPeek(peekedAtPlayerNumber int, peekingPlayerNumber int) bool {
	if g.players[peekedAtPlayerNumber].ShowCardsIfPeek {
		g.players[peekedAtPlayerNumber].ShowCardsToPlayerNumbers = append(g.players[peekedAtPlayerNumber].ShowCardsToPlayerNumbers, peekingPlayerNumber)
		return true
	}
	return false
}

func (g *game) getPlayerNumber(playerName string) int {
	for i := range g.players {
		if g.players[i].Name == playerName {
			return i
		}
	}
	return -1
}

func (g *game) getCommunityCardsCurrentRound() []poker.Card {
	switch g.currentRound {
	case 1:
		return []poker.Card{}
	case 2:
		return g.communityCards[:3]
	case 3:
		return g.communityCards[:4]
	case 4:
		return g.communityCards
	}

	return nil
}

/*
	Assumes that an active game exists that has ended
 	(relevant checks should be performed before calling this method).

	Currently, every winner receives an equal amount. Side pots are not implemented.
*/
func (g *game) computeShowdownData() {
	// Ensure that data is reset
	pokerGameShowdwon = gameShowdown{}

	pokerGameShowdwon.CommunityCards = pokerGame.communityCards
	pokerGameShowdwon.Players = pokerGame.players

	potMoneyAmount := 0
	winningCardScore := 0
	winningPlayers := []int{}
	for i, player := range pokerGameShowdwon.Players {
		potMoneyAmount += player.TotalMoneyBetAmount

		if !player.HasFolded {
			if player.RelativeCardScore > winningCardScore {
				winningCardScore = player.RelativeCardScore
				winningPlayers = []int{i}
			} else if player.RelativeCardScore == winningCardScore {
				winningPlayers = append(winningPlayers, i)
			}
		}
	}
	pokerGameShowdwon.PotMoneyAmount = potMoneyAmount

	winningMoneyAmount := potMoneyAmount / len(winningPlayers)
	for i := range winningPlayers {
		pokerGameShowdwon.Winners = append(pokerGameShowdwon.Winners, pokerGameShowdwon.Players[winningPlayers[i]])
		pokerGameShowdwon.WinningMoneyAmounts = append(pokerGameShowdwon.WinningMoneyAmounts, winningMoneyAmount)
	}

	pokerGameShowdwon.WinningReason = pokerGameShowdwon.Winners[0].VerboseScore
}

/*
	Start new game with existing players and existing small blind amount.
	This method should only be called after showdown data has been computed for the current ended game.

	Uses both pokerGame and pokerGameShowdwon objects for computing values for the new poker game.
*/
func (g *game) startNewGame() {
	// Reset game attributes
	pokerGame.hasEnded = false
	pokerGame.currentRound = 1
	pokerGame.currentPlayer = 0
	pokerGame.lastBetAmountCurrentRound = 0

	pokerGame.deck = poker.NewDeck()

	// Reset player attributes, give each player two new cards,
	for i := range pokerGame.players {
		pokerGame.players[i].HasFolded = false
		pokerGame.players[i].LastBetAmount = 0
		pokerGame.players[i].TotalMoneyBetAmount = 0
		pokerGame.players[i].AllIn = false
		pokerGame.players[i].ShowCardsMe = false
		pokerGame.players[i].ShowCardsIfPeek = false
		pokerGame.players[i].ShowCardsToPlayerNumbers = []int{}
		pokerGame.players[i].FailedPeekAttemptsCurrentGame = 0

		pokerGame.players[i].Hand = pokerGame.deck.Draw(2)

		for j := range pokerGameShowdwon.Winners {
			if pokerGame.players[i].Name == pokerGameShowdwon.Winners[j].Name {
				pokerGame.players[i].MoneyAvailableAmount += pokerGameShowdwon.WinningMoneyAmounts[j]
			}
		}
	}

	// Determine which other cards will appear in game
	communityCards := pokerGame.deck.Draw(5)

	// Move dealer button by one
	pokerGame.players = sortPlayersAccordingToBlind(pokerGame.players, (getDealerPlayerIdx(pokerGame.players)+1)%len(pokerGame.players))

	pokerGame.players = allocateRelativeCardScores(pokerGame.players, communityCards)
}