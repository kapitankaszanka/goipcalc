// Copyright (c) 2025 Mateusz Krupczyński
// Licensed under the MIT License.
// See LICENSE file in the project root for details.

// Package pkg allow for calculate address space.
package pkg

import (
	"bytes"
	"fmt"
	"goipcalc/pkg/ipcalc"
	"io"
	"text/tabwriter"
)

func NicePrint(w io.Writer, ipList []ipcalc.IP) error {
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.StripEscape)

	for _, p := range ipList {
		items := p.Pretty(true)
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

func PrintErrors(w io.Writer, errs []string) error {
	var buf bytes.Buffer
	buf.WriteString("------ Error")
	for _, err := range errs {
		buf.WriteString(err)
	}
	_, err := w.Write(buf.Bytes())
	return err
}
