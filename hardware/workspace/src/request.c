/*
 * request.c
 *
 *  Created on: 15 Mar 2021
 *      Author: IBricchi
 */

#include "request.h"
//#include "jsonDecode.h"
#include "jsonEncode.h"
#include "timerLoop.h"
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
			fscanf(fp, "%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d", &inputData.isTurn, &inputData.currentPlayerNumber, &inputData.allowFold, &inputData.allowCheck, &inputData.allowBet, &inputData.allowCall, &inputData.allowRaise, &inputData.moneyAvailableAmount, &inputData.minimumNextBetAmount, &inputData.relativeCardScore, &inputData.failedPeekAttemptsCurrentGame);

			// print output;
			fprintf(fp, "<data>\n");
			writeOutput(fp, &outputData);

			// resetting values:
			outputData.isActiveData = 0;
			outputData.newTryPeek = 0;
		}
	}
}
