#!/bin/bash

set -eou pipefail

PLAYERS=("player1" "player2" "player3")
INITIAL_PLAYER_MONEY=2000
SMALL_BLIND_VALUE=5
NUMBER_OF_GAMES=3
ADDRESS="localhost:3000"

# Seed for RANDOM
RANDOM=$(date +%s)

# Open game
curl --header "Content-Type: application/json; charset=UTF-8" \
 --request POST \
 --data '{"initialPlayerMoney":'"${INITIAL_PLAYER_MONEY}"',"smallBlindValue":'"${SMALL_BLIND_VALUE}"'}' \
 http://test:test@$ADDRESS/poker/openGame

# Join Game
for player in ${PLAYERS[@]}; do
  curl --header "Content-Type: application/json; charset=UTF-8" \
   --request POST \
   http://${player}:${player}@$ADDRESS/poker/joinGame
done

# Start Game
 curl --header "Content-Type: application/json; charset=UTF-8" \
  --request POST \
  http://test:test@$ADDRESS/poker/startGame

# Play
for (( i=0; i<$NUMBER_OF_GAMES; i++ )) do
    tput setaf 1; echo    "Game number $i:"; tput sgr0; printf "\n"

    HAS_ENDED=false
    while ! $HAS_ENDED; do
        for player in ${PLAYERS[@]}; do
            HAS_ENDED=$(curl -s --show-error http://${PLAYERS[0]}:${PLAYERS[0]}@$ADDRESS/poker/activeGameStatus | jq -r '.hasEnded')
            if $HAS_ENDED; then
                break
            fi

            IS_TURN=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.isTurn')
            if ! $IS_TURN; then 
                # Randomly try to peek
                if [ $((RANDOM % 3)) -ne 0 ]; then
                    CURRENT_PLAYER=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.currentPlayerNumber')

                    curl --header "Content-Type: application/json; charset=UTF-8" \
                     --request POST \
                     --data '{"isActiveData":false,"showCardsMe":false,"showCardsIfPeek":false,"newTryPeek":true,"newTryPeekPlayerNumber":'"${CURRENT_PLAYER}"'}' \
                     http://${player}:${player}@$ADDRESS/poker/fpgaData
                fi

                continue
            fi

            MOVE=""
            BET_AMOUNT=0

            # Determine which of two possible sets of moves is available
            AVAILABLE_MOVES=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.availableNextMoves[]')
            MIN_NEXT_BET_AMOUNT=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.minimumNextBetAmount')
            if [[ "${AVAILABLE_MOVES[@]}" =~ "bet" ]] && [[ ! "${AVAILABLE_MOVES[@]}" =~ "check" ]]; then
                # Small blind
                MOVE="bet"
                BET_AMOUNT=$MIN_NEXT_BET_AMOUNT
            elif [[ "${AVAILABLE_MOVES[@]}" =~ "raise" ]] && [[ ! "${AVAILABLE_MOVES[@]}" =~ "fold" ]]; then
                # Big blind
                MOVE="raise"
                BET_AMOUNT=$(($MIN_NEXT_BET_AMOUNT + 1))
            elif [[ "${AVAILABLE_MOVES[@]}" =~ "check" ]]; then
                case $((RANDOM % 2)) in
                    0)
                        MOVE="check"
                        ;;
                    1)
                        MOVE="bet"
                        BET_AMOUNT=$(($MIN_NEXT_BET_AMOUNT + RANDOM % 100))
                        ;;
                esac
            else
                case $((RANDOM % 6)) in
                    0)
                        MOVE="fold"
                        ;;
                    1|2|3)
                        MOVE="call"
                        ;;
                    4|5)
                        MOVE="raise"
                        BET_AMOUNT=$(($MIN_NEXT_BET_AMOUNT + 1 + RANDOM % 100))
                        ;;
                esac
            fi

            # Randomly tild cards too much
            SHOW_CARDS_IF_PEEK=false
            if [ $((RANDOM % 2)) -eq 0 ]; then
                SHOW_CARDS_IF_PEEK=true
            fi

            curl --header "Content-Type: application/json; charset=UTF-8" \
             --request POST \
             --data '{"isActiveData":true,"showCardsMe":'"$SHOW_CARDS_IF_PEEK"',"showCardsIfPeek":'"$SHOW_CARDS_IF_PEEK"',"newTryPeek":false,"newTryPeekPlayerNumber":0,"newMoveType":"'"$MOVE"'","newBetAmount":'"${BET_AMOUNT}"'}' \
             http://${player}:${player}@$ADDRESS/poker/fpgaData
        done
    done

    # Display showdown data
    curl -s --show-error http://test:test@$ADDRESS/poker/activeGameStatus/showdown | jq '.'

    # Start new game with same players
    curl --header "Content-Type: application/json; charset=UTF-8" \
     --request POST \
     http://test:test@$ADDRESS/poker/startNewGameSamePlayers
done

# Terminate game
curl --header "Content-Type: application/json; charset=UTF-8" \
 --request POST \
 http://test:test@$ADDRESS/poker/terminateGame
