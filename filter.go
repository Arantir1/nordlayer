package main

import (
	"net"

	"golang.org/x/sys/unix"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
)

// CreateFilter creates object of u32 filter with passed paramteres
func CreateFilter(ifaceIndex uint32, ip string, mask int, class_id uint32) tc.Object {
	tmpMsg := tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: ifaceIndex,
		Handle:  core.BuildHandle(0x8000, 0x0011), //???
		Parent:  core.BuildHandle(0x0001, 0x0),    //???
		Info:    0x10008,                          // protocol
	}

	tmpAttr := tc.Attribute{
		Kind: "u32",
		U32: &tc.U32{
			ClassID: &class_id,
			Sel: &tc.U32Sel{
				Flags: 0x1,
				NKeys: uint8(0x1),
				Keys: []tc.U32Key{
					{
						Mask: (0xFFFFFFFF << (32 - mask)) & 0xFFFFFFFF,
						Val:  convertIp(ip),
						Off:  0x10,
					},
				},
			},
		},
	}

	u32_filter := tc.Object{Msg: tmpMsg, Attribute: tmpAttr}

	return u32_filter
}

// convert IPv4 to hex format
func convertIp(ip string) uint32 {
	var res uint32 = 0
	for i, digit := range net.ParseIP(ip).To4() {
		res += (uint32(digit) << (8 * i))
	}
	return res
}
