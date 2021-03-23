/*
 * globals.h
 *
 *  Created on: 15 Mar 2021
 *      Author: IBricchi
 */

#ifndef GLOBALS_H_
#define GLOBALS_H_

#include <stdio.h>
#include <stdint.h>
#include <string.h>

#include "system.h"

#include "altera_up_avalon_accelerometer_spi.h"
#include "altera_avalon_timer_regs.h"
#include "altera_avalon_timer.h"
#include "altera_avalon_pio_regs.h"
#include "sys/alt_irq.h"

// setup data struct
typedef struct{
	uint32_t acc_x_read;
	uint32_t acc_y_read;
	uint32_t acc_z_read;
	uint16_t switch_read;
	uint8_t button_read;
} Data;
extern Data data;

typedef struct{
	int isTurn;
	int currentPlayerNumber;
	char **availableNextMoves;
	int aviablableNextMovesCount;
	int moneyAvailableAmount;
	int minimumNextBetAmount;
	int relativeCardScore;
	int failedPeekAttemptsCurrentGame;
} InputData;
extern InputData inputData;

typedef struct{
	alt_up_accelerometer_spi_dev * acc;
} DataSrc;
extern DataSrc dataSrc;

#endif /* GLOBALS_H_ */
