/*
 * request.c
 *
 *  Created on: 15 Mar 2021
 *      Author: IBricchi
 */

#include "request.h"
#include "jsonDecode.h"
#include "jsonEncode.h"

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

			// process input strings
			inputData.allowFold = 0;
			inputData.allowCheck = 0;
			inputData.allowBet = 0;
			inputData.allowCall = 0;
			inputData.allowRaise = 0;
			for(int i = 0; i < inputData.aviablableNextMovesCount; i++){
				if(!strcmp(inputData.availableNextMoves[i], "fold")){
					inputData.allowFold = 1;
					continue;
				}
				else if(!strcmp(inputData.availableNextMoves[i], "check")){
					inputData.allowCheck = 1;
					continue;
				}
				else if(!strcmp(inputData.availableNextMoves[i], "bet")){
					inputData.allowBet = 1;
					continue;
				}
				else if(!strcmp(inputData.availableNextMoves[i], "call")){
					inputData.allowCall = 1;
					continue;
				}
				else if(!strcmp(inputData.availableNextMoves[i], "raise")){
					inputData.allowRaise = 1;
					continue;
				}
			}

			// print output;
			printf("<data>\n");
			writeOutput(fp, &outputData);

			// resetting values:
			outputData.isActiveData = 0;
			outputData.newTryPeek = 0;
		}
	}
}
