// Copyright (c) 2025 Mateusz Krupczyński
// Licensed under the MIT License.
// See LICENSE file in the project root for details.

// Package output allows calculation and formatted
// output of IP address space.
package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goipcalc/pkg/ipcalc"
	"math/big"
	"os"
	"strconv"
	"text/tabwriter"
)

// JSONOut represent structured version of complete IPOut list and errors
// This type is used for stable JSON output.
type JSONOut struct {
	Results []IPOut  `json:"results"`
	Errors  []string `json:"errors,omitempty"`
}

// IPOut represents a structured version of IP address calculation
// results. This type is used for stable JSON encoding output.
type IPOut struct {
	FullAddress string   `json:"full_address"`
	Network     string   `json:"network,omitempty"`
	Broadcast   string   `json:"broadcast,omitempty"`
	LastAddress string   `json:"last_address,omitempty"`
	Address     string   `json:"address,omitempty"`
	Mask        int      `json:"mask,omitempty"`
	MaskAddress string   `json:"mask_address,omitempty"`
	HostsNumber *big.Int `json:"hosts_number,omitempty"`
}

// nicePrintCLI formats and writes the IP address calculation results
// in a human-readable CLI table format.
//
// It uses the tabwriter to align output in columns. If `detail` is true,
// additional information such as mask and host count is included.
//
// Example output:
// ---
// Full address:  10.0.1.1/24
// Network:       10.0.1.0
// Broadcast:     10.0.1.255
// Address:       10.0.1.1
// Mask:          24
// Mask address:  255.255.255.0
// Hosts number:  256
func nicePrintCLI(b *bytes.Buffer, ipList []ipcalc.IP, detail bool) error {
	tw := tabwriter.NewWriter(b, 0, 0, 2, ' ', tabwriter.StripEscape)

	for _, p := range ipList {
		items := p.Pretty(detail, true)
		fmt.Fprintf(tw, "---\n")
		for _, kv := range items {
			fmt.Fprintf(tw, "%s:\t%s\n", kv[0], kv[1])
		}
	}

	// flush tabwriter to buffer
	return tw.Flush()
}

// nicePrintJSON encodes IP address calculation results as JSON
// and writes them to the provided writer.
//
// If `indent` is true, the JSON is pretty-printed with indentation.
// The `detail` flag controls whether additional fields (mask, hosts, etc.)
// are included in the output.
func nicePrintJSON(buf *bytes.Buffer, ips []ipcalc.IP, errs []string, d, i bool) error {
	out := JSONOut{
		Results: make([]IPOut, 0, len(ips)),
		Errors:  errs,
	}

	for _, ip := range ips {
		var o IPOut
		for _, kv := range ip.Pretty(d, false) {
			switch kv[0] {
			case "Full address":
				o.FullAddress = kv[1]
			case "Network":
				o.Network = kv[1]
			case "Broadcast":
				o.Broadcast = kv[1]
			case "Last address":
				o.LastAddress = kv[1]
			case "Address":
				o.Address = kv[1]
			case "Mask":
				if v, err := strconv.Atoi(kv[1]); err == nil {
					o.Mask = v
				}
			case "Mask address":
				o.MaskAddress = kv[1]
			case "Hosts number":
				if v, ok := new(big.Int).SetString(kv[1], 10); ok {
					o.HostsNumber = v
				}
			}
		}
		out.Results = append(out.Results, o)
	}

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if i {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(out); err != nil {
		return err
	}

	return nil
}

// errorsCLI writes a list of error messages to the given writer.
// It is used to output errors in a simple, human-readable format.
//
// Each error is written on its own line with a prefixed header.
func errorsCLI(buf *bytes.Buffer, errs []string) {
	buf.WriteString("--- Error\n")
	for _, e := range errs {
		buf.WriteString(e)
	}
}

// PrintOutput renders results and errors to stdout/stderr and returns an exit status.
// Behavior:
//   - Exit 0 if at least one input produced results (even if some failed).
//   - Exit 1 if none succeeded and there were errors.
//   - JSON mode still writes human-readable errors to stderr; stdout is JSON
//     (or "[]\n" if there are no results).
func PrintOutput(
	jsonOut, jsonIndent, d bool,
	errList []string,
	ipList []ipcalc.IP,
) (int, error) {
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	// status devicion
	hadResults := len(ipList) > 0
	hadErrors := len(errList) > 0
	status := 0
	if !hadResults && hadErrors {
		status = 1
	}
	// print errors (stderr) in CLI mode
	if !jsonOut && hadErrors {
		errorsCLI(errBuf, errList)
	}

	// handel corect output
	if jsonOut {
		if hadResults {
			if err := nicePrintJSON(outBuf, ipList, errList, d, jsonIndent); err != nil {
				return 1, err
			}
		} else {
			outBuf.WriteString("[]\n")
		}
	} else {
		if err := nicePrintCLI(outBuf, ipList, d); err != nil {
			return 1, err
		}
	}

	if !jsonOut {
		_, err := os.Stderr.Write(errBuf.Bytes())
		if err != nil {
			return 1, err
		}
	}
	_, err := os.Stdout.Write(outBuf.Bytes())
	if err != nil {
		return 1, err
	}

	return status, nil
}
