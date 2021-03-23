/*
 * bet.h
 *
 *  Created on: 22 Mar 2021
 *      Author: bjs3118
 */

#ifndef BET_H_
#define BET_H_

#include "globals.h"

typedef struct{
	int segvalue;
	int bcount;
	int maxQ;
	int m_digits[6];
	int bet_value[6];
} BetData;

int Bet(alt_32  *count, int *segvalue, int *maxQ, alt_32 x_value, uint16_t switch_read, uint button_read, int m_digits[], int bet_value[]);

#endif /* BET_H_ */
