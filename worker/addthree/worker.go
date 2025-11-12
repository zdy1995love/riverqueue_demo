package addthree

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riverqueue/river"
)

// Args represents the worker arguments
type Args struct {
	Number int `json:"number"`
}

// Kind returns the task type name, which River uses to identify the worker
func (Args) Kind() string { return "add_three" }

// Worker implements the add-three functionality
type Worker struct {
	river.WorkerDefaults[Args]
	logger *slog.Logger
}

// NewWorker creates a new AddThree Worker
func NewWorker(logger *slog.Logger) *Worker {
	return &Worker{
		logger: logger,
	}
}

// Work executes the actual job
func (w *Worker) Work(ctx context.Context, job *river.Job[Args]) error {
	input := job.Args.Number
	result := input + 3

	w.logger.Info("AddThreeWorker processing",
		"job_id", job.ID,
		"input", input,
		"result", result,
	)

	// Simulate some processing time
	// time.Sleep(time.Second)

	fmt.Printf("✅ AddThree: %d + 3 = %d (chained task completed)\n", input, result)

	return nil
}
