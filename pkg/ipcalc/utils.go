// Copyright (c) 2025 Mateusz Krupczyński
// Licensed under the MIT License.
// See LICENSE file in the project root for details.

// Package ipcalc
package ipcalc

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type IP struct {
	addr []uint16
	mask []uint16
	pfx  uint8
}

func (ip IP) Pretty(detail bool, pretty bool) [][2]string {
	result := [][2]string{
		{"Full address", ip.GetAddrMask()},
		ip.GetFirstAddress(),
		ip.GetLastAddr(),
	}

	if detail {
		tmp := [][2]string{
			{"Address", NiceAddr(ip.addr)},
			{"Mask", strconv.Itoa(int(ip.pfx))},
			{"Mask address", NiceAddr(ip.mask)},
			{"Hosts number", ip.GetHostsNumberStr(pretty)},
		}
		result = append(result, tmp...)
	}
	return result
}

// NiceAddr format IP.addr [][uint16 to string ipv4/6 address string
func NiceAddr(ip []uint16) string {
	switch len(ip) {
	case 2:
		return fmt.Sprintf("%d.%d.%d.%d",
			byte(ip[0]>>8),
			byte(ip[0]&0x00ff),
			byte(ip[1]>>8),
			byte(ip[1]&0x00ff),
		)
	case 8:
		return fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x",
			ip[0],
			ip[1],
			ip[2],
			ip[3],
			ip[4],
			ip[5],
			ip[6],
			ip[7],
		)
	default:
		return ""
	}
}

func (ip IP) GetAddrMask() string {
	return fmt.Sprintf("%s/%d", NiceAddr(ip.addr), ip.pfx)
}

// GetFirstAddress return string with calcualted network address
func (ip IP) GetFirstAddress() [2]string {
	r := append([]uint16(nil), ip.addr...)
	m := ip.mask
	p := ip.pfx
	topic := "Network"

	switch p {
	case 0:
		for i := range r {
			r[i] = 0x0
		}
		return [2]string{topic, NiceAddr(r)}
	case 32:
		return [2]string{topic, NiceAddr(r)}
	case 128:
		return [2]string{topic, NiceAddr(r)}
	default:
		// to see which hextext may be changed
		hextetNum := int(p / 16)
		// to flick corect bits
		rem := uint8(16 - (p % 16))
		for i := range r {
			switch {
			// when the hextet does not need to be changed, omit
			case i < hextetNum:
				continue
			// when the hextet may be changed
			case i == hextetNum:
				if rem == 0 {
					r[i] = 0x0
				} else {
					r[i] = r[i] & m[i]
				}
			default:
				r[i] = 0x0000
			}
		}
		return [2]string{topic, NiceAddr(r)}
	}
}

// GetLastAddr return last of ipv6 address
func (ip IP) GetLastAddr() [2]string {
	r := append([]uint16(nil), ip.addr...)
	m := ip.mask
	p := ip.pfx
	var topic string

	switch len(r) {
	case 2:
		topic = "Broadcast"
	default:
		topic = "Last address"
	}

	switch p {
	case 0:
		for i := range r {
			r[i] = 0xFFFF
		}
		return [2]string{topic, NiceAddr(r)}
	case 32:
		return [2]string{topic, NiceAddr(r)}
	case 128:
		return [2]string{topic, NiceAddr(r)}
	default:
		// to see which hextext may be changed
		hextetNum := int(p / 16)
		// to flick corect bits
		rem := uint8(16 - (p % 16))
		for i := range r {
			switch {
			// when the hextet does not need to be changed, omit
			case i < hextetNum:
				continue
			// when the hextet may be changed
			case i == hextetNum:
				if rem == 0 {
					r[i] = 0xFFFF
				} else {
					r[i] = r[i] | ^m[i]
				}
			default:
				r[i] = 0xFFFF
			}
		}
		return [2]string{topic, NiceAddr(r)}
	}
}

// formatBigIntWithSpaces sperate big int value on space sparate string
func formatBigIntWithSpaces(n *big.Int) string {
	s := n.String()
	// Insert spaces every 3 digits from the right
	var parts []string
	for len(s) > 3 {
		parts = append(parts, s[len(s)-3:])
		s = s[:len(s)-3]
	}
	parts = append(parts, s)

	// Reverse and join
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, " ")
}

func (ip IP) GetHostsNumberStr(format bool) string {
	var totalBits uint
	if len(ip.addr) == 2 { // IPv4
		totalBits = 32
	} else { // IPv6
		totalBits = 128
	}

	mask := uint(ip.pfx)
	hostBits := totalBits - mask

	// Special case for IPv4 /31 networks
	if totalBits == 32 && mask == 31 {
		return "2"
	}

	// Use big.Int for 2^hostBits
	result := new(big.Int).Lsh(big.NewInt(1), hostBits)

	if format {
		return formatBigIntWithSpaces(result)
	} else {
		return result.String()
	}
}
