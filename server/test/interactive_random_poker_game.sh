#!/bin/bash

# Infinite loop => Need to break out manually

set -eou pipefail

PLAYERS=("player1" "player2" "player3")
ADDRESS="localhost:3000"

# Play
while true; do
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

            echo "Player: ${player}"
            read -p "Press any key to make a random move ..."
            printf "\n"

            MOVE=""
            BET_AMOUNT=0

            # Determine which of two possible sets of moves is available
            AVAILABLE_MOVES=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.availableNextMoves[]')
            MIN_NEXT_BET_AMOUNT=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.minimumNextBetAmount')
            if [[ "${AVAILABLE_MOVES[@]}" =~ "check" ]]; then
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
                case $((RANDOM % 3)) in
                    0)
                        MOVE="fold"
                        ;;
                    1)
                        MOVE="call"
                        ;;
                    2)
                        MOVE="raise"
                        BET_AMOUNT=$(($MIN_NEXT_BET_AMOUNT + RANDOM % 100))
                        ;;
                esac
            fi

            # Randomly tild cards too much
            SHOW_CARDS_IF_PEEK=false
            if [ $((RANDOM % 3)) -ne 0 ]; then
                SHOW_CARDS_IF_PEEK=true
            fi

            curl --header "Content-Type: application/json; charset=UTF-8" \
             --request POST \
             --data '{"isActiveData":true,"showCardsMe":'"$SHOW_CARDS_IF_PEEK"',"showCardsIfPeek":'"$SHOW_CARDS_IF_PEEK"',"newTryPeek":false,"newTryPeekPlayerNumber":0,"newMoveType":"'"$MOVE"'","newBetAmount":'"${BET_AMOUNT}"'}' \
             http://${player}:${player}@$ADDRESS/poker/fpgaData
        done
    done
done
