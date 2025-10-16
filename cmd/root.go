// Package cmd handle command line for ipcalc
package cmd

import (
	"flag"
	"fmt"
	"goipcalc/pkg/ipcalc"
	"goipcalc/pkg/output"
	"os"
	"strings"
)

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

	detail := flag.Bool("d", false, "IPv4 address to calculate")
	jsonOut := flag.Bool("j", false, "json output")
	jsonIndent := flag.Bool("json-indent", false, "change json output to indentation")

	flag.Parse()

	ips := flag.Args()
	if len(ips) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no address provided.")
		flag.Usage()
		os.Exit(1)
	}

	objList := make([]ipcalc.IP, 0, len(ips))
	var errors []string
	if len(ips) > 0 {
		for _, v := range ips {
			if strings.Contains(v, ":") {
				obj, err := ipcalc.ParseIPv6Prefix(v)
				if err != nil {
					errors = append(
						errors,
						fmt.Sprintf("skip %q: %v\n", v, err),
					)
					continue
				}
				objList = append(objList, obj)
			} else {
				obj, err := ipcalc.ParseIPv4Prefix(v)
				if err != nil {
					errors = append(
						errors,
						fmt.Sprintf("skip %q: %v\n", v, err),
					)
					continue
				}
				objList = append(objList, obj)
			}
		}
	}

	status, err := output.PrintOutput(*jsonOut, *jsonIndent, *detail, errors, objList)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(status)

}
