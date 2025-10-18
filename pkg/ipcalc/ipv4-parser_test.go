package ipcalc_test

import (
	"goipcalc/pkg/ipcalc"
	"testing"
)

var testCasesIPv4 = []struct {
	input   string
	expAddr []uint16
	expMask []uint16
	expPfx  uint8
}{
	// valid addresses
	{
		"192.168.0.0/24",
		[]uint16{0xc0a8, 0x0000},
		[]uint16{0xffff, 0xff00},
		uint8(24),
	},
	{
		"10.0.0.1/8",
		[]uint16{0x0a00, 0x0001},
		[]uint16{0xff00, 0x0000},
		uint8(8),
	},
	{
		"127.0.0.1/32",
		[]uint16{0x7f00, 0x0001},
		[]uint16{0xffff, 0xffff},
		uint8(32),
	},
	{
		"0.0.0.0/0",
		[]uint16{0x0000, 0x0000},
		[]uint16{0x0000, 0x0000},
		uint8(0),
	},

	// invalid
	{"256.0.0.1/24", nil, nil, 0},    // wrong octet
	{"192.168.0.1/33", nil, nil, 0},  // wrong mask
	{"192.168.0.1", nil, nil, 0},     // no mask
	{"192.192.168.0.1", nil, nil, 0}, // to long
}

func TestParseIPv4Prefix(t *testing.T) {
	for _, tt := range testCasesIPv4 {
		ip, err := ipcalc.ParseIPv4Prefix(tt.input)
		if tt.expAddr == nil {
			if err == nil {
				t.Errorf("%q expected error, got none", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q unexpected error: %v", tt.input, err)
			continue
		}

		gotAddr := ip.Addr
		gotMask := ip.Mask

		if !EqualU16(gotAddr, tt.expAddr) {
			t.Errorf("%q addr got %v, want %v", tt.input, gotAddr, tt.expAddr)
		}
		if !EqualU16(gotMask, tt.expMask) {
			t.Errorf("%q mask got %v, want %v", tt.input, gotMask, tt.expMask)
		}
		if ip.Pfx != tt.expPfx {
			t.Errorf("%q prefix got %d, want %d", tt.input, ip.Pfx, tt.expPfx)
		}
	}

}
