#!/bin/bash
rm out.txt
(echo $1 | nios2-terminal.exe | ./main > out.txt) &
while ! [ -s out.txt ]; do
    echo "file is empty - keep checking it "
    sleep 1 # throttle the check
done
echo "file is not empty "
cat out.txt
killall nios2-terminal.exe