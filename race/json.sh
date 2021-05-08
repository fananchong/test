#!/bin/bash

time=$(date "+%Y-%m-%d %H:%M:%S.000000")
i=0
msg="{\"Time\":\"${time}\", \"Msg\":\""
while read line; do
    if [ "${line}" == "==================" ]; then
        i=$(expr $i + 1)
    fi
    if [ "$i" == "0" ]; then
        echo ${line}
    elif [ "$i" == "1" ]; then
        msg="${msg}${line}\n"
    elif [ "$i" == "2" ]; then
        msg="${msg}\"}"
        echo $msg
        time=$(date "+%Y-%m-%d %H:%M:%S.000000")
        i=0
        msg="{\"Time\":\"${time}\", \"Msg\":\""
    fi
done
