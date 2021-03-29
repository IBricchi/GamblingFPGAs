#!/bin/bash

# Infinite loop => Need to break out manually

set -eou pipefail

PLAYERS=("player1" "player2")
ADDRESS="localhost:3000"

# Seed for RANDOM
RANDOM=$(date +%s)

# Add duplicates to PLAYERS for more peeking
for (( i=${#PLAYERS[@]}-1; i>=0; i-- )); do
    PLAYERS+=(${PLAYERS[i]})
done

# Play
GAME_NUMBER=0
while true; do
    # Leave time for starting new game
    HAS_ENDED=true
    while $HAS_ENDED; do 
        HAS_ENDED=$(curl -s --show-error http://${PLAYERS[0]}:${PLAYERS[0]}@$ADDRESS/poker/activeGameStatus | jq -r '.hasEnded')
    done

    GAME_NUMBER=$((GAME_NUMBER+1))
    tput setaf 1; echo    "Game number $GAME_NUMBER:"; tput sgr0; printf "\n"

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

            MOVE=""
            BET_AMOUNT=0

            # Determine which of two possible sets of moves is available
            AVAILABLE_MOVES=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.availableNextMoves[]')
            MIN_NEXT_BET_AMOUNT=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.minimumNextBetAmount')
            if [[ "${AVAILABLE_MOVES[@]}" =~ "bet" ]] && [[ ! "${AVAILABLE_MOVES[@]}" =~ "check" ]]; then
                # Small blind
                MOVE="bet"
                BET_AMOUNT=$MIN_NEXT_BET_AMOUNT
                echo "bet => amount=$BET_AMOUNT"
            elif [[ "${AVAILABLE_MOVES[@]}" =~ "raise" ]] && [[ ! "${AVAILABLE_MOVES[@]}" =~ "fold" ]]; then
                # Big blind
                MOVE="raise"
                BET_AMOUNT=$(($MIN_NEXT_BET_AMOUNT + 1))
                echo "raise => amount=$BET_AMOUNT"
            elif [[ "${AVAILABLE_MOVES[@]}" =~ "check" ]]; then
                case $((RANDOM % 2)) in
                    0)
                        MOVE="check"
                        echo "check"
                        ;;
                    1)
                        MOVE="bet"
                        BET_AMOUNT=$(($MIN_NEXT_BET_AMOUNT + RANDOM % 100))
                        echo "bet => amount=$BET_AMOUNT"
                        ;;
                esac
            else
                case $((RANDOM % 6)) in
                    0)
                        MOVE="fold"
                        echo "fold"
                        ;;
                    1|2|3)
                        MOVE="call"
                        echo "call"
                        ;;
                    4|5)
                        MOVE="raise"
                        BET_AMOUNT=$(($MIN_NEXT_BET_AMOUNT + 1 + RANDOM % 100))
                        echo "raise => amount=$BET_AMOUNT"
                        ;;
                esac
            fi

            printf "\n"

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
