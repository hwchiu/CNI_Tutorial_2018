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

func createBridge(name string) (*netlink.Bridge, error) {
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

	err := netlink.LinkAdd(br)
	if err != nil && err != syscall.EEXIST {
		return nil, err
	}

	//Fetch the bridge Object, we need to use it for the veth
	l, err := netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("could not lookup %q: %v", name, err)
	}
	newBr, ok := l.(*netlink.Bridge)
	if !ok {
		return nil, fmt.Errorf("%q already exists but is not a bridge", name)
	}

	if err := netlink.LinkSetUp(br); err != nil {
		return nil, err
	}

	return newBr, nil
}

func setupVeth(netns ns.NetNS, br *netlink.Bridge, ifName string) error {
	guestIface := &current.Interface{}
	hostIface := &current.Interface{}
	mtu := 1500
	err := netns.Do(func(hostNS ns.NetNS) error {
		// create the veth pair in the container and move host end into host netns
		hostVeth, containerVeth, err := ip.SetupVeth(ifName, mtu, hostNS)
		if err != nil {
			return err
		}
		guestIface.Name = containerVeth.Name
		guestIface.Mac = containerVeth.HardwareAddr.String()
		guestIface.Sandbox = netns.Path()
		hostIface.Name = hostVeth.Name
		return nil
	})
	//
	fmt.Println(hostIface)
	if err != nil {
		return err
	}

	// need to lookup hostVeth again as its index has changed during ns move
	hostVeth, err := netlink.LinkByName(hostIface.Name)
	if err != nil {
		return fmt.Errorf("failed to lookup %q: %v", hostIface.Name, err)
	}
	hostIface.Mac = hostVeth.Attrs().HardwareAddr.String()

	// connect host veth end to the bridge
	if err := netlink.LinkSetMaster(hostVeth, br); err != nil {
		return fmt.Errorf("failed to connect %q to bridge %v: %v", hostVeth.Attrs().Name, br.Attrs().Name, err)
	}

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
	br, err := createBridge(sb.BridgeName)
	if err != nil {
		return err
	}

	//Create a Veth Pair
	///Get the NS Object
	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		return err
	}
	if err := setupVeth(netns, br, args.IfName); err != nil {
		return err
	}
	//Setup a IP address
	err = netns.Do(func(hostNS ns.NetNS) error {
		// create the veth pair in the container and move host end into host netns

		link, err := netlink.LinkByName(args.IfName)
		if err != nil {
			return err
		}
		ipv4Addr, ipv4Net, err := net.ParseCIDR(sb.IP)
		addr := &netlink.Addr{IPNet: ipv4Net, Label: ""}
		ipv4Net.IP = ipv4Addr
		if err = netlink.AddrAdd(link, addr); err != nil {
			return fmt.Errorf("failed to add IP addr %v to %q: %v", ipv4Net, args.IfName, err)
		}
		return nil
	})

	return err
}

func cmdDel(args *skel.CmdArgs) error {
	fmt.Println(args)
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}
