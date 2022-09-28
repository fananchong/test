#!/bin/bash

#set -x

corefile=$1
hall_http=$2
# 打印所有入参
echo $@

(gdb ./a.out ${corefile} >./${corefile}_1.log 2>&1) <<GDBEOF
set pagination 0
bt
GDBEOF

(gdb ./a.out ${corefile} >./${corefile}_2.log 2>&1) <<GDBEOF
set pagination 0
thread apply all bt
GDBEOF

filename=$(md5sum ./${corefile}_1.log | awk -F' ' '{print $1}')
echo "|||============================ cserver crash ============================|||" >"./${filename}"
cat ./cserver_cpp.info >>"./${filename}"

(gdb ./a.out ${corefile} >./${corefile}_0.log 2>&1) <<GDBEOF
info threads
GDBEOF

lwp=$(cat ./${corefile}_0.log | grep "\*" | awk -F' ' '{print $4}')
tids=($(cat ./tids.txt | tr ',' ' '))
latest_area_ids=($(cat ./latest_area_id.txt | tr ',' ' '))
latest_player_ids=($(cat ./latest_player_id.txt | tr ',' ' '))

for i in $(seq 0 $(expr ${#tids[@]} - 1)); do
    _tid=${tids[i]}
    _latest_area_id=${latest_area_ids[i]}
    _latest_player_id=${latest_player_ids[i]}
    if [ "${_tid}" == "${lwp}" ]; then
        echo "${_latest_area_id}" >./latest_area_id.txt
        echo "${_latest_player_id}" >./latest_player_id.txt
    fi
done

if [ -f "./latest_area_id.txt" ]; then
    cat ./latest_area_id.txt >>"./${filename}"
fi
playerId=""
if [ -f "./latest_player_id.txt" ]; then
    playerId=$(cat ./latest_player_id.txt)
    echo playerId:${playerId} >>"./${filename}"
fi

echo "TAG: ${TAG}" >>"./${filename}"

cat ./${corefile}_2.log >>"./${filename}"

cat ./${corefile}_1.log
