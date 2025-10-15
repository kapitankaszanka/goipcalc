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

func (ip IP) Pretty(format bool) [][2]string {

	return [][2]string{
		{"Addr/Pfx", ip.GetAddrMask()},
		{"Address", ip.NiceAddr()},
		{"Mask", strconv.Itoa(int(ip.pfx))},
		ip.GetNetwork(),
		ip.GetLastAddr(),
		{"Host number", ip.GetHostsNumberStr(format)},
	}
}

// NiceAddr format IP.addr [][uint16 to string ipv4/6 address string
func (ip IP) NiceAddr() string {
	switch len(ip.addr) {
	case 2:
		return fmt.Sprintf("%d.%d.%d.%d",
			byte(ip.addr[0]>>8),
			byte(ip.addr[0]&0x00ff),
			byte(ip.addr[1]>>8),
			byte(ip.addr[1]&0x00ff),
		)
	case 8:
		return fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x",
			ip.addr[0],
			ip.addr[1],
			ip.addr[2],
			ip.addr[3],
			ip.addr[4],
			ip.addr[5],
			ip.addr[6],
			ip.addr[7],
		)
	default:
		return ""
	}
}

func (ip IP) GetAddrMask() string {
	return fmt.Sprintf("%s/%d", ip.NiceAddr(), ip.pfx)
}

// GetNetwork return string with calcualted network address
func (ip IP) GetNetwork() [2]string {
	r := ip
	topic := "Network"

	switch r.pfx {
	case 0:
		for i := range len(r.addr) {
			r.addr[i] = 0x0
		}
		return [2]string{topic, r.NiceAddr()}
	case 32:
		return [2]string{topic, r.NiceAddr()}
	case 128:
		return [2]string{topic, r.NiceAddr()}
	default:
		// to see which hextext may be changed
		hextetNum := int(r.pfx / 16)
		// to flick corect bits
		rem := uint8(16 - (r.pfx % 16))
		for i := range r.addr {
			switch {
			// when the hextet does not need to be changed, omit
			case i < hextetNum:
				continue
			// when the hextet may be changed
			case i == hextetNum:
				if rem == 0 {
					r.addr[i] = 0x0
				} else {
					r.addr[i] = r.addr[i] & r.mask[i]
				}
			default:
				r.addr[i] = 0xFFFF
			}
		}
		return [2]string{topic, r.NiceAddr()}
	}
}

// GetLastAddr return last of ipv6 address
func (ip IP) GetLastAddr() [2]string {
	r := ip
	var topic string

	switch len(r.addr) {
	case 4:
		topic = "Broadcast"
	default:
		topic = "Last address"
	}

	switch r.pfx {
	case 0:
		for i := range len(r.addr) {
			r.addr[i] = 0xFFFF
		}
		return [2]string{topic, r.NiceAddr()}
	case 32:
		return [2]string{topic, r.NiceAddr()}
	case 128:
		return [2]string{topic, r.NiceAddr()}
	default:
		// to see which hextext may be changed
		hextetNum := int(r.pfx / 16)
		// to flick corect bits
		rem := uint8(16 - (r.pfx % 16))
		for i := range r.addr {
			switch {
			// when the hextet does not need to be changed, omit
			case i < hextetNum:
				continue
			// when the hextet may be changed
			case i == hextetNum:
				if rem == 0 {
					r.addr[i] = 0xFFFF
				} else {
					r.addr[i] = r.addr[i] | ^r.mask[i]
				}
			default:
				r.addr[i] = 0xFFFF
			}
		}
		return [2]string{topic, r.NiceAddr()}
	}
}

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
