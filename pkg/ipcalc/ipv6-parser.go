package ipcalc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ParseIPv6Prefix format x:x:x:x:x:x:x:x to [8]uint16
func ParseIPv6Prefix(s string) (IPv6, error) {
	var out IPv6

	// CIDR
	addr, pfxStr, ok := strings.Cut(s, "/")
	if !ok {
		p := fmt.Errorf("invalid addr, expected <ipv6>/mask, given: %s", s)
		return out, p
	}
	if strings.Contains(addr, ".") {
		return out, fmt.Errorf("invalid addr, not allowe ipv4 format: %s", s)
	}
	if strings.Contains(addr, "%") {
		return out, fmt.Errorf("invalid addr, not allowed zones: %s", s)
	}
	pfxU, err := strconv.ParseUint(pfxStr, 10, 8)
	if err != nil || pfxU > 128 {
		return out, fmt.Errorf("invalid addr, prefix to long: %s", s)
	}
	out.mask = uint8(pfxU)

	// Address
	if strings.Count(addr, "::") > 1 {
		return out, fmt.Errorf("invalid addr, multiple '::'")
	}

	if strings.Contains(addr, "::") {
		leftRight := strings.SplitN(addr, "::", 2)
		left, right := leftRight[0], leftRight[1]
		var leftParts, rightParts []string
		if left != "" {
			leftParts = strings.Split(left, ":")
		}
		if right != "" {
			rightParts = strings.Split(right, ":")
		}
		if len(leftParts)+len(rightParts) > 8 {
			return out, fmt.Errorf("invalid addr, to many hextet")
		}

		zeros := 8 - (len(leftParts) + len(rightParts)) // how many 0x0 to add
		idx := 0                                        // index for all loops
		for _, p := range leftParts {
			v, err := parseHextet(p)
			if err != nil {
				return out, err
			}
			out.hextet[idx] = v
			idx++
		}
		for range zeros {
			out.hextet[idx] = 0
			idx++
		}
		for _, p := range rightParts {
			v, err := parseHextet(p)
			if err != nil {
				return out, err
			}
			out.hextet[idx] = v
			idx++
		}
	} else {
		parts := strings.Split(addr, ":")
		if len(parts) != 8 {
			p := fmt.Errorf("invalid addr, expected 8 hextet is %d", len(parts))
			return out, p
		}
		for i, p := range parts {
			v, err := parseHextet(p)
			if err != nil {
				return out, err
			}
			out.hextet[i] = v
		}
	}

	return out, nil
}

func parseHextet(p string) (uint16, error) {
	if len(p) == 0 || len(p) > 4 {
		return 0, fmt.Errorf("invalid hextet: %q", p)
	}
	for _, r := range p {
		if !unicode.Is(unicode.ASCII_Hex_Digit, r) {
			return 0, fmt.Errorf("not hext in hextet: %q", p)
		}
	}
	u, err := strconv.ParseUint(p, 16, 16)
	if err != nil {
		return 0, err
	}
	return uint16(u), nil
}
