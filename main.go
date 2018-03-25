package main

import (
	"fmt"
	"encoding/json"

	"github.com/vishvananda/netlink"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
)

type SimpleBridge struct {
	BridgeName string `json:"bridgeName"`
	IP string `json:"ip"`
}

func createBridge(name string) error {
	br := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
			MTU:  1500,
			// Let kernel use default txqueuelen; leaving it unset
			// means 0, and a zero-length TX queue messes up FIFO
			// traffic shapers which use TX queue length as the
			// default packet limit
			TxQLen: -1,
		},
	}

	netlink.LinkAdd(br)
	return nil
}

func cmdAdd(args *skel.CmdArgs) error {
	fmt.Println(args.ContainerID)
	fmt.Println(args.Netns)
	fmt.Println(args.IfName)
	fmt.Println(args.Args)
	fmt.Println(args.Path)

	sb := SimpleBridge{}
	if err := json.Unmarshal(args.StdinData, &sb); err != nil {
		return err
	}
	fmt.Println(sb)

	//Create a Linux Bridge
	createBridge(sb.BridgeName)
	//Create a Veth Pair
	//Setup a IP address
	return nil
}

func cmdDel(args *skel.CmdArgs) error {
	fmt.Println(args)
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}
