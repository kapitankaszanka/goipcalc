// Copyright (c) 2025 Mateusz Krupczyński
// Licensed under the MIT License.
// See LICENSE file in the project root for details.

package main

import (
	"flag"
	"fmt"
	"goipcalc/pkg"
	"goipcalc/pkg/ipcalc"
	"os"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: goipcalc [-ip4 ADDR/PLEN]... [-ip6 ADDR/PLEN]...")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  goipcalc -[ip | ip4] 10.0.0.1/24 -ip4 192.168.1.10/24")
		fmt.Fprintln(os.Stderr, "  goipcalc -ip6 2001:db8::1/64")
		flag.PrintDefaults()
	}

	var ip4 stringSlice
	var ip6 stringSlice

	flag.Var(&ip4, "ip", "IPv4 address to calculate")
	flag.Var(&ip4, "ip4", "IPv4 address to calculate")
	flag.Var(&ip6, "ip6", "IPv6 address to calculate")

	flag.Parse()

	pretty := make([]ipcalc.Pretty, 0, (len(ip6) + len(ip4)))
	var errors []string
	if len(ip4) > 0 {
		for _, v := range ip4 {
			obj, err := ipcalc.ParseIPv4Prefix(v)
			if err != nil {
				errors = append(errors, fmt.Sprintf("skip %q: %v\n", v, err))
				continue
			}
			pretty = append(pretty, obj)
		}
	}

	if len(ip6) > 0 {
		for _, v := range ip6 {
			obj, err := ipcalc.ParseIPv6Prefix(v)
			if err != nil {
				errors = append(errors, fmt.Sprintf("skip %q: %v\n", v, err))
				continue
			}
			pretty = append(pretty, obj)
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
		fmt.Println("------ End")
		os.Exit(0)
	}

}
