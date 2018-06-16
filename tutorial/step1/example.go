package main

import (
	"fmt"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
)

func cmdAdd(args *skel.CmdArgs) error {
	fmt.Printf("interfance Name: %s\n", args.IfName)
	fmt.Printf("netns path: %s\n", args.Netns)
	fmt.Printf("the config data: %s\n", args.StdinData)
	return nil
}

func cmdDel(args *skel.CmdArgs) error {
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}
