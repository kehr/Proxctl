package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type Printer struct {
	Out    io.Writer
	Err    io.Writer
	JSON   bool
	Color  bool
	Quiet  bool
	Failed int
	Warned int
	Passed int
}

func (p *Printer) Section(name string) {
	if p.JSON || p.Quiet {
		return
	}
	fmt.Fprintf(p.Out, "\n============================================================\n%s\n============================================================\n", name)
}

func (p *Printer) Info(format string, args ...any) {
	if p.JSON || p.Quiet {
		return
	}
	fmt.Fprintf(p.Out, "[INFO] "+format+"\n", args...)
}

func (p *Printer) Pass(format string, args ...any) {
	p.Passed++
	if p.JSON || p.Quiet {
		return
	}
	fmt.Fprintf(p.Out, "[PASS] "+format+"\n", args...)
}

func (p *Printer) Warn(format string, args ...any) {
	p.Warned++
	if p.JSON || p.Quiet {
		return
	}
	fmt.Fprintf(p.Out, "[WARN] "+format+"\n", args...)
}

func (p *Printer) Fail(format string, args ...any) {
	p.Failed++
	if p.JSON || p.Quiet {
		return
	}
	fmt.Fprintf(p.Out, "[FAIL] "+format+"\n", args...)
}

func (p *Printer) JSONValue(v any) error {
	enc := json.NewEncoder(p.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
