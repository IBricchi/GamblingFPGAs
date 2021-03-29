/*
 * jsonEncode.c
 *
 *  Created on: 23 Mar 2021
 *      Author: IBricchi
 */

#include "jsonEncode.h"

void writeBool(FILE* out, char* name, int val){
	if(val){
		fprintf(out, "\"%s\":true", name, val);
	}
	else{
		fprintf(out, "\"%s\":false", name, val);
	}
}

void writeInt(FILE* out, char* name, int val){
	fprintf(out, "\"%s\":%i", name, val);
}

void writeString(FILE* out, char* name, char* val){
	fprintf(out, "\"%s\":\"%s\"", name, val);
}

int writeOutput(FILE* out, OutputData* in){
	fprintf(out,"{");
	writeBool(out,"isActiveData",outputData.isActiveData);
	fprintf(out,",");
	writeBool(out,"showCardsMe",outputData.showCardsMe);
	fprintf(out,",");
	writeBool(out,"showCardsIfPeek",outputData.showCardsEveryone);
	fprintf(out,",");
	writeBool(out,"newTryPeek",outputData.newTryPeek);
	fprintf(out,",");
	writeInt(out,"newTryPeekPlayerNumber",outputData.newTryPeekPlayerNumber);
	fprintf(out,",");
	writeString(out,"newMoveType",outputData.newMoveType);
	fprintf(out,",");
	writeInt(out,"newBetAmount",outputData.newBetAmount);
	fprintf(out,"}");
}
