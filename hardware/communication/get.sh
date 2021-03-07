#!/bin/bash
rm out.txt
(echo $1 | nios2-terminal.exe | ./main > out.txt) &
i=0
while ! [ -s out.txt ]; do
    ((i++))
    if [[ $i -eq 10 ]]; then
        break
    fi
    sleep 0.1
done
cat out.txt
killall nios2-terminal.exe