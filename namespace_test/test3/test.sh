#!/bin/bash


### ubuntu 系统限制要打开，不然 veth0 ping veth1 会失败（arp 会没响应）
# echo 1 > /proc/sys/net/ipv4/conf/veth1/accept_local
# echo 1 > /proc/sys/net/ipv4/conf/veth0/accept_local
# echo 0 > /proc/sys/net/ipv4/conf/all/rp_filter
# echo 0 > /proc/sys/net/ipv4/conf/veth0/rp_filter
# echo 0 > /proc/sys/net/ipv4/conf/veth1/rp_filter


# 添加网桥 br0 ，也可以使用命令 sudo brctl addbr br0
sudo ip link add name br0 type bridge
sudo ip link set br0 up
bridge link
brctl show

# 创建 veth pair 网卡
sudo ip link add veth0 type veth peer name veth1
sudo ip addr add 192.168.1.201/24 dev veth0
sudo ip addr add 192.168.1.202/24 dev veth1
sudo ip link set veth0 up
sudo ip link set veth1 up
ip link

# 将 veth0 添加到网桥 br0 ，也可以使用命令 sudo brctl addif br0 veth0
sudo ip link set dev veth0 master br0
bridge link
brctl show

# shell2 另外控制台抓包
sudo tcpdump -n -i veth1

ping -c 1 -I veth0 192.168.1.202


# veth0 的 ip 给网桥
sudo ip addr del 192.168.1.201/24 dev veth0
sudo ip addr add 192.168.1.201/24 dev br0

# shell2 另外控制台抓包
sudo tcpdump -n -i veth1

ping -c 1 -I br0 192.168.1.202

# 把物理网卡 eth0 给 br0
sudo ip link set dev eth0 master br0
bridge link
brctl show
