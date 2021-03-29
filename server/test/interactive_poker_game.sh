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
                read -p "Should $player try to peek (y|n)? "
                if [ "${REPLY}" = "y" ] || [ "${REPLY}" = "yes" ]; then
                    CURRENT_PLAYER=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.currentPlayerNumber')

                    curl --header "Content-Type: application/json; charset=UTF-8" \
                     --request POST \
                     --data '{"isActiveData":false,"showCardsMe":false,"showCardsIfPeek":false,"newTryPeek":true,"newTryPeekPlayerNumber":'"${CURRENT_PLAYER}"'}' \
                     http://${player}:${player}@$ADDRESS/poker/fpgaData
                    
                    echo "  $player attempted to peek"
                else
                    echo "  $player did not attempt to peek"
                fi

                continue
            fi

            echo "Player: ${player}"
            AVAILABLE_MOVES=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.availableNextMoves[]')
            while true; do
                read -p "Enter a move (bet|call|check|fold|raise): "
                MOVE=${REPLY}

                if [[ "${AVAILABLE_MOVES[@]}" =~ "$MOVE" ]]; then
                    break
                fi
                echo "  $MOVE is not one of the available moves! Please try again."
            done

            BET_AMOUNT=0
            if [ "$MOVE" = "bet" ] || [ "$MOVE" = "raise" ]; then
                MIN_NEXT_BET_AMOUNT=$(curl -s --show-error http://${player}:${player}@$ADDRESS/poker/fpgaData | jq -r '.minimumNextBetAmount')
                while true; do
                    read -p "Enter the amount to bet: "
                    BET_AMOUNT=${REPLY}

                    if [ "$MOVE" = "bet" ] && [ $BET_AMOUNT -ge $MIN_NEXT_BET_AMOUNT ]; then
                        break
                    elif [ "$MOVE" = "raise" ] && [ $BET_AMOUNT -gt $MIN_NEXT_BET_AMOUNT ]; then
                        break
                    fi
                    echo "  $BET_AMOUNT is less than the minimum amount! Please try again."
                done
            fi

            echo "  Move: $MOVE, BET_AMOUNT: $BET_AMOUNT"

            SHOW_CARDS_ME=false
            read -p "Should $player look at their cards (y|n)? "
            if [ "${REPLY}" = "y" ] || [ "${REPLY}" = "yes" ]; then
                SHOW_CARDS_ME=true

                curl --header "Content-Type: application/json; charset=UTF-8" \
                 --request POST \
                 --data '{"isActiveData":false,"showCardsMe":'"$SHOW_CARDS_ME"',"showCardsIfPeek":false,"newTryPeek":false,"newTryPeekPlayerNumber":0,"newMoveType":"bet","newBetAmount":0}' \
                 http://${player}:${player}@$ADDRESS/poker/fpgaData

                echo "  $player looked at their cards"
            else
                echo "  $player did not look at their cards"
            fi

            SHOW_CARDS_IF_PEEK=false
            read -p "Should $player tild their cards too much (y|n)? "
            if [ "${REPLY}" = "y" ] || [ "${REPLY}" = "yes" ]; then
                SHOW_CARDS_IF_PEEK=true

                echo "  $player tilted their cards too much"
            else
                echo "  $player did not tild their cards too much"
            fi

            curl --header "Content-Type: application/json; charset=UTF-8" \
             --request POST \
             --data '{"isActiveData":true,"showCardsMe":'"$SHOW_CARDS_ME"',"showCardsIfPeek":'"$SHOW_CARDS_IF_PEEK"',"newTryPeek":false,"newTryPeekPlayerNumber":0,"newMoveType":"'"$MOVE"'","newBetAmount":'"${BET_AMOUNT}"'}' \
             http://${player}:${player}@$ADDRESS/poker/fpgaData

            printf "\n"
        done
    done
done
