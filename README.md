# Job Scheduler — Critical Path Calculator

**Wingie EnUygun Group — Software Developer Case Study**

A CLI tool that computes the minimum completion time for a job consisting of dependent tasks, using the **Critical Path Method (CPM)**.

> For a detailed explanation of the algorithm, assumptions, and design decisions, see **[SOLUTION.md](SOLUTION.md)**.

## Project Structure

```
.
├── main.go                  # Entry point, dependency injection
├── model/
│   ├── task.go              # Task entity
│   ├── job.go               # Job entity
│   └── schedule_result.go   # Scheduling output model
├── input/
│   └── reader.go            # Reader interface + CLIReader
├── validator/
│   └── validator.go         # Validator interface + GraphValidator
├── scheduler/
│   └── scheduler.go         # Scheduler interface + CriticalPathScheduler
├── output/
│   └── printer.go           # Printer interface + ConsolePrinter
├── Dockerfile               # Multi-stage build
├── .dockerignore
├── go.mod
├── SOLUTION.md              # Algorithm explanation and reasoning
└── README.md
```

## Running

### With Docker (recommended)

```bash
docker build -t job-scheduler .
docker run -it job-scheduler
```

### With Go (locally)

```bash
go run .   # start the application
```

## Example

```
Enter job name (e.g. J): J
How many tasks?: 6

--- Task 1 ---
Task ID (e.g. A): A
Duration for task 'A' (positive integer): 3
Dependencies for task 'A' (comma-separated, or leave empty):

--- Task 2 ---
Task ID (e.g. A): B
Duration for task 'B' (positive integer): 2
Dependencies for task 'B' (comma-separated, or leave empty):

--- Task 3 ---
Task ID (e.g. A): C
Duration for task 'C' (positive integer): 4
Dependencies for task 'C' (comma-separated, or leave empty):

--- Task 4 ---
Task ID (e.g. A): D
Duration for task 'D' (positive integer): 5
Dependencies for task 'D' (comma-separated, or leave empty): A

--- Task 5 ---
Task ID (e.g. A): E
Duration for task 'E' (positive integer): 2
Dependencies for task 'E' (comma-separated, or leave empty): B,C

--- Task 6 ---
Task ID (e.g. A): F
Duration for task 'F' (positive integer): 3
Dependencies for task 'F' (comma-separated, or leave empty): D,E
```

### Output

```
============================================================
  Job: J
============================================================
  Minimum completion time : 11 unit(s)
  Critical path           : A -> D -> F
------------------------------------------------------------
  Execution Plan (assuming parallel execution):
------------------------------------------------------------
  Task            Start       Finish     Duration
------------------------------------------------------------
  A                   0            3            3
  B                   0            2            2
  C                   0            4            4
  D                   3            8            5
  E                   4            6            2
  F                   8           11            3
------------------------------------------------------------
  Execution order: [A, B, C, D, E, F]
============================================================
```
