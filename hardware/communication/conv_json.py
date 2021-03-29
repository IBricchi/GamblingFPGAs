import sys, json;

# {"isActiveData":false,"showCardsMe":false,"showCardsIfPeek":false,"newTryPeek":false,"newTryPeekPlayerNumber":0,"newMoveType":"(null)","newBetAmount":0}

o = {
    "isActiveData":False,
    "showCardsMe":False,
    "showCardsIfPeek":False,
    "newTryPeek":False,
    "newTryPeekPlayerNumber":0,
    "newMoveType":"",
    "newBetAmount":0
    }

data = sys.stdin.readline().split('|')

o["isActiveData"] = data[0] == '1'
o["showCardsMe"] = data[1] == '1'
o["showCardsIfPeek"] = data[2] == '1'
o["newTryPeek"] = data[3] == '1'
o["newTryPeekPlayerNumber"] = int(data[4])
o["newMoveType"] = data[5]
o["newBetAmount"] = int(data[6])

print(json.dumps(o))