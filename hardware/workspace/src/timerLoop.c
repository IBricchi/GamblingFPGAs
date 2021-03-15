/*
 * timerLoop.c
 *
 *  Created on: 15 Mar 2021
 *      Author: IBricchi
 */

#include "timerLoop.h"

// setup timer information
#define PWM_PERIOD 16

uint8_t pwm = 0;
void sys_timer_isr() {
    IOWR_ALTERA_AVALON_TIMER_STATUS(TIMER_BASE, 0);
    if (pwm > PWM_PERIOD) {
    	// get data
    	alt_up_accelerometer_spi_read_x_axis(dataSrc.acc, & data.acc_x_read);
    	alt_up_accelerometer_spi_read_x_axis(dataSrc.acc, & data.acc_y_read);
    	alt_up_accelerometer_spi_read_x_axis(dataSrc.acc, & data.acc_z_read);
    	data.switch_read = IORD_ALTERA_AVALON_PIO_DATA(SWITCH_BASE);
    	data.button_read = IORD_ALTERA_AVALON_PIO_DATA(BUTTON_BASE);
        pwm = 0;
    } else {
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
