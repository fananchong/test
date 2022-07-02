#!/bin/bash

# 创建一个名为 netns1 的 network namespace
sudo ip netns add netns1

# 使用 ip netns exec 命令进入 network namespace
sudo ip netns exec netns1 ip link list

# 进入 netns1 这个 network namespace ，把设备状态设置成 UP
sudo ip netns exec netns1 ip link set dev lo up

# 尝试 ping netns1 这个 network namespace 的 127.0.0.1
sudo ip netns exec netns1 ping 127.0.0.1

# 查看系统中有哪些 network namespace
ip netns list

# 删除 network namespace
sudo ip netns delete netns1
