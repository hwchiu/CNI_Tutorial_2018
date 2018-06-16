#!/bin/sh
sudo ifconfig test down
sudo brctl delbr test
go build -o example .

echo "Ready to call the step2 example"
sudo CNI_COMMAND=ADD CNI_CONTAINERID=ns1 CNI_NETNS=/var/run/netns/ns1 CNI_IFNAME=eth10 CNI_PATH=`pwd` ./example < config
echo "The CNI has been called, see the following results"
echo "The bridge has been created"
sudo brctl show test
