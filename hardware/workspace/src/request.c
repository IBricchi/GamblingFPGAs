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
			// read input
			if(readInput(fp, &inputData) == -1){
				fprintf(stderr, "ERROR: Unable to parse input json.");
				return;
			}

			// print output;
			// should replace with an output json file
			fprintf(fp, "<data>\n%s", inputData.AvailableNextMoves[0]);
		}
	}
}
