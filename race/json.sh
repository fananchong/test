#!/bin/bash

time=$(date "+%Y-%m-%d %H:%M:%S.000000")
i=0
msg="{\"Time\":\"${time}\", \"Msg\":\""
while read line; do
    if [ "${line}" == "==================" ]; then
        i=$(expr $i + 1)
        msg="${msg}${line}\n"
        if [ "$i" == "2" ]; then
            msg="${msg}\"}"
            echo $msg
            time=$(date "+%Y-%m-%d %H:%M:%S.000000")
            i=0
            msg="{\"Time\":\"${time}\", \"Msg\":\""
        fi
    else
        msg="${msg}${line}\n"
    fi
done
