package ipcalc_test

import (
	"goipcalc/pkg/ipcalc"
	"testing"
)

var testCasesIPv6 = []struct {
	input   string
	expAddr []uint16
	expMask []uint16
	expPfx  uint8
}{
	// valid
	{
		"2001:db8::1/64",
		[]uint16{0x2001, 0x0db8, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0001},
		[]uint16{0xffff, 0xffff, 0xffff, 0xffff, 0x0000, 0x0000, 0x0000, 0x0000},
		uint8(64),
	},
	{
		"::1/128",
		[]uint16{0, 0, 0, 0, 0, 0, 0, 1},
		[]uint16{0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff},
		uint8(128),
	},
	{
		"::/0",
		[]uint16{0, 0, 0, 0, 0, 0, 0, 0},
		[]uint16{0, 0, 0, 0, 0, 0, 0, 0},
		uint8(0),
	},
	{
		"1234:1234:1234:1234:1234:1234:1234:1234/64",
		[]uint16{0x1234, 0x1234, 0x1234, 0x1234, 0x1234, 0x1234, 0x1234, 0x1234},
		[]uint16{0xffff, 0xffff, 0xffff, 0xffff, 0, 0, 0, 0},
		uint8(64),
	},

	// invalid
	{"2001:db8::1/129", nil, nil, 0},             // wrong mask
	{"2001:db8::1", nil, nil, 0},                 // no mask
	{"2001:::1/64", nil, nil, 0},                 // triple :
	{"2001::1234::1/64", nil, nil, 0},            // double ::
	{"2001::1234:192.168.11.24/64", nil, nil, 0}, // no supported
}

func TestParseIPv6Prefix(t *testing.T) {
	for _, tt := range testCasesIPv6 {
		ip, err := ipcalc.ParseIPv6Prefix(tt.input)
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
			t.Errorf("%q prefix got %d, want %d, type: %T", tt.input, ip.Pfx, tt.expPfx, ip.Pfx)
		}
	}

}
