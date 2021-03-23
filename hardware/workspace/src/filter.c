/*
 * filter.c
 *
 *  Created on: 22 Mar 2021
 *      Author: bjs3118
 */
#include "filter.h"


int float_to_fixed(float input)
{
		return (int) (input * (1 << 8));
}

float fixed_to_float(int input)
{
	return ((float)input / (float) (1 << 8));
}

int fixed_mult(int x, int y)
{
	return (x * y) / (1 << 4);
}

void quantised_filt(int coef[], int buffer[], int x_read, int * quantised, int N)
{
	int intermediate = 0;

	for(int i=(N-1); i>0; i--)
	        {  buffer[i] = buffer[i-1];  }
	        buffer[0] = x_read;
	for(int i=0; i<(N-1); i++)
			{ intermediate = intermediate + (fixed_mult(buffer[i], coef[i]));
			//printf("QUAN: i: %d  buf: %d   coef: %d   inter: %d\n", i, buffer[i], coef[i], intermediate);
			}
			//printf("\n qfilt: %d  unq: %d\n", intermediate, (int)fixed_to_float(intermediate));
	*quantised = intermediate;

}

void filt(float buffer[], int x_read, float * filtered, int N)
{
	float coef[24] = {0.00464135470656760, 0.00737747226463043, -0.00240768675012549, -0.00711018685736960, 0.00326564674118811, 6.11463173516297e-05, -0.00935761974859676, 0.00397493281996669, 0.00437887161977042, -0.0133160721439149, 0.00304771783859210, 0.0114361953193935, -0.0179286984033957, -0.00107408161324030, 0.0222597890359562, -0.0224772654507762, -0.0108744542661829, 0.0395972756447093, -0.0263221720611839, -0.0337570326573828, 0.0751987217099385, -0.0288978194901786, -0.120354853218164, 0.287921968939103};
	float intermediate=0;
	for(int i=(N-1); i>0; i--)
	        {  buffer[i] = buffer[i-1];  }
	        buffer[0] = (float)x_read;
	for(int i=0; i<(N-1); i++)
			{ intermediate = intermediate + (buffer[i]*coef[i]);
			//printf("FILT: i: %d  buf: %d   coef: %d   inter: %d\n", i, (int)buffer[i], (int)coef[i], (int)intermediate);
			}
			//printf("\nfilt: %d\n", (int)intermediate);
	*filtered = intermediate;

}

void filterAccelerometer(FilterData* src){
	//Filtering x-axis values
	filt(src->xbuffer, data.acc_x_read, &src->xfiltered, 24);
	//Filtering x-axis values
	filt(src->ybuffer, data.acc_y_read, &src->yfiltered, 24);
	//Filtering x-axis values
	filt(src->zbuffer, data.acc_z_read, &src->zfiltered, 24);
}
