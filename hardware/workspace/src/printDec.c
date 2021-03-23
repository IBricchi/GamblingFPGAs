/*
 * printDec.c
 *
 *  Created on: 22 Mar 2021
 *      Author: bjs3118
 */

#include "printDec.h"

void print_dec(int dig, int base)
{
	int prod;
 	switch(dig){
 			case 0: prod = 0b1000000; break;
 			case 1: prod = 0b1111001; break;
 			case 2: prod = 0b0100100; break;
 			case 3: prod = 0b0110000; break;
 			case 4: prod = 0b0011001; break;
 			case 5: prod = 0b0010010; break;
 			case 6: prod = 0b0000010; break;
 			case 7: prod = 0b1111000; break;
 			case 8: prod = 0b0000000; break;
 			case 9: prod = 0b0011000; break;
 			case 10: prod = 0b1111111; break;

 	}
	switch(base){
 			case 0: IOWR_ALTERA_AVALON_PIO_DATA(HEX_0_BASE, prod); break;
 			case 1: IOWR_ALTERA_AVALON_PIO_DATA(HEX_1_BASE, prod); break;
 			case 2: IOWR_ALTERA_AVALON_PIO_DATA(HEX_2_BASE, prod); break;
 			case 3: IOWR_ALTERA_AVALON_PIO_DATA(HEX_3_BASE, prod); break;
 			case 4: IOWR_ALTERA_AVALON_PIO_DATA(HEX_4_BASE, prod); break;
 			case 5: IOWR_ALTERA_AVALON_PIO_DATA(HEX_5_BASE, prod); break;
 			}

}

void clear_dec()
{
	 IOWR_ALTERA_AVALON_PIO_DATA(HEX_0_BASE, 0b1111111);
	 IOWR_ALTERA_AVALON_PIO_DATA(HEX_1_BASE, 0b1111111);
	 IOWR_ALTERA_AVALON_PIO_DATA(HEX_2_BASE, 0b1111111);
	 IOWR_ALTERA_AVALON_PIO_DATA(HEX_3_BASE, 0b1111111);
	 IOWR_ALTERA_AVALON_PIO_DATA(HEX_4_BASE, 0b1111111);
	 IOWR_ALTERA_AVALON_PIO_DATA(HEX_5_BASE, 0b1111111);
}
