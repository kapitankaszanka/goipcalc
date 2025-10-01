// Copyright (c) 2025 Mateusz Krupczyński
// Licensed under the MIT License.
// See LICENSE file in the project root for details.

package ipcalc

import "fmt"

// NiceAddrIPv4 format ipv4 struct to ipv4 string addres
func NiceAddrIPv4(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24),
		byte(ip>>16),
		byte(ip>>8),
		byte(ip),
	)
}

// NiceAddrIPv6 format ipv6 struct to string ipv6 address
func NiceAddrIPv6(ip IPv6) string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x",
		ip.hextet[0],
		ip.hextet[1],
		ip.hextet[2],
		ip.hextet[3],
		ip.hextet[4],
		ip.hextet[5],
		ip.hextet[6],
		ip.hextet[7],
	)
}
