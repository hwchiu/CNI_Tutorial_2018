#!/bin/sh
go build -o example .
echo "Ready to call the step1 example"
sudo CNI_COMMAND=ADD CNI_CONTAINERID=ns1 CNI_NETNS=/var/run/netns/ns1 CNI_IFNAME=eth10 CNI_PATH=`pwd` ./example < config
echo "The CNI has been called, see the following results"
