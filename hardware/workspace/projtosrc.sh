#!/bin/bash

rm -rf bak

mkdir -p bak/
cp -rf src/ bak/src/
cp main.c bak/main.c