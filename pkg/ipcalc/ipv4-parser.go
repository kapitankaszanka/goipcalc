// Copyright (c) 2025 Jan Kowalski
// Licensed under the MIT License. See LICENSE file in the project root for details.

package ipcalc

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseIPv4Prefix parse x.x.x.x/y
func ParseIPv4Prefix(s string) (IPv4, error) {
	var out IPv4

	// CIDR
	addr, pfxStr, ok := strings.Cut(s, "/")
	if !ok {
		p := fmt.Errorf("invalid addr, expected <ipv4>/mask, given: %s", s)
		return out, p
	}

	// Mask
	v, err := parseMask(pfxStr)
	if err != nil {
		return out, err
	}
	out.mask = v

	// Addres
	ip := uint32(0x0)
	parts := strings.Split(addr, ".")
	if len(parts) > 4 {
		return out, fmt.Errorf("invalid addr, to many octets: %s", addr)
	}
	for i, p := range parts {
		v, err := parseOctet(p)
		if err != nil {
			return out, fmt.Errorf("invalid addr: %s, %s", addr, err)
		}
		ip = ip | (v << (24 - (8 * i)))
	}

	out.ipAddr = ip
	return out, nil
}

// convertToByteIPv4 format x.x.x.x to uint32
func parseOctet(o string) (uint32, error) {
	v, err := strconv.ParseUint(o, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid octet: %q", o)
	}

	return uint32(v), nil
}

// parseMask format /x to -> uint32
func parseMask(m string) (uint32, error) {
	v, err := strconv.ParseUint(m, 10, 9)
	if err != nil {
		return 0, fmt.Errorf("invalid prefix: %q", m)
	}
	if v > 32 {
		return 0, fmt.Errorf("invalid mask: %q", m)
	}

	pfx := uint32(v)
	mask := ^uint32(0) << (32 - pfx)
	return mask, nil
}
