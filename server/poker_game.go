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
	lastRaisePlayerNumber is the last player that bet/raised. It is 0 (first player) if no player raised this round.
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
	maxBetAmountCurrentRound  int
	smallBlindAmount          int
	lastRaisePlayerNumber     int
	smallBlindPlayed          bool
	bigBlindPlayed            bool
}

/*
	Index in Winners corresponds with the index in WinningMoneyAmounts.
	Active is true when there exists an active poker game. This stays true when the active poker game has ended.
	NewGameStarted is true when a new game has been started with the same players (handlePokerStarNewGameSamePlayers).
*/
type gameShowdown struct {
	Active              bool         `json:"active"`
	NewGameStarted      bool         `json:"newGameStarted"`
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

	sortedPlayers := sortPlayersAccordingToRandomBlind(players)

	allocateRelativeCardScores(sortedPlayers, communityCards)

	return game{
		active:                    true,
		hasEnded:                  false,
		deck:                      deck,
		communityCards:            communityCards,
		players:                   sortedPlayers,
		currentRound:              1,
		currentPlayer:             0,
		lastBetAmountCurrentRound: 0,
		maxBetAmountCurrentRound:  0,
		smallBlindAmount:          smallBlindAmount,
		lastRaisePlayerNumber:     0,
		smallBlindPlayed:          false,
		bigBlindPlayed:            false,
	}, nil
}

/*
	Go to next player.
	Go to next round if last player of this round.
*/
func (g *game) next() {
	// Check if only one player remaining
	foldedPlayerAmount := 0
	for _, player := range g.players {
		if player.HasFolded {
			foldedPlayerAmount++
		}
	}
	if foldedPlayerAmount == len(g.players)-1 {
		g.hasEnded = true
		g.computeShowdownData()
		return
	}

	if g.lastRaisePlayerNumber != (g.currentPlayer+1)%len(g.players) {
		g.currentPlayer = (g.currentPlayer + 1) % len(g.players)
	} else if g.currentRound < 4 {
		g.currentRound++
		resetRoundSpecificGameData(g)
		resetRoundSpecificPlayerData(g.players)
	} else {
		g.hasEnded = true
		g.computeShowdownData()
		return
	}

	// Skip player if no action possible
	if g.players[g.currentPlayer].AllIn || g.players[g.currentPlayer].HasFolded {
		g.next()
	}
}

func (g *game) updateWithFPGAData(player *player, data incomingFPGAData) error {
	player.ShowCardsMe = data.ShowCardsMe

	if data.ShowCardsIfPeek {
		// Will be set back to false at end of round
		player.ShowCardsIfPeek = data.ShowCardsIfPeek

		peekSucceeded := g.tryPeek(g.players, g.getPlayerNumber(player.Name))
		if !peekSucceeded {
			player.FailedPeekAttemptsCurrentGame++
		}
	} else if data.NewTryPeek && data.NewTryPeekPlayerNumber == g.currentPlayer {
		// Will be set back to []int{} at end of round
		player.TryPeekPlayerNumbers = append(player.TryPeekPlayerNumbers, data.NewTryPeekPlayerNumber)

		peekSucceeded := g.tryPeek(g.players, g.getPlayerNumber(player.Name))
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

	if !isMoveAnAvailableNextMove(data.NewMoveType) {
		return fmt.Errorf("server: poker: move %v is not one of the available moves", data.NewMoveType)
	}

	switch data.NewMoveType {
	case "fold":
		player.fold()
	case "check":
		// Do nothing
	case "bet":
		if err := player.bet(data.NewBetAmount); err != nil {
			return fmt.Errorf("server: poker: failed to place bet: %w", err)
		}
	case "call":
		player.call()
	case "raise":
		if err := player.raise(data.NewBetAmount); err != nil {
			return fmt.Errorf("server: poker: failed to place raise: %w", err)
		}
	}

	g.next()

	return nil
}

// Return true if peek succeeded
func (g *game) tryPeek(players []player, peekingPlayerNumber int) bool {
	for i := range players {
		for _, peekedAtPlayerNumber := range players[i].TryPeekPlayerNumbers {
			if g.players[peekedAtPlayerNumber].ShowCardsIfPeek {
				g.players[peekedAtPlayerNumber].ShowCardsToPlayerNumbers = append(g.players[peekedAtPlayerNumber].ShowCardsToPlayerNumbers, peekingPlayerNumber)
				return true
			}
		}
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
	pokerGameShowdwon.Active = true

	pokerGameShowdwon.CommunityCards = pokerGame.communityCards
	pokerGameShowdwon.Players = pokerGame.players

	potMoneyAmount := 0
	winningCardScore := -1
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

	// Check if won because only one player remaining
	foldedPlayerAmount := 0
	for _, player := range g.players {
		if player.HasFolded {
			foldedPlayerAmount++
		}
	}
	if foldedPlayerAmount == len(g.players)-1 {
		pokerGameShowdwon.WinningReason = "Last player remaining"
	} else {
		pokerGameShowdwon.WinningReason = pokerGameShowdwon.Winners[0].VerboseScore
	}
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
	pokerGame.maxBetAmountCurrentRound = 0
	pokerGame.smallBlindPlayed = false
	pokerGame.bigBlindPlayed = false

	pokerGame.deck = poker.NewDeck()

	// Reset player attributes, give each player two new cards,
	for i := range pokerGame.players {
		name := pokerGame.players[i].Name
		moneyAvailableAmount := pokerGame.players[i].MoneyAvailableAmount
		isDealer := pokerGame.players[i].IsDealer
		isSmallBlind := pokerGame.players[i].IsSmallBlind
		isBigBlind := pokerGame.players[i].IsBigBlind

		pokerGame.players[i] = player{
			Name:                 name,
			Hand:                 pokerGame.deck.Draw(2),
			MoneyAvailableAmount: moneyAvailableAmount,
			IsDealer:             isDealer,
			IsSmallBlind:         isSmallBlind,
			IsBigBlind:           isBigBlind,
		}

		for j := range pokerGameShowdwon.Winners {
			if pokerGame.players[i].Name == pokerGameShowdwon.Winners[j].Name {
				pokerGame.players[i].MoneyAvailableAmount += pokerGameShowdwon.WinningMoneyAmounts[j]
			}
		}

		if pokerGame.players[i].MoneyAvailableAmount == 0 {
			pokerGame.players[i].HasFolded = true
		}
	}

	// Determine which other cards will appear in game
	pokerGame.communityCards = pokerGame.deck.Draw(5)

	// Move dealer button by one
	sortedPlayers := sortPlayersAccordingToBlind(pokerGame.players, (getDealerPlayerIdx(pokerGame.players)+1)%len(pokerGame.players))

	allocateRelativeCardScores(sortedPlayers, pokerGame.communityCards)
	pokerGame.players = sortedPlayers

	pokerGameShowdwon.NewGameStarted = true
}
