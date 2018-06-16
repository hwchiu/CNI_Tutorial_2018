#!/bin/sh
sudo ip netns del ns1
sudo ifconfig test down
sudo brctl delbr test
sudo ip netns add ns1
go build -o example .

echo "Ready to call the step3 example"
sudo CNI_COMMAND=ADD CNI_CONTAINERID=ns1 CNI_NETNS=/var/run/netns/ns1 CNI_IFNAME=eth10 CNI_PATH=`pwd` ./example < config
echo "The CNI has been called, see the following results"
echo "The bridge and the veth has been attatch to"
sudo brctl show test
echo "The interface in the netns"
sudo ip netns exec ns1 ifconfig -a
