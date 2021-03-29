import sys, json;

# {"isTurn":true,"currentPlayerNumber":0,"availableNextMoves":["bet"],"moneyAvailableAmount":2000,"minimumNextBetAmount":5,"relativeCardScore":72,"failedPeekAttemptsCurrentGame":0}

i = json.load(sys.stdin)
out = ""

out += "1|" if i["isTurn"] else "0|"

out += f'{i["currentPlayerNumber"]}|'

out += "1|" if "fold" in i["availableNextMoves"] else "0|"
out += "1|" if "check" in i["availableNextMoves"] else "0|"
out += "1|" if "bet" in i["availableNextMoves"] else "0|"
out += "1|" if "call" in i["availableNextMoves"] else "0|"
out += "1|" if "raise" in i["availableNextMoves"] else "0|"

out += f'{i["moneyAvailableAmount"]}|'

out += f'{i["minimumNextBetAmount"]}|'

out += f'{i["relativeCardScore"]}|'

out += f'{i["failedPeekAttemptsCurrentGame"]}'

print(out)