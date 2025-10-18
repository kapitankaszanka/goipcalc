package ipcalc_test

import "fmt"

func EqualU16(a, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func PrintAddress(ip []uint16) string {
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
