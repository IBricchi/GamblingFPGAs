/*
 * jsonDecode.c
 *
 *  Created on: 15 Mar 2021
 *      Author: IBricchi
 */

#include "jsonDecode.h"

void setupInputData(){
	int a_length = 5;
	int w_length = 100;
	inputData.availableNextMoves = malloc(sizeof(char[a_length][w_length]));
	for(int i = 0; i < a_length; i++){
		inputData.availableNextMoves[i] = malloc(sizeof(char) * w_length);
	}
}

int checkName(FILE* in, char* c, char* name, int size){
    // check that name is correct
    for(int i = 0; i < size; i++){
        *c = getc(in);
        if(*c != name[i]){
            return 1;
        }
    }
    // read semicolon
    *c = getc(in);
    if(*c != ':') return 1;
    // if everything checks out
    return 0;
}

// assumes that if string is opened it will terminate
// does not check for un closed strings
// assumes that enough space is avaiable in store for string
int readStr(FILE* in, char* c, char* store){
    if(*c != '"')
        return 1;
    int i = 0;
    *c = getc(in);
    while(*c != '"'){
        store[i++] = *c;
        *c = getc(in);
    }
    store[i] = '\0';
    *c = getc(in);
    return 0;
}

int readBoolInput(FILE* in, char* c, char* name, int size, int* store){
    // check name and colon
    if(checkName(in, c, name, size))
        return 1;

    // read value
    *c = getc(in);
    // check if true
    if(*c == 't' | *c == 'T'){
        const char* comp = "rue";
        for(int i = 0; i < 3; i++){
            *c = getc(in);
            if(*c != comp[i])
                return 1;
        }
        *store = 1;
        *c = getc(in);
        return 0;
    }
    // check if false
    else if(*c == 'f' | *c == 'F'){
        const char* comp = "alse";
        for(int i = 0; i < 4; i++){
            *c = getc(in);
            if(*c != comp[i])
                return 1;
        }
        *store = 0;
        *c = getc(in);
        return 0;
    }
    // anything else is an error
    return 1;
}

int readIntInput(FILE* in, char* c, char* name, int size, int* store){
    // check name and semi colon
    if(checkName(in, c, name, size) == 1)
        return 1;

    // read value
    *store = 0;
    *c = getc(in);
    if('0' > *c | *c > '9')
        return 1;
    do{
        *store = *store * 10 + *c - '0';
        *c = getc(in);
    }
    while('0' <= *c & *c <= '9');
    return 0;
}

int readStrInput(FILE* in, char* c, char* name, int size, char* store){
    if(checkName(in, c, name, size) == 1)
        return 1;

    // get value
    *c = getc(in);
    if(readStr(in, c, store) == 1)
        return 1;

    return 0;
}

// if starting [ is found assumes remaining ] will be found
// does not check for unclosed ]
// assumes that stroe will have enough space for array
int readStrArrayInput(FILE* in, char* c, char* name, int size, char** store, int* storeCount){
    if(checkName(in, c, name, size) == -1)
        return -1;

    // get value
    *c = getc(in);
    if(*c != '[')
        return 1;
    *c = getc(in);
    int i = 0;
    while(*c != ']'){
        if(readStr(in, c, store[i++]))
            return 1;
        if(*c != ',' & *c != ']')
            return 1;
        if(*c == ',')
            *c = getc(in);
    }
    *storeCount = i;
    *c = getc(in);
    return 0;
}

int readInput(FILE* in, InputData* out){
    char c = getc(in);
    // errors return 1
    if(c == '{'){
        if(readBoolInput(in, &c, "\"isTurn\"", 8, &out->isTurn))
            return 1;
        // check for comma
        if(c != ',')
            return 1;
        if(readIntInput(in, &c, "\"currentPlayerNumber\"", 21, out->currentPlayerNumber))
            return 1;
        // check for comma
        if(c != ',')
            return 1;
        if(readStrArrayInput(in, &c, "\"availableNextMoves\"", 20, out->availableNextMoves, &out->aviablableNextMovesCount))
            return 1;
        // check for comma
        if(c != ',')
            return 1;
        if(readIntInput(in, &c, "\"moneyAvailableAmount\"", 22, &out->moneyAvailableAmount))
            return 1;
        // check for comma
        if(c != ',')
            return 1;
        if(readIntInput(in, &c, "\"minimumNextBetAmount\"", 22, &out->minimumNextBetAmount))
            return 1;
        // check for comma
        if(c != ',')
            return 1;
        if(readIntInput(in, &c, "\"relativeCardScore\"", 19, &out->relativeCardScore))
            return 1;
        // check for comma
        if(c != ',')
            return 1;
        if(readIntInput(in, &c, "\"failedPeekAttemptsCurrentGame\"", 31, &out->failedPeekAttemptsCurrentGame))
            return 1;
        // read final }
        if(c != '}')
            return 1;
        return 0;
    }
    return -1;
}
