/*
 * digify.c
 *
 *  Created on: 22 Mar 2021
 *      Author: bjs3118
 */

#include "digify.h"

void digify(int digs[], int num )
{
	 int getrid;
	 digs[5] = num/100000;
	 getrid = (digs[5]*10);
	 digs[4] = num/10000 - getrid;
	 getrid = (getrid*10) + (digs[4]*10);
	 digs[3] = num/1000  - getrid;
	 getrid = (getrid*10) + (digs[3]*10);
	 digs[2] = num/100 - getrid;
	 getrid = (getrid*10) + (digs[2]*10);
	 digs[1] = num/10 - getrid;
	 getrid = (getrid*10) + (digs[1]*10);
	 digs[0] = num - getrid;
}
