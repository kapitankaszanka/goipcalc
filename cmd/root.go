// Package cmd handle command line for ipcalc
package cmd

import (
	"flag"
	"fmt"
	"goipcalc/pkg"
	"goipcalc/pkg/ipcalc"
	"os"
	"strings"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func RootCMD() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: goipcalc [-ip4 ADDR/PLEN]... [-ip6 ADDR/PLEN]...")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  goipcalc -[ip | ip4] 10.0.0.1/24 -ip4 192.168.1.10/24")
		fmt.Fprintln(os.Stderr, "  goipcalc -ip6 2001:db8::1/64")
		flag.PrintDefaults()
	}

	var ips stringSlice

	flag.Var(&ips, "ip", "IPv4 address to calculate")
	flag.Var(&ips, "ip4", "IPv4 address to calculate")
	flag.Var(&ips, "ip6", "IPv6 address to calculate")

	flag.Parse()

	pretty := make([]ipcalc.IP, 0, len(ips))
	var errors []string
	if len(ips) > 0 {
		for _, v := range ips {
			if strings.Contains(v, ":") {
				obj, err := ipcalc.ParseIPv6Prefix(v)
				if err != nil {
					errors = append(errors, fmt.Sprintf("skip %q: %v\n", v, err))
					continue
				}
				pretty = append(pretty, obj)
			} else {
				obj, err := ipcalc.ParseIPv4Prefix(v)
				if err != nil {
					errors = append(errors, fmt.Sprintf("skip %q: %v\n", v, err))
					continue
				}
				pretty = append(pretty, obj)
			}
		}
	}

	// errors
	if len(errors) > 0 {
		pkg.PrintErrors(os.Stderr, errors)
		if len(pretty) == 0 {
			os.Exit(2)
		}
	}

	if len(pretty) > 0 {
		pkg.NicePrint(os.Stdout, pretty)
		os.Exit(0)
	}
}
