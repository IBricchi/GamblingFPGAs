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
struct Data{
	alt_up_accelerometer_spi_dev * acc_dev;
	uint32_t acc_x_read;
	uint32_t acc_y_read;
	uint32_t acc_z_read;
	uint16_t switch_read;
	uint button_read;
}data;

// setup timer information
#define PWM_PERIOD 16

uint8_t pwm = 0;
void sys_timer_isr() {
    IOWR_ALTERA_AVALON_TIMER_STATUS(TIMER_BASE, 0);
    if (pwm > PWM_PERIOD) {
    	// get data
    	alt_up_accelerometer_spi_read_x_axis(data.acc_dev, & data.acc_x_read);
    	alt_up_accelerometer_spi_read_x_axis(data.acc_dev, & data.acc_y_read);
    	alt_up_accelerometer_spi_read_x_axis(data.acc_dev, & data.acc_z_read);
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

int main()
{
	printf("Checking Peripherals..\n");

	// setup accelerometer
	data.acc_dev = alt_up_accelerometer_spi_open_dev("/dev/accelerometer_spi");
	if (data.acc_dev == NULL) { // if return 1, check if the spi ip name is "accelerometer_spi"
		printf("Unable to access accelerometer.\n");
		return 1;
	}

	// setup jtag
	// create file pointer to jtag_uart port
	FILE* fp;
	fp = fopen ("/dev/jtag_uart", "r+");
	if(!fp){
		printf("Unable to access jtag.\n");
		return 1;
	}

	// setup timer
	timer_init(sys_timer_isr);

	// start main execution
	printf("Running ..\n");


	char prompt = 0;
	if (fp) {
		// here 'v' is used as the character to stop the program
		while (prompt != 'v') {
			// accept the character that has been sent down
			prompt = getc(fp);
			switch(prompt){
			case 'x':
				fprintf(fp, "<data>\n%lu\n", data.acc_x_read);
				break;
			case 'y':
				fprintf(fp, "<data>\n%lu\n", data.acc_y_read);
				break;
			case 'z':
				fprintf(fp, "<data>\n%lu\n", data.acc_z_read);
				break;
			case 's':
				fprintf(fp, "<data>\n%u\n", data.switch_read);
				break;
			case 'b':
				fprintf(fp, "<data>\n%u\n", data.button_read);
				break;
			}
			if (ferror(fp)) {
				clearerr(fp);
			}
		}
		fprintf(fp, "Closing the JTAG UART file handle.\n %c",0x4);
		fclose(fp);
	}
	printf("Complete\n");

	return 0;
}
