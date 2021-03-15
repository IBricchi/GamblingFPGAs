/*
 * setup.c
 *
 *  Created on: 15 Mar 2021
 *      Author: IBricchi
 */

#include "setup.h"

int setupAccelerometer(){
	dataSrc.acc = alt_up_accelerometer_spi_open_dev("/dev/accelerometer_spi");
	if (dataSrc.acc == NULL) { // if return 1, check if the spi ip name is "accelerometer_spi"
		printf("ERROR: Unable to access accelerometer.\n");
		return 1;
	}
	return 0;
}

int setupPeripherals(){
	// try to setup accelerometer
	if(setupAccelerometer())
		return 1;

	return 0;
}
