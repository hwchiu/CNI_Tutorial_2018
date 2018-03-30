# CNI_Tutorial_2018
A example repository about how to write your own CNI (Container Network Interface) plugin.

# How to test it
To build the CNI binary
```sh
$ vagrant up
$ vagrant ssh
$ cd $GOPATH/src/github.com/hwchiu/CNI_Tutorial_2018
$ go build .
```

Setup a netns and use the CNI to handle the IP address
```sh
$ ip netns add ns1
```

Execute a CNI:
```sh
$ sh test.sh
$ brctl show
$ ip netns exec ns1 ifconfig -a
```

Delete a netns:
```sh
$ ip netnd del ns1
```

Delete a bridge:
```sh
$ ifconfig test down
$ brctl delbr test
```
