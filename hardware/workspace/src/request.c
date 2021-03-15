/*
 * request.c
 *
 *  Created on: 15 Mar 2021
 *      Author: IBricchi
 */

#include "request.h"

FILE* fp;

int openFP(){
	fp = fopen ("/dev/jtag_uart", "r+");
	if(!fp)
		return 1;
	return 0;
}

void closeFP(){
	fprintf(fp, "Closing the JTAG UART file handle.\n %c",0x4);
	fclose(fp);
}

void requestLoop(){
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
	}
}
