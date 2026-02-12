package main

import (
	"fmt"
	"os"

	"wingie_case/input"
	"wingie_case/output"
	"wingie_case/scheduler"
	"wingie_case/validator"
)

// App wires together all components via interfaces.
type App struct {
	reader    input.Reader
	validator validator.Validator
	scheduler scheduler.Scheduler
	printer   output.Printer
}

// NewApp creates an App with the given dependencies.
func NewApp(
	reader input.Reader,
	val validator.Validator,
	sched scheduler.Scheduler,
	printer output.Printer,
) *App {
	return &App{
		reader:    reader,
		validator: val,
		scheduler: sched,
		printer:   printer,
	}
}

// Run executes the full pipeline: read → validate → schedule → print.
func (a *App) Run() error {
	printWelcome()

	in, err := a.reader.ReadJob()
	if err != nil {
		return fmt.Errorf("input error: %w", err)
	}

	if err := a.validator.Validate(in.Job); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	result, err := a.scheduler.Schedule(in.Job, in.Workers)
	if err != nil {
		return fmt.Errorf("scheduling error: %w", err)
	}

	a.printer.Print(result)
	return nil
}

func printWelcome() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║       Job Scheduler - Critical Path Calculator          ║")
	fmt.Println("║       Wingie EnUygun Group - Case Study                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func main() {
	app := NewApp(
		input.NewCLIReader(os.Stdin),
		validator.NewGraphValidator(),
		scheduler.NewWorkerScheduler(),
		output.NewConsolePrinter(),
	)

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
		os.Exit(1)
	}
}
