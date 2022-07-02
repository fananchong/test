#!/bin/bash

# 创建 network namespace ns0 ns1
sudo ip netns add ns0
sudo ip netns add ns1

# 使用 veth pair 创建 2 张虚拟网卡；分别加到 ns0 ns1 wetowrk namespace
sudo ip link add veth0 type veth peer name veth1
sudo ip link set veth0 netns ns0
sudo ip link set veth1 netns ns1

# 绑定 ip
sudo ip netns exec ns0 ifconfig veth0 10.1.1.1/24 up
sudo ip netns exec ns1 ifconfig veth1 10.1.1.2/24 up

# 删除 network namespace
sudo ip netns delete ns0
sudo ip netns delete ns1
