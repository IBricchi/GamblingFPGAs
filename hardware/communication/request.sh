#!/bin/bash
rm -f out.txt

IP=18.132.52.158:3000
SRC="http://${1}:${2}@${IP}/poker/fpgaData"

DATAIN=$(curl $SRC)

echo $DATAIN | ./checkError

if ! [ $? ] 
then
    echo "Server responded with an error message: \"${DATAIN}\". Terminating Request."
    exit 1
fi

echo $DATAIN | nios2-terminal.exe | ./readResponse > out.txt&

sleep 0.9

killall nios2-terminal.exe

DATAOUT=$(cat out.txt)
echo $DATAOUT

curl --header "Content-Type: application/json; charset=UTF-8" \
    --request POST \
    --data "${DATAOUT}" \
    $SRC