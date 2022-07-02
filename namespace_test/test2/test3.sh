#!/bin/bash

### centos
## yum install -y bridge-utils
### ubuntu
## apt-get install -y bridge-utils

# 添加网桥 br0
sudo brctl addbr br0

# 启动网桥 br0
sudo ip link set br0 up

# 创建 network namespace ns0 ns1 ns2
sudo ip netns add ns0
sudo ip netns add ns1
sudo ip netns add ns2

# 创建 veth peer
sudo ip link add veth0-ns type veth peer name veth0-br
sudo ip link add veth1-ns type veth peer name veth1-br
sudo ip link add veth2-ns type veth peer name veth2-br

# 将 veth 的一端移动到netns中
sudo ip link set veth0-ns netns ns0
sudo ip link set veth1-ns netns ns1
sudo ip link set veth2-ns netns ns2

# 绑定 ip
sudo ip netns exec ns0 ip link set dev lo up
sudo ip netns exec ns1 ip link set dev lo up
sudo ip netns exec ns2 ip link set dev lo up
sudo ip netns exec ns0 ifconfig veth0-ns 10.1.1.1/24 up
sudo ip netns exec ns1 ifconfig veth1-ns 10.1.1.2/24 up
sudo ip netns exec ns2 ifconfig veth2-ns 10.1.1.3/24 up

# 设置默认路由，可以通宿主机
sudo ip netns exec ns0 route add default gw 10.1.1.254 veth0-ns
sudo ip netns exec ns1 route add default gw 10.1.1.254 veth1-ns
sudo ip netns exec ns2 route add default gw 10.1.1.254 veth2-ns

# 将 veth 的另一端启动并挂载到网桥上
sudo ip link set veth0-br up
sudo ip link set veth1-br up
sudo ip link set veth2-br up
sudo brctl addif br0 veth0-br
sudo brctl addif br0 veth1-br
sudo brctl addif br0 veth2-br

sudo ip addr add 10.1.1.254/24 dev br0

# 删除网桥
sudo ifconfig veth0-br 0
sudo ifconfig veth1-br 0
sudo ifconfig veth2-br 0
sudo brctl delif br0 veth0-br
sudo brctl delif br0 veth1-br
sudo brctl delif br0 veth2-br
sudo ip link set br0 down
sudo brctl delbr br0
