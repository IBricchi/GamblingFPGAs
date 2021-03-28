echo "Starting"

killall nios2-terminal.exe
coproc n2t { nios2-terminal.exe; }

rm -f out.txt

IP=18.132.52.158:3000
SRC="http://${1}:${2}@${IP}/poker/fpgaData"

while :
do
    
    DATAIN=$(curl -s $SRC)

    # echo $DATAIN

    # check for error messages from server
    if [[ "${DATAIN}" != "{"* ]] 
    then
        echo "Server responded with an error message: \"${DATAIN}\". Terminating Request."
        break
    fi

    # for timing
    start=`date +%s.%N`

    echo $DATAIN >&"${n2t[1]}"

    # loop until <data> is found
    read output <&"${n2t[0]}"
    while [[ "${output}" != "<data>"* ]]
    do
        read output <&"${n2t[0]}"
    done

    read output <&"${n2t[0]}"
    # echo $output

    end=`date +%s.%N`
    echo $(echo "$end - $start" | bc -l)

    if [[ "${output}" != "{"* ]] 
    then
        echo "FPGA responded with an error message: \"${output}\". Terminating Request."
        break
    fi

    curl --header "Content-Type: application/json; charset=UTF-8" \
        --request POST \
        --data "${output}" \
        $SRC
done

echo "Ending"