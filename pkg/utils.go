// Package pkg allow for calculate address space.
package pkg

import (
	"bytes"
	"fmt"
	"io"
	"ipcalculator/pkg/ipcalc"
	"text/tabwriter"
)

func NicePrint(w io.Writer, pl []ipcalc.Pretty) error {
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.StripEscape)

	fmt.Fprintln(tw, "------ Start")
	for _, p := range pl {
		items := p.Pretty()
		fmt.Fprintf(tw, "--- %s %s\n", items[0][0], items[0][1])
		for _, kv := range items[1:] {
			fmt.Fprintf(tw, "%s\t : %s\n", kv[0], kv[1])
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
