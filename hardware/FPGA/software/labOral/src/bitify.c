/*
 * bitify.c
 *
 *  Created on: 22 Mar 2021
 *      Author: bjs3118
 */

#include "bitify.h"

void bitify(int bits[], uint16_t num )
{
	bits[0] = num & 0b0000000001;
	bits[1] = (num & 0b0000000010)/2;
	bits[2] = (num & 0b0000000100)/4;
	bits[3] = (num & 0b0000001000)/8;
	bits[4] = (num & 0b0000010000)/16;
	bits[5] = (num & 0b0000100000)/32;

}
