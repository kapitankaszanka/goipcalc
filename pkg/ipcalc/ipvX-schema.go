// Package ipcalc allow for calculate address space.
package ipcalc

type IPv4 struct {
	ipAddr uint32
	mask   uint32
}

type IPv6 struct {
	hextet [8]uint16
	mask   uint8
}

type Pretty interface {
	Pretty() [][2]string
}
