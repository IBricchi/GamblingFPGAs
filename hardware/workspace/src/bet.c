/*
 * bet.c
 *
 *  Created on: 22 Mar 2021
 *      Author: bjs3118
 */

#include "bet.h"
#include "bitify.h"

int Bet(alt_32  *count, int *segvalue, int *maxQ, alt_32 x_value, uint16_t switch_read, uint button_read, int m_digits[], int bet_value[])
{
	// Setting up and defining
	int range = 43;
	int min = -25;
	int xval = x_value;

	// Digifying switch values
	int s_digits[6];
	bitify(s_digits, switch_read);

	// If switch locked
	if(s_digits[*segvalue] == 1 ){
		//printf("locked switch: %d ", *segvalue);
		print_dec(bet_value[*segvalue], *segvalue);
		if(bet_value[*segvalue] != m_digits[*segvalue]){
			*maxQ=1;
		}
		//printf(" %d \n", mx);
		if(*segvalue == 0){
			*segvalue = 5;
		}
		else{
			 *segvalue = *segvalue - 1;
			 if(*maxQ == 1){
				 m_digits[*segvalue] = 9;
				//printf(" %d \n", m_digits[*segvalue]);
			 }
		 }
	}
	// If switch unlocked
	else if(s_digits[*segvalue] == 0){
		if(m_digits[*segvalue] == 0){
			bet_value[*segvalue] = 0;
		}
		else{
			int intermedA = ((m_digits[*segvalue]*100)/range) + 1;
			int intermedB = (intermedA*xval) - (intermedA*min);
			bet_value[*segvalue] = m_digits[*segvalue] - (intermedB/100);
		}
		print_dec(bet_value[*segvalue], *segvalue);
	}

	return ((bet_value[5]*100000)+(bet_value[4]*10000)+(bet_value[3]*1000)+(bet_value[2]*100)+(bet_value[1]*10)+(bet_value[0]));
}

