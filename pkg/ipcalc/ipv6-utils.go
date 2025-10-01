// Package ipcalc allow for calculate address space.
package ipcalc

import (
	"fmt"
	"strconv"
)

func (ip IPv6) Pretty() [][2]string {
	return [][2]string{
		{"Version", "IPv6"},
		{"Addr/Pref", ip.NiceAddrMask()},
		{"Address", ip.NiceAddr()},
		{"Mask", ip.NiceMask()},
		{"Network", ip.GetNetwork()},
		{"Last address", ip.GetLastAddr()},
		{"Host number", ip.GetNumberOfPossibleHosts()},
	}
}

// NiceAddrMask return nicly formated ipv6 address with mask
func (ip IPv6) NiceAddrMask() string {
	v := fmt.Sprintf("%v/%v", NiceAddrIPv6(ip), ip.mask)
	return v
}

// NiceAddr return nicly formated ipv6 address
func (ip IPv6) NiceAddr() string {
	return NiceAddrIPv6(ip)
}

// NiceMask return mask in string format
func (ip IPv6) NiceMask() string {
	return strconv.Itoa(int(ip.mask))
}

// GetLastAddr return last of ipv6 address
func (ip IPv6) GetLastAddr() string {
	r := IPv6{}
	r.mask = ip.mask
	hextetNum := int(r.mask / 16)
	rem := uint8(r.mask % 16)

	if r.mask == 0 {
		r.hextet = [8]uint16{
			0xFFFF,
			0xFFFF,
			0xFFFF,
			0xFFFF,
			0xFFFF,
			0xFFFF,
			0xFFFF,
			0xFFFF,
		}
		return r.NiceAddrMask()
	}

	for i, h := range ip.hextet {
		switch {
		case i < hextetNum:
			r.hextet[i] = h
		case i == hextetNum:
			if rem == 0 {
				r.hextet[i] = 0xFFFF
			} else {
				m1 := uint16(0xFFFF << (16 - rem)) // mask for leading ones
				m2 := uint16(0xFFFF >> rem)        // mask for remains
				r.hextet[i] = (m1 & h) | m2
			}
		default:
			r.hextet[i] = 0xFFFF
		}
	}
	return r.NiceAddrMask()
}

// GetNetwork return firs of ipv6 address
func (ip IPv6) GetNetwork() string {
	r := IPv6{}
	r.mask = ip.mask
	hextetNum := int(ip.mask / 16)
	m := uint16(0xFFFF << (16 - ip.mask%16))

	for i, h := range ip.hextet {
		switch {
		case i < hextetNum:
			r.hextet[i] = h
		case i == hextetNum:
			r.hextet[i] = m & h
		default:
			r.hextet[i] = 0x0
		}
	}
	return r.NiceAddrMask()
}

// GetNumberOfPossibleHosts always return "To many to bother..."
func (ip IPv6) GetNumberOfPossibleHosts() string {
	return "To many to bother..."
}
