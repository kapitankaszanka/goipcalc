// Copyright (c) 2025 Mateusz Krupczyński
// Licensed under the MIT License.
// See LICENSE file in the project root for details.

// Package pkg allows calculation and formatted output of IP address space.
package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goipcalc/pkg/ipcalc"
	"io"
	"text/tabwriter"
)

// NicePrintCLI formats and writes the IP address calculation results
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
func NicePrintCLI(w io.Writer, ipList []ipcalc.IP, detail bool) error {
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.StripEscape)

	for _, p := range ipList {
		items := p.Pretty(detail, true)
		fmt.Fprintf(tw, "---\n")
		for _, kv := range items {
			fmt.Fprintf(tw, "%s:\t%s\n", kv[0], kv[1])
		}
	}

	// flush tabwriter to buffor
	if err := tw.Flush(); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())

	return err
}

// IPOut represents a structured version of IP address calculation
// results. This type is used for stable JSON encoding output.
type IPOut struct {
	FullAddress string `json:"full_address"`
	Network     string `json:"network,omitempty"`
	Broadcast   string `json:"broadcast,omitempty"`
	LastAddress string `json:"last_address,omitempty"`
	Address     string `json:"address,omitempty"`
	Mask        string `json:"mask,omitempty"`
	MaskAddress string `json:"mask_address,omitempty"`
	HostsNumber string `json:"hosts_number,omitempty"`
}

// NicePrintJSON encodes IP address calculation results as JSON
// and writes them to the provided writer.
//
// If `indent` is true, the JSON is pretty-printed with indentation.
// The `detail` flag controls whether additional fields (mask, hosts, etc.)
// are included in the output.
func NicePrintJSON(w io.Writer, ipList []ipcalc.IP, detail, indent bool) error {
	var buf bytes.Buffer
	out := make([]IPOut, 0, len(ipList))

	for _, ip := range ipList {
		var o IPOut
		for _, kv := range ip.Pretty(detail, false) {
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
				o.Mask = kv[1]
			case "Mask address":
				o.MaskAddress = kv[1]
			case "Hosts number":
				o.HostsNumber = kv[1]
			}
		}
		out = append(out, o)
	}

	enc := json.NewEncoder(&buf)
	if indent {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(out); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// PrintErrors writes a list of error messages to the given writer.
// It is used to output errors in a simple, human-readable format.
//
// Each error is written on its own line with a prefixed header.
func PrintErrors(w io.Writer, errs []string) error {
	var buf bytes.Buffer
	buf.WriteString("------ Error\n")
	for _, err := range errs {
		buf.WriteString(err)
	}
	_, err := w.Write(buf.Bytes())
	return err
}
