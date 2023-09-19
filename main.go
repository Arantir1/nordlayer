package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/mdlayher/netlink"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
)

func main() {
	// Parse arguments or set default
	iface_name := flag.String("iface", "wlp5s0", "main interface")
	ip := flag.String("ip", "80.249.99.148", "ip adress for reducing bandwidth")
	mask := flag.Int("mask", 32, "mask of ip adress")
	kbits := flag.Int("kbits", 100, "speed in kbits")
	flag.Parse()
	cmd := flag.Args()[0]

	if cmd == "help" {
		help()
		return
	}

	// open a rtnetlink socket
	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rtnetlink socket: %v\n", err)
		return
	}
	defer func() {
		if err := rtnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	// For enhanced error messages from the kernel, it is recommended to set
	// option `NETLINK_EXT_ACK`, which is supported since 4.12 kernel.
	//
	// If not supported, `unix.ENOPROTOOPT` is returned.
	err = rtnl.SetOption(netlink.ExtendedAcknowledge, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not set option ExtendedAcknowledge: %v\n", err)
		return
	}

	// get network interface by name
	iface, err := net.InterfaceByName(*iface_name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not find interface %s: %v\n", *iface_name, err)
		return
	}

	// fetch all the qdiscs from all interfaces
	qdiscs, err := rtnl.Qdisc().Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get qdiscs: %v\n", err)
		return
	}

	var htb_qdisc tc.Object
	// get htb qdisc if exist or create it otherwise
	if tmp := SelectHtbQdisc(qdiscs, iface.Index); tmp == nil {
		htb_qdisc = CreateHtbQdisc(iface.Index)
		if err := rtnl.Qdisc().Add(&htb_qdisc); err != nil {
			fmt.Fprintf(os.Stderr, "could not add new qdisc: %v\n", err)
		}
	} else {
		htb_qdisc = *tmp
	}

	// get all the classes from interface
	classes, err := rtnl.Class().Get(&tc.Msg{Ifindex: uint32(iface.Index)})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get classes: %v\n", err)
		return
	}

	var htb_class tc.Object
	rate := uint32(*kbits * 128) // convert kbits to bytes

	// get class if exist or create it otherwise
	if tmp := FilterHtbClass(classes); tmp == nil {
		htb_class = CreateClass(uint32(iface.Index), rate)
		if err := rtnl.Class().Add(&htb_class); err != nil {
			fmt.Fprintf(os.Stderr, "could not add new class: %v\n", err)
		}
	} else {
		htb_class = *tmp
	}

	// get reqiered filter from interface
	u32_filters, err := rtnl.Filter().Get(&tc.Msg{Ifindex: uint32(iface.Index), Handle: core.BuildHandle(0x8000, 0x0011)})
	var u32_filter tc.Object

	switch cmd {
	case "create":
		if len(u32_filters) > 2 {
			u32_filter = u32_filters[2]
			fmt.Fprintf(os.Stdout, "Filter already exist: %s\n", u32_filter.Attribute.Kind)
		} else {
			u32_filter = CreateFilter(uint32(iface.Index), *ip, *mask, htb_class.Handle)
			if err := rtnl.Filter().Add(&u32_filter); err != nil {
				fmt.Fprintf(os.Stderr, "Could not add new filter: %v\n", err)
			}
		}
	case "delete":
		if len(u32_filters) > 2 {
			u32_filter = u32_filters[2]
			if err := rtnl.Filter().Delete(&u32_filter); err != nil {
				fmt.Fprintf(os.Stderr, "Could not delete filter: %v\n", err)
			}
		} else {
			fmt.Println("No filter found!")
		}
	default:
		fmt.Println("Command not found!")
	}
}

func help() {
	fmt.Println(
		`Usage:  limit [ OPTIONS ] COMMAND
where  
OPTIONS := { 
	--iface=eth0   				network interface name
	--ip=80.249.99.148			IPv4 format address 
	--mask=32				mask for ip address
	--kbits=100				bandwidth limit in Kbits
}
COMMAND := {create | delete | help}`)
}
