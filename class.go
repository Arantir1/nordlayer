package main

import (
	"golang.org/x/sys/unix"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
)

// CreateClass creates object of htb class with passed paramteres
func CreateClass(ifaceIndex uint32, rate uint32) tc.Object {
	tmpMsg := tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: ifaceIndex,
		Handle:  core.BuildHandle(0x0001, 0x0011),
		Parent:  tc.HandleRoot,
		Info:    0,
	}

	tmpAttr := tc.Attribute{
		Kind: "htb",
		Htb: &tc.Htb{
			Parms: &tc.HtbOpt{
				Rate: tc.RateSpec{
					Rate: rate,
				},
				Ceil: tc.RateSpec{
					Rate: rate,
				},
			},
		},
	}

	htb_class := tc.Object{Msg: tmpMsg, Attribute: tmpAttr}

	return htb_class
}

// FilterHtbClass return htb class if exist or nil otherwise
func FilterHtbClass(classes []tc.Object) *tc.Object {
	for _, class := range classes {
		if class.Kind == "htb" {
			return &class
		}
	}
	return nil
}
