# CNI_Tutorial_2018
A example repository about how to write your own CNI (Container Network Interface) plugin.

# How to test it
## Build the CNI binary
go build .

## Setup a netns and use the CNI to handle the IP address
ip netns add ns1


## Execute a CNI
```
$ sh test.sh
$ brctl show 
$ ip netns exec ns1 ifconfig -a
```

## Delete a netns
ip netnd del ns1
