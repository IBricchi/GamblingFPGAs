#!/bin/bash

rm -rf bak

mkdir -p bak/
cp -rf src/ bak/src/
cp main.c bak/main.c

cp -f software/fpga/src/* src/
cp -f software/fpga/main.c main.c