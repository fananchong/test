#!/bin/bash

i=0
msg=""
while read line; do
    if [ "${line}" == "==================" ]; then
        i=$(expr $i + 1)
    fi
    if [ "$i" == "0" ]; then
        echo ${line}
    elif [ "$i" == "1" ]; then
        if [ "${msg}" == "" ]; then
            time=$(date "+%Y-%m-%d %H:%M:%S.000000")
            msg="{\"Time\":\"${time}\", \"Msg\":\""
        fi
        msg="${msg}${line}\n"
    elif [ "$i" == "2" ]; then
        msg="${msg}\"}"
        echo $msg
        i=0
        msg=""
    fi
done
