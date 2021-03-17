package server

/*
	Data/state object that is send from FPGA nodes to server.
	IsActiveData is used to differentiating data that is send due to player action (active data)
	from data that is send regardless of player actions (passive data). Active data can only be send
	once during a player's turn (e.g. placing a new bet) while passive data is send multiple times during
	a player's turn (e.g. ShowCardsMe). A player's turn ends once active data has been received. Active data
	is ignored if it is not the player's turn.
	NewMoveType can be any of 'fold', 'check', 'bet', 'call', or 'raise'. Not all of these might be available,
	depending on the actions taken by the previous players. A list of the available next moves can be obtained
	from outgoingFPGAData:AvailableNextMoves.
	NewBetAmount is only used when NewMoveType is either 'bet' or 'raise'.
	NewTryPeek is only counted once during a player's turn and is ignored if NewTryPeekPlayerNumber does not
	correspond to the current player number.
	NewTryPeekPlayerNumber should correspond to CurrentPlayerNumber from outgoingFPGAData.
	ShowCardsIfPeek indicates that any peek attempts by other players should succeed.
*/
type incomingFPGAData struct {
	IsActiveData           bool   `json:"isActiveData"`
	ShowCardsMe            bool   `json:"showCardsMe"`
	ShowCardsIfPeek        bool   `json:"showCardsIfPeek"`
	NewTryPeek             bool   `json:"newTryPeek"`
	NewTryPeekPlayerNumber int    `json:"newTryPeekPlayerNumber"`
	NewMoveType            string `json:"newMoveType"`
	NewBetAmount           int    `json:"newBetAmount"`
}

/*
	Data/state object that is send from server to FPGA nodes.
*/
type outgoingFPGAData struct {
	IsTurn                        bool     `json:"isTurn"`
	CurrentPlayerNumber           int      `json:"currentPlayerNumber"`
	AvailableNextMoves            []string `json:"availableNextMoves"`
	MoneyAvailableAmount          int      `json:"moneyAvailableAmount"`
	MinimumNextBetAmount          int      `json:"minimumNextBetAmount"`
	RelativeCardScore             int      `json:"relativeCardScore"`
	FailedPeekAttemptsCurrentGame int      `json:"failedPeekAttemptsCurrentGame"`
}
