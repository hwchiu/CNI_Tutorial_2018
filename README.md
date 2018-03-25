# CNI_Tutorial_2018
A example repository about how to write your own CNI (Container Network Interface) plugin.

# How to test it
## Build the CNI binary
go build .

## Setup a netns and use the CNI to handle the IP address

## Execute a CNI
```
$ sudo CNI_COMMAND=ADD CNI_CONTAINERID=ns1 CNI_NETNS=/var/run/netns/ns1 CNI_IFNAME=eth0 CNI_PATH=`pwd` .main < config
```
