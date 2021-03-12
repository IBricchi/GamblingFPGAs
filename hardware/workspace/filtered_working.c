#include <stdio.h>
#include <stdint.h>
#include <string.h>

#include "system.h"

#include "altera_up_avalon_accelerometer_spi.h"
#include "altera_avalon_timer_regs.h"
#include "altera_avalon_timer.h"
#include "altera_avalon_pio_regs.h"
#include "sys/alt_irq.h"



//----------------------------//
//     setup data struct      //
//----------------------------//
struct Data{
	alt_up_accelerometer_spi_dev * acc_dev;
	alt_32 acc_x_read;
	alt_32 acc_y_read;
	alt_32 acc_z_read;
	uint16_t switch_read;
	uint button_read;
}data;

//----------------------------//
//     Filter functions       //
//----------------------------//

#define FIXED_POINT_FRACTIONAL_BITS 8 //is this a global variable?? think so
typedef alt_32 fixed_point_t;


alt_32 float_to_fixed(float input)
{
		return (fixed_point_t) (input * (1 << FIXED_POINT_FRACTIONAL_BITS));
}

float fixed_to_float(alt_32 input)
{
	return ((float)input / (float) (1 << FIXED_POINT_FRACTIONAL_BITS));
}

alt_32 fixed_mult(alt_32 x, alt_32 y)
{
	return (x * y) / (1 << FIXED_POINT_FRACTIONAL_BITS);
}

void quantised_filt(alt_32 coef[], alt_32 buffer[], alt_32 x_read, alt_32 * quantised, int N)
{
	alt_32 intermediate = 0;

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


//----------------------------//
//   setup timer information  //
//----------------------------//

#define PWM_PERIOD 16

uint8_t pwm = 0;
void sys_timer_isr() {
    IOWR_ALTERA_AVALON_TIMER_STATUS(TIMER_BASE, 0);
    if (pwm > PWM_PERIOD) {
    	// get data
    	alt_up_accelerometer_spi_read_x_axis(data.acc_dev, & data.acc_x_read);
    	alt_up_accelerometer_spi_read_y_axis(data.acc_dev, & data.acc_y_read);
    	alt_up_accelerometer_spi_read_z_axis(data.acc_dev, & data.acc_z_read);
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

//----------------------------//
//  	Printing to seg       //
//----------------------------//

void print_seg(alt_32 hex)
{

	alt_32 hex_a[] = {(hex << 28)>>28, (hex << 24)>>28, (hex << 20)>>28, (hex << 16)>>28, (hex << 12)>>28, (hex << 8)>>28};

 	for(int i=0; i<6; i++)
 	{
 		switch(hex_a[i]){
 					case 0: hex_a[i] = 0b1000000; break;
 					case 1: hex_a[i] = 0b1111001; break;
 					case 2: hex_a[i] = 0b0100100; break;
 					case 3: hex_a[i] = 0b0110000; break;
 					case 4: hex_a[i] = 0b0011001; break;
 					case 5: hex_a[i] = 0b0010010; break;
 					case 6: hex_a[i] = 0b0000010; break;
 					case 7: hex_a[i] = 0b1111000; break;
 					case 8: hex_a[i] = 0b0000000; break;
 					case 9: hex_a[i] = 0b0011000; break;
 					case 10: hex_a[i] = 0b0001000; break;
 					case 11: hex_a[i] = 0b0000011; break;
 					case 12: hex_a[i] = 0b1000110; break;
 					case 13: hex_a[i] = 0b0100001; break;
 					case 14: hex_a[i] = 0b0000110; break;
 					case 15: hex_a[i] = 0b0001110; break;
 					}
 	}
 	IOWR_ALTERA_AVALON_PIO_DATA(HEX_0_BASE, hex_a[0]);
 	IOWR_ALTERA_AVALON_PIO_DATA(HEX_1_BASE, hex_a[1]);
 	IOWR_ALTERA_AVALON_PIO_DATA(HEX_2_BASE, hex_a[2]);
 	IOWR_ALTERA_AVALON_PIO_DATA(HEX_3_BASE, hex_a[3]);
 	IOWR_ALTERA_AVALON_PIO_DATA(HEX_4_BASE, hex_a[4]);
 	IOWR_ALTERA_AVALON_PIO_DATA(HEX_5_BASE, hex_a[5]);

}

//----------------------------//
//     find peek level        //
//----------------------------//
/*
 * max value is 7000(for given filter and values)
 * for me flat is approx -300
 * this is dependant on the coeffcicents of filter so could definatly be worth creating something that calibrates ...
 * probability value is between 0-9
 * high probability will give small difference between user peek and all peek
 * low probability will give large difference between user peek and all peek
 * working in alt_32 as my filter was but this can be changed
 */


void print_peek(alt_32 prob, alt_32 y_axis, uint16_t switch_read )
{
	// Finding angle value for all peak
	alt_32 range = 300;
	alt_32 all = ((10-(prob))*range) + 4000;

	// Checking if locked
	if(switch_read >=512){
		printf("locked\n");

	// Checking what angle
	}else{
		if( y_axis <= 250) 							    	{ printf("no peak\n");   }
		else if((y_axis > 250) && (y_axis <= all )) 		{ printf("user peek\n"); }
		else if( y_axis > all)								{ printf("all peek\n");  }
	}

}

//----------------------------//
//           shake            //
//----------------------------//

void shake(alt_32 * count, alt_32 * previous_value, alt_32 new_data )
{
	alt_32 intermediate = 0;

	if(*count == 10000 )
	{
		intermediate = abs(*previous_value - new_data);
		if(intermediate > 1000)
		{printf("Shake \n");}
		else
		{printf("No shake \n");}

		*count = -1;
		*previous_value = new_data;
	}

	*count = *count + 1;


}



int main()
{
	//setup filter
	float coef[24] = {0.00464135470656760, 0.00737747226463043, -0.00240768675012549, -0.00711018685736960, 0.00326564674118811, 6.11463173516297e-05, -0.00935761974859676, 0.00397493281996669, 0.00437887161977042, -0.0133160721439149, 0.00304771783859210, 0.0114361953193935, -0.0179286984033957, -0.00107408161324030, 0.0222597890359562, -0.0224772654507762, -0.0108744542661829, 0.0395972756447093, -0.0263221720611839, -0.0337570326573828, 0.0751987217099385, -0.0288978194901786, -0.120354853218164, 0.287921968939103};
	int N = 24;
	alt_32 qx_read;
	alt_32 quantised;
	alt_32 q_coef[24];
	alt_32 q_buffer[24];
	float q_filtered;
	for(int i = 0; i < (N); i++)
		{ q_coef[i] = float_to_fixed(coef[i]);}

	// shake & peek
	int prob=0;
	alt_32 count;
	alt_32 previous_value;



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
			prob = prompt - '0';

			printf("probability : %d\n", prob);


			while(1){
				qx_read = float_to_fixed((float)data.acc_z_read);
				quantised_filt(q_coef, q_buffer, qx_read, & quantised, N);
				//printf("data : %d\n", quantised);
				//print_seg(quantised);
				//print_peek(prob, quantised, data.switch_read);
				shake(& count, & previous_value, quantised );
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

