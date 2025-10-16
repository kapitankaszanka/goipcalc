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
		fmt.Fprintln(os.Stderr, "Usage: goipcalc [OPTIONS] [ADDR/PLEN]")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  goipcalc -d 10.0.0.1/24")
		fmt.Fprintln(os.Stderr, "  goipcalc 2001:db8::1/64 192.168.10.11/28")
		fmt.Fprintln(os.Stderr, "Options:")
		fmt.Fprintln(os.Stderr, "  [ADDR/PLEN] address/prefix lenght, can be multiple")
		flag.PrintDefaults()
	}

	do := flag.Bool("d", false, "IPv4 address to calculate")
	jo := flag.Bool("j", false, "json output")
	ji := flag.Bool("json-indent", false, "change json output to indentation")

	flag.Parse()

	ips := flag.Args()
	if len(ips) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no address provided.")
		flag.Usage()
		os.Exit(1)
	}

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

	// output
	if len(pretty) > 0 {
		if *jo {
			err := pkg.NicePrintJSON(
				os.Stdout,
				pretty,
				*do,
				*ji,
			)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
			os.Exit(0)
		} else {
			err := pkg.NicePrintCLI(os.Stdout, pretty, *do)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
}
