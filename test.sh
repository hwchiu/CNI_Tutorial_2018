sudo ip netns del ns1
sudo brctl delbr test

sudo ip netns add ns1

sudo CNI_COMMAND=ADD CNI_CONTAINERID=ns1 CNI_NETNS=/var/run/netns/ns1 CNI_IFNAME=eth10 CNI_PATH=/home/vagrant/go/src/github.com/hwchiu/CNI_Tutorial_2018 ./CNI_Tutorial_2018  < config
