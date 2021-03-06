#!/bin/bash
g++ -o main main.cpp
nios2-terminal.exe <<< $1 | ./main