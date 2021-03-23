#!/bin/bash
rm -f err.txt
rm -f out.txt

IP=18.132.52.158:3000
SRC="http://${1}:${2}@${IP}/poker/fpgaData"



DATAIN=$(curl $SRC)

if [ $(echo $DATAIN | ./checkError) ] 
then
    echo "Server responded with an error message: \"${DATAIN}\". Terminating Request."
fi

(echo $DATAIN | nios2-terminal.exe | ./main > out.txt) &
i=0
while ! [ -s out.txt ]; do
    ((i++))
    if [ $i -eq 10 ]
    then
        echo "FPGA took too long to respond. Terminating Request."
        echo "FAIL" > err.txt
        break
    fi
    echo "File empty, continue checking."
    sleep 0.1
done

if [[ -s err.txt ]
then
    cat out.txt
else
    echo "Unable to complete request."
fi
killall nios2-terminal.exe