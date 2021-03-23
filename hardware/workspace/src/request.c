/*
 * request.c
 *
 *  Created on: 15 Mar 2021
 *      Author: IBricchi
 */

#include "request.h"
#include "jsonDecode.h"

FILE* fp;

int openFP(){
	fp = fopen ("/dev/jtag_uart", "r+");
	if(!fp)
		return 1;
	return 0;
}

void closeFP(){
	fprintf(fp, "Closing the JTAG UART file handle.\n %c",0x4);
	fclose(fp);
}

void requestLoop(){
	char prompt = 0;
	if (fp) {
		while(1){
			// Setting values for testing
//			inputData.isTurn =1;
//			inputData.moneyAvailableAmount = 12345;
//			inputData.relativeCardScore = 50;
//			outputData.newBetAmount = 0;

			// read input
			if(readInput(fp, &inputData)){
				fprintf(stderr, "ERROR: Unable to parse input json.\n");
				return;
			}

			// print output;
			printf("<data>\n");
			printf("'{\"isActiveData\":%d,\"showCardsMe\":%d,\"showCardsEveryone\":%d,\"newTryPeek\":%d,\"newTryPeekPlayerNumber\":%d,\"newMoveType\":%c,\"newBetAmount\":%d}'\n",
					outputData.isActiveData,
					outputData.showCardsMe,
					outputData.showCardsEveryone,
					outputData.newTryPeek,
					outputData.newTryPeekPlayerNumber,
					outputData.newMoveType,
					outputData.newBetAmount);

			// resetting values:
			outputData.isActiveData = 0;
			outputData.newMoveType = '0';
		}
	}
}
