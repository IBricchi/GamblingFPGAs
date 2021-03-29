#include "system.h"
#include "altera_up_avalon_accelerometer_spi.h"
#include "altera_avalon_timer_regs.h"
#include "altera_avalon_timer.h"
#include "altera_avalon_pio_regs.h"
#include "sys/alt_irq.h"
#include <stdlib.h>
//#include <stdio.h>
#define OFFSET -32
#define PWM_PERIOD 16
#define FIXED_POINT_FRACTIONAL_BITS 8		// Number of fractional bits

typedef alt_32 fixed_point_t;

alt_8 pwm = 0;
alt_u8 led;
int level;
int N = 5; // number of taps


void led_write(alt_u8 led_pattern) {
    IOWR(LED_BASE, 0, led_pattern);
}

void convert_read(alt_32 acc_read, int * level, alt_u8 * led) {
    acc_read += OFFSET;
    alt_u8 val = (acc_read >> 6) & 0x07;
    * led = (8 >> val) | (8 << (8 - val));
    * level = (acc_read >> 1) & 0x1f;
}

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

void filt(float coef[], float buffer[], alt_32 x_read, float * filtered, int N)
{
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



void sys_timer_isr() {
    IOWR_ALTERA_AVALON_TIMER_STATUS(TIMER_BASE, 0);

    if (pwm < abs(level)) {

        if (level < 0) {
            led_write(led << 1);
        } else {
            led_write(led >> 1);
        }

    } else {
        led_write(led);
    }

    if (pwm > PWM_PERIOD) {
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

int main() {
	// filter values
	int N = 24;
	float coef[24] = {0.00464135470656760, 0.00737747226463043, -0.00240768675012549, -0.00711018685736960, 0.00326564674118811, 6.11463173516297e-05, -0.00935761974859676, 0.00397493281996669, 0.00437887161977042, -0.0133160721439149, 0.00304771783859210, 0.0114361953193935, -0.0179286984033957, -0.00107408161324030, 0.0222597890359562, -0.0224772654507762, -0.0108744542661829, 0.0395972756447093, -0.0263221720611839, -0.0337570326573828, 0.0751987217099385, -0.0288978194901786, -0.120354853218164, 0.287921968939103};

	//{-0.0188274046745848, 0.00620850422084166, 0.0457717365594247, -0.0915651540471274, 0.00846817901264134, 0.542078578589013, 0.542078578589013, 0.00846817901264134, -0.0915651540471274, 0.0457717365594247, 0.00620850422084166, -0.0188274046745848};

	//{0.00464135470656760, 0.00737747226463043, -0.00240768675012549, -0.00711018685736960, 0.00326564674118811};
	float buffer[24];
	float filtered = 0;

	// quantised filter values
	alt_32 qx_read;
	alt_32 quantised;
	alt_32 q_coef[24];
	alt_32 q_buffer[24];
	float q_filtered;

	//quantising coef
	for(int i = 0; i < (N); i++)
	{
		q_coef[i] = float_to_fixed(coef[i]);
	}

	int data_out;
	int bet;
	int fir;
	int tilt4;

    alt_32 x_read;
    alt_up_accelerometer_spi_dev * acc_dev;
    acc_dev = alt_up_accelerometer_spi_open_dev("/dev/accelerometer_spi");
    if (acc_dev == NULL) { // if return 1, check if the spi ip name is "accelerometer_spi"
        return 1;
    }

    timer_init(sys_timer_isr);
    while (1) {

        alt_up_accelerometer_spi_read_x_axis(acc_dev, & x_read);
        //alt_printf("raw data: %x", x_read);


        filt(coef, buffer, x_read, & filtered, N);

//////////////////////////////////////////////////////////////////////////////////////////
	// Fixed Point filter
        qx_read = float_to_fixed((float)x_read);
        quantised_filt(q_coef, q_buffer, qx_read, & quantised, N);
        q_filtered = fixed_to_float((int)quantised);
        convert_read((int)q_filtered, & level, & led);

//////////////////////////////////////////////////////////////////////////////////////////
	// Hardware filter
	//IOWR_ALTERA_AVALON_PIO_DATA(FIR_OUT_BASE, x_read);
	//data_out = IORD_ALTERA_AVALON_PIO_DATA(FIR_IN_BASE);

        bet = ALT_CI_BET1_0(5,5);
       // data_out = ALT_CI_TILT_0(((int)filtered+30),50);
        data_out = ALT_CI_TILT3_0(((int)filtered+30),50);
        fir = ALT_CI_FIR5_0(x_read,5);
       // tilt4 =  ALT_CI_TILT4_0((int)filtered+31),50);

       printf("software :: %d  || tilt: %d  || bet: %d  || fir: %d\n", (int)filtered, data_out, bet, fir);


    }

    return 0;
}

