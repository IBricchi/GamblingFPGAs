#include <stdio.h>
#include <stdint.h>
#include <string.h>

#include "system.h"


#include "src/globals.h"
#include "src/setup.h"
#include "src/request.h"
#include "src/timerLoop.h"
#include "src/jsonDecode.h"
#include "src/bet.h"
#include "src/shake.h"
#include "src/filter.h"
#include "src/bitify.h"
#include "src/digify.h"
#include "src/printDec.h"

Data data;
InputData inputData;
OutputData outputData;
DataSrc dataSrc;



int main(){
	// clearing segments
	clear_dec();

	// setup inputData
	setupInputData();

	// setup peripherals
	if (setupPeripherals()) {
		return 1;
	}

	// setup jtag
	if(openFP()){
		fprintf(stderr, "Unable to access jtag.\n");
		return 1;
	}

	// setup timer
	setupTimerLoop();

	// start main execution
	printf("Running ..\n");

	// run request loop
	requestLoop();

	// close jtag
	closeFP();

	printf("Complete\n");

	return 0;
}
