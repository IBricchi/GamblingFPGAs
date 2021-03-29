/*
 * timerLoop.c
 *
 *  Created on: 15 Mar 2021
 *      Author: IBricchi
 */

#include "timerLoop.h"
#include "bet.h"
#include "filter.h"
#include "digify.h"
#include "printDec.h"

// setup timer information
#define PWM_PERIOD 16
int pwm = -1;

extern FILE* fp;

FilterData filterData;
BetData betData = {5,0,0,{0},{0},0};
void sys_timer_isr() {
    IOWR_ALTERA_AVALON_TIMER_STATUS(TIMER_BASE, 0);

    if (pwm > PWM_PERIOD) {
		// get data
		alt_up_accelerometer_spi_read_x_axis(dataSrc.acc, & data.acc_x_read);
		alt_up_accelerometer_spi_read_y_axis(dataSrc.acc, & data.acc_y_read);
		alt_up_accelerometer_spi_read_z_axis(dataSrc.acc, & data.acc_z_read);
		data.switch_read = IORD_ALTERA_AVALON_PIO_DATA(SWITCH_BASE);
		data.button_read = IORD_ALTERA_AVALON_PIO_DATA(BUTTON_BASE);

		//Filtering x-axis values
		filt(filterData.xbuffer, data.acc_x_read, &filterData.xfiltered, 24);
		//Filtering x-axis values
		filt(filterData.ybuffer, data.acc_y_read, &filterData.yfiltered, 24);
		//Filtering x-axis values
		filt(filterData.zbuffer, data.acc_z_read, &filterData.zfiltered, 24);

		//-----------------------------------------------//
		// Peek/tilt function --- values set in hardware //
		//-----------------------------------------------//
		// Takes yfiltered and relative card score  	  //
		// Returns 2 bits, 				  //
		// LSB is show cards me		          //
		// MSB is show cards all 			  //

		int peek = ALT_CI_TILT_0((((int)filterData.yfiltered)+30), inputData.relativeCardScore);
		outputData.showCardsMe = (peek & 0b01);
		outputData.showCardsEveryone = (peek & 0b10);    // If peek attempt calculations going on in hardware, need extra input from server

		//-----------------------------------------------//
		//            Peek attempt function 		  //
		//-----------------------------------------------//
		// checks is turn and button val		  //

		if(inputData.isTurn == 0 && data.button_read == 2)
		{
			outputData.newTryPeek = 1;
			outputData.newTryPeekPlayerNumber = inputData.currentPlayerNumber;
		}

		// TO DO: not sure how to implement new try peek player?

		//-----------------------------------------------//
		//            Bet option			  //
		//-----------------------------------------------//
		// checks not folding or all in then runs bet 	  //
		// returns integer value of bet amount           //
		// only occurs during go			  //

		if(inputData.isTurn == 1)
		{
			if(inputData.allowFold && data.button_read == 2){
				outputData.newMoveType = "fold";
				outputData.isActiveData = 1;
			}
			else if((inputData.allowCheck|inputData.allowCall) & ((data.switch_read & 0b0010000000) == 0b0010000000) & data.button_read == 1){
				outputData.newMoveType = inputData.allowCheck?"check":"call";
				outputData.isActiveData = 1;
			}
			else if(inputData.allowBet | inputData.allowRaise){
				if((data.switch_read & 0b0100000000) == 0b0100000000){
					betData.bet_total = Bet(&betData.bcount, &betData.segvalue, &betData.maxQ, filterData.xfiltered, data.switch_read, data.button_read, betData.m_digits, betData.bet_value);
				}
				else{
					betData.segvalue = 5;
					betData.bcount = 0;
					betData.maxQ = 0;
					betData.bet_total = 0;
					for(int i = 0; i < 6; i++){
						betData.bet_value[i] = 0;
						betData.m_digits[i] = 0;
						print_dec(10, i);
					}
					digify(betData.m_digits, inputData.moneyAvailableAmount);
					betData.bet_total = Bet(&betData.bcount, &betData.segvalue, &betData.maxQ, filterData.xfiltered, data.switch_read, data.button_read, betData.m_digits, betData.bet_value);
				}
				if(data.button_read == 1){
					outputData.newMoveType = inputData.allowBet?"bet":"raise";
					int b = betData.bet_total;
					if(b < inputData.moneyAvailableAmount && b >= inputData.minimumNextBetAmount)	// Fixing edge case
					{
						outputData.newBetAmount = b;
						outputData.isActiveData = 1;
					}
					else
					{
						outputData.newBetAmount = 0;
					}
				}
			}
			// check if too many peak attempts were made
			if(inputData.failedPeekAttemptsCurrentGame > 5){
				outputData.showCardsEveryone = 1;
			}
		}
		pwm = 0;
    }
    else{
    	pwm++;
    }
}

void timer_init(void * isr) {
    IOWR_ALTERA_AVALON_TIMER_CONTROL(TIMER_BASE, 0x0003);
    IOWR_ALTERA_AVALON_TIMER_STATUS(TIMER_BASE, 0);
    IOWR_ALTERA_AVALON_TIMER_PERIODL(TIMER_BASE, 0x0900);
    IOWR_ALTERA_AVALON_TIMER_PERIODH(TIMER_BASE, 0x0000);
    alt_irq_register(TIMER_IRQ, 0, isr);
    IOWR_ALTERA_AVALON_TIMER_CONTROL(TIMER_BASE, 0x0007);
}

void setupTimerLoop(){
	timer_init(sys_timer_isr);
}
