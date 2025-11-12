package multiplytwo

import (
	"context"
	"fmt"
	"log/slog"
	"zdy/worker/addthree"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// Args represents the worker arguments
type Args struct {
	Number int `json:"number"`
}

// Kind returns the task type name, which River uses to identify the worker
func (Args) Kind() string { return "multiply_two" }

// RiverClient defines the River client interface
type RiverClient interface {
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// Worker implements the multiply-by-two functionality
type Worker struct {
	river.WorkerDefaults[Args]
	logger      *slog.Logger
	riverClient RiverClient
}

// NewWorker creates a new MultiplyTwo Worker
func NewWorker(logger *slog.Logger, riverClient RiverClient) *Worker {
	return &Worker{
		logger:      logger,
		riverClient: riverClient,
	}
}

// SetClient sets the River client (for chained task support)
func (w *Worker) SetClient(client RiverClient) {
	w.riverClient = client
}

// Work executes the actual job
func (w *Worker) Work(ctx context.Context, job *river.Job[Args]) error {
	input := job.Args.Number
	result := input * 2

	w.logger.Info("MultiplyTwoWorker processing",
		"job_id", job.ID,
		"input", input,
		"result", result,
	)

	// Simulate some processing time
	// time.Sleep(time.Second)

	fmt.Printf("✅ MultiplyTwo: %d × 2 = %d\n", input, result)

	// Chain next task: pass result to AddThree
	if w.riverClient != nil {
		_, err := w.riverClient.Insert(ctx, addthree.Args{Number: result}, nil)
		if err != nil {
			w.logger.Error("Failed to insert AddThree task", "error", err, "result", result)
			return fmt.Errorf("failed to insert chained task: %w", err)
		}
		w.logger.Info("🔗 Chained task inserted", "next_task", "add_three", "input", result)
	}

	return nil
}
