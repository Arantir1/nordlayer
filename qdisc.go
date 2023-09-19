package main

import (
	"golang.org/x/sys/unix"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
)

// CreateHtbQdisc creates object of htb qdisc with passed paramteres
func CreateHtbQdisc(ifaceIndex int) tc.Object {
	tmpMsg := tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: uint32(ifaceIndex),
		Handle:  core.BuildHandle(0x1, 0x0),
		Parent:  tc.HandleRoot,
		Info:    0,
	}
	tmpAttr := tc.Attribute{
		Kind: "htb",
		Htb: &tc.Htb{
			Init: &tc.HtbGlob{
				Version: 0x3,
				Defcls:  0x10,
			},
		},
	}
	htb_qdisc := tc.Object{Msg: tmpMsg, Attribute: tmpAttr}

	return htb_qdisc
}

// SelectHtbQdisc return htb qdisc if exist or nil otherwise
func SelectHtbQdisc(qdiscs []tc.Object, ifaceIndex int) *tc.Object {
	for _, qdisc := range qdiscs {
		if qdisc.Kind == "htb" && int(qdisc.Ifindex) == ifaceIndex {
			return &qdisc
		}
	}
	return nil
}
