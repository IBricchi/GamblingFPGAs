#!/bin/bash

mkdir -p bak/
cat software/fpga/hello_world.c > bak/hw_proj.c
cat hello_world.c > bak/hw_src.c
cat software/fpga/hello_world.c > hello_world.c