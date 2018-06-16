package main

import (
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"syscall"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
)

type SimpleBridge struct {
	BridgeName string `json:"bridgeName"`
	IP         string `json:"ip"`
}

func init() {
	// this ensures that main runs only on main thread (thread group leader).
	// since namespace ops (unshare, setns) are done for a single thread, we
	// must ensure that the goroutine does not jump from OS thread to thread
	runtime.LockOSThread()
}

func cmdAdd(args *skel.CmdArgs) error {
	sb := SimpleBridge{}
	if err := json.Unmarshal(args.StdinData, &sb); err != nil {
		return err
	}
	fmt.Println(sb)

	br := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: sb.BridgeName,
			MTU:  1500,
			// Let kernel use default txqueuelen; leaving it unset
			// means 0, and a zero-length TX queue messes up FIFO
			// traffic shapers which use TX queue length as the
			// default packet limit
			TxQLen: -1,
		},
	}

	err := netlink.LinkAdd(br)
	if err != nil && err != syscall.EEXIST {
		return err
	}

	if err := netlink.LinkSetUp(br); err != nil {
		return err
	}

	l, err := netlink.LinkByName(sb.BridgeName)
	if err != nil {
		return fmt.Errorf("could not lookup %q: %v", sb.BridgeName, err)
	}

	newBr, ok := l.(*netlink.Bridge)
	if !ok {
		return fmt.Errorf("%q already exists but is not a bridge", sb.BridgeName)
	}

	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		return err
	}

	hostIface := &current.Interface{}
	var handler = func(hostNS ns.NetNS) error {
		hostVeth, containerVeth, err := ip.SetupVeth(args.IfName, 1500, hostNS)
		if err != nil {
			return err
		}
		hostIface.Name = hostVeth.Name

		ipv4Addr, ipv4Net, err := net.ParseCIDR(sb.IP)
		if err != nil {
			return err
		}

		link, err := netlink.LinkByName(containerVeth.Name)
		if err != nil {
			return err
		}

		ipv4Net.IP = ipv4Addr

		addr := &netlink.Addr{IPNet: ipv4Net, Label: ""}
		if err = netlink.AddrAdd(link, addr); err != nil {
			return err
		}
		return nil
	}

	if err := netns.Do(handler); err != nil {
		return err
	}

	hostVeth, err := netlink.LinkByName(hostIface.Name)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetMaster(hostVeth, newBr); err != nil {
		return err
	}

	return nil
}

func cmdDel(args *skel.CmdArgs) error {
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}
