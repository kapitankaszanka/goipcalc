// Copyright (c) 2025 Mateusz Krupczyński
// Licensed under the MIT License.
// See LICENSE file in the project root for details.

package ipcalc

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseIPv4Prefix parse x.x.x.x/y
func ParseIPv4Prefix(s string) (IP, error) {
	var out IP

	// CIDR
	addr, pfxStr, ok := strings.Cut(s, "/")
	if !ok {
		p := fmt.Errorf("invalid addr, expected <ipv4>/mask, given: %s", s)
		return out, p
	}

	// Mask
	mask, pfx, err := parseMask(pfxStr)
	if err != nil {
		return out, err
	}

	ip, err := parseOctets(addr)
	if err != nil {
		return out, err
	}

	out.addr = ip
	out.mask = mask
	out.pfx = pfx
	return out, nil
}

// parseOctets from string and return array uint16 or error.
func parseOctets(addrStr string) ([]uint16, error) {
	result := make([]uint16, 2)

	// check if string is correct ipv4 address
	parts := strings.Split(addrStr, ".")
	if len(parts) > 4 {
		p := fmt.Errorf("invalid addr, to many octets: %s", addrStr)
		return result, p
	}

	var tmp uint8 = 0
	// parse octets
	for i, v := range parts {
		o, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return []uint16{}, fmt.Errorf("invalid address: %s", addrStr)
		}

		if i%2 == 0 {
			result[tmp] = 0x0 | (uint16(o) << 8)
		} else {
			result[tmp] = result[tmp] | uint16(o)
			tmp += 1
		}
	}

	return result, nil
}

// parseMask valid if mask is corect and return []uint16 with mask and prefix
func parseMask(m string) ([]uint16, uint8, error) {
	mask := make([]uint16, 2)
	v, err := strconv.ParseUint(m, 10, 9)
	if err != nil {
		return mask, 0, fmt.Errorf("invalid prefix: %q", m)
	}
	if v > 32 {
		return mask, 0, fmt.Errorf("invalid mask: %q", m)
	}

	pfx := uint8(v)
	if pfx > 16 {
		mask[0] = 0xFFFF
		mask[1] = 0xFFFF << (16 - pfx)
	} else {
		mask[0] = 0xFFFF << (16 - pfx)
		mask[1] = 0x0000
	}

	return mask, pfx, nil
}
