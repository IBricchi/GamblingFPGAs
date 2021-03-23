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



int main()
{
	// Filter
	float xbuffer[24];
	float xfiltered;
	float ybuffer[24];
	float yfiltered;
	float zbuffer[24];
	float zfiltered;

	// shake
	int count = 0;
	int previous_value=0;

	// bet
	int segvalue=5;
	int bcount = 0;
	int maxQ = 0;
	int m_digits[6];
	int bet_value[6];

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

	// Setting values for testing 
	inputData.IsTurn =1;
	inputData.MoneyAvailableAmount = 12345;
	inputData.RelativeCardScor = 50;
	outputData.newBetAmount = 0;

	// Commented out for testing 

	/*
	// run request loop
	requestLoop();

	// close jtag
	closeFP();

	*/
while(1){
	//-----------------------------------------------//
	// Floating point filter --- Range from -30 to 21//
	//-----------------------------------------------//
	//Filtering x-axis values
	  filt( xbuffer, data.acc_x_read, &xfiltered, 24);
	//Filtering x-axis values
	  filt( ybuffer, data.acc_y_read, &yfiltered, 24);
	//Filtering x-axis values
	  filt( zbuffer, data.acc_z_read, &zfiltered, 24);

	 //-----------------------------------------------//
	 // Peek/tilt function --- values set in hardware //
	 //-----------------------------------------------//
	 // Takes yfiltered and relative card score  	  //
	 // Returns 2 bits, 				  //
	 // LSB is show cards me		          //
	 // MSB is show cards all 			  //

	 int peek = ALT_CI_TILT_0((((int)yfiltered)+30), inputData.RelativeCardScor);
	 outputData.showCardsMe = (peek & 0b01);
	 outputData.showCardsEveryone = (peek & 0b10);    // If peek attempt calculations going on in hardware, need extra input from server


	 //-----------------------------------------------//
	 //            Peek attempt function 		  //
	 //-----------------------------------------------//
	 // checks is turn and button val		  //

	 if(inputData.IsTurn == 0 && data.button_read == 2)
	 {
		 outputData.newTryPeek = 1;
		 outputData.isActiveData = 1;
	 }
	 else
	 { outputData.newTryPeek = 0;}


	 // TO DO: not sure how to implement new try peek player?


	 //-----------------------------------------------//
	 //            Bet option			  //
	 //-----------------------------------------------//
	 // checks not folding or all in then runs bet 	  //
	 // returns integer value of bet amount           //
	 // only occurs during go			  //

	 if(inputData.IsTurn == 1)
	 {
		 if(data.button_read == 2)
		 {
			 outputData.newMoveType = 'a'; // All in
			 outputData.isActiveData = 1;

		 }

		 else if( (data.switch_read & 0b0010000000) == 0b0010000000 ) //shake was not accurate enough //shake(&count, &previous_value, zfiltered ) == 1)
		 {
			 outputData.newMoveType = 'f'; // Fold
			 outputData.isActiveData = 1;

		 }
		 else if((data.switch_read & 0b0100000000) == 0b0100000000)
		 {
			 outputData.newMoveType = 'b'; // Bet
			 digify(m_digits, inputData.MoneyAvailableAmount);
			 int b = Bet(&bcount, &segvalue, &maxQ, xfiltered, data.switch_read, data.button_read, m_digits, bet_value);
			 if(b < inputData.MoneyAvailableAmount)	// Fixing edge case
			 {
				 outputData.newBetAmount = b;
				 outputData.isActiveData = 1;
			 }
			 else
			 {outputData.newBetAmount = 0;}
		 }
	 }

	// Printing for testing
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
	printf("Complete\n");

	return 0;
}
