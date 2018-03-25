package main

import (
	"fmt"
	"encoding/json"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
)

type SimpleBridge struct {
	BrigdeName string `json:"bridgeName"`
	IP string `json:"ip"`
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
	return nil
}

func cmdDel(args *skel.CmdArgs) error {
	fmt.Println(args)
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}
