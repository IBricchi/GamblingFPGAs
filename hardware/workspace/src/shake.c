/*
 * shake.c
 *
 *  Created on: 22 Mar 2021
 *      Author: bjs3118
 */

#include "shake.h"


int shake(alt_32 * count, alt_32 * previous_value, alt_32 new_data )
{
	alt_32 intermediate = 0;

	if(*count == 10000 )
	{
		intermediate = abs(*previous_value - new_data);
		if(intermediate > 10)
		{return 1;}
		else
		{return 0;}

		*count = -1;
		*previous_value = new_data;
	}

	*count = *count + 1;


}
