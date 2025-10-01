// Package ipcalc allow for calculate address space.
package ipcalc

import (
	"fmt"
	"math/bits"
	"strconv"
)

func (ip IPv4) Pretty() [][2]string {
	return [][2]string{
		{"Version", "IPv4"},
		{"Addr/Pref", ip.NiceAddrMask()},
		{"Address", ip.NiceAddr()},
		{"Mask", ip.NiceMask()},
		{"Network", ip.GetNetwork()},
		{"Broadcast", ip.GetBroadcast()},
		{"Host number", ip.GetNumberOfPossibleHosts()},
	}
}

// NiceAddrMask return string in ipv4/mask format
func (ip IPv4) NiceAddrMask() string {
	v := fmt.Sprintf("%v/%v", NiceAddrIPv4(ip.ipAddr), ip.NicePrefix())
	return v
}

// NiceAddr return string with nice formated ipv4 address
func (ip IPv4) NiceAddr() string {
	return NiceAddrIPv4(ip.ipAddr)
}

// NiceMask return string with nice formated ipv4 mask address
func (ip IPv4) NiceMask() string {
	return NiceAddrIPv4(ip.mask)
}

// NicePrefix return string with nice formated decimal mask
func (ip IPv4) NicePrefix() int {
	return bits.OnesCount32(ip.mask)
}

// getNetwork calulate network address
func (ip IPv4) getNetwork() uint32 {
	r := ip.mask & ip.ipAddr
	return r
}

// GetNetwork return string with calcualted network address
func (ip IPv4) GetNetwork() string {
	r := ip.getNetwork()
	return NiceAddrIPv4(r)
}

// getBroadcast calulate broadcast address
func (ip IPv4) getBroadcast() uint32 {
	r := ^ip.mask | ip.ipAddr
	return r
}

// GetBroadcast return string with calcualted broadcast address
func (ip IPv4) GetBroadcast() string {
	r := ip.getBroadcast()
	return NiceAddrIPv4(r)
}

// GetNumberOfPossibleHosts return number of posible hosts in string format
func (ip IPv4) GetNumberOfPossibleHosts() string {
	hostMask := ^ip.mask
	hostBits := bits.OnesCount32(hostMask)
	switch hostBits {
	case 0: // /32
		return "1"
	case 1: // /31
		return "2"
	default:
		return strconv.Itoa((1 << hostBits) - 2)
	}
}
