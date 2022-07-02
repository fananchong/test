#!/bin/bash

# 创建一个名为 netns1 的 network namespace
sudo ip netns add netns1

# 使用 ip netns exec 命令进入 network namespace
sudo ip netns exec netns1 ip link list

# 查看系统中有哪些 network namespace
ip netns list

# 删除 network namespace
sudo ip netns delete netns1
