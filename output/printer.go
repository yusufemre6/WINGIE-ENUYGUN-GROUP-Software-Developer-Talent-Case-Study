// Package output formats and prints scheduling results.
package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"wingie_case/model"
)

// Printer is the interface for rendering a ScheduleResult.
type Printer interface {
	Print(result *model.ScheduleResult)
}

// ConsolePrinter writes a human-readable schedule to an io.Writer.
type ConsolePrinter struct {
	writer io.Writer
}

// NewConsolePrinter creates a printer that writes to stdout.
func NewConsolePrinter() *ConsolePrinter {
	return &ConsolePrinter{writer: os.Stdout}
}

// NewConsolePrinterWithWriter creates a printer that writes to the given writer.
func NewConsolePrinterWithWriter(w io.Writer) *ConsolePrinter {
	return &ConsolePrinter{writer: w}
}

// Print renders the schedule: summary, task table, and execution order.
func (p *ConsolePrinter) Print(result *model.ScheduleResult) {
	w := p.writer
	line := strings.Repeat("=", 60)
	dash := strings.Repeat("-", 60)

	fmt.Fprintln(w)
	fmt.Fprintln(w, line)
	fmt.Fprintf(w, "  Job: %s\n", result.JobName)
	fmt.Fprintln(w, line)

	fmt.Fprintf(w, "  Minimum completion time : %d unit(s)\n", result.MinCompletionTime)
	if len(result.CriticalPath) > 0 {
		fmt.Fprintf(w, "  Critical path           : %s\n", strings.Join(result.CriticalPath, " -> "))
	}

	fmt.Fprintln(w, dash)
	fmt.Fprintln(w, "  Execution Plan (assuming parallel execution):")
	fmt.Fprintln(w, dash)
	fmt.Fprintf(w, "  %-8s %12s %12s %12s\n", "Task", "Start", "Finish", "Duration")
	fmt.Fprintln(w, dash)

	for _, ts := range result.TaskSchedules {
		dur := ts.EarliestFinish - ts.EarliestStart
		fmt.Fprintf(w, "  %-8s %12d %12d %12d\n",
			ts.TaskID, ts.EarliestStart, ts.EarliestFinish, dur)
	}

	fmt.Fprintln(w, dash)
	fmt.Fprintf(w, "  Execution order: [%s]\n", strings.Join(result.ExecutionOrder, ", "))
	fmt.Fprintln(w, line)
}

