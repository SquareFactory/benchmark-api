package benchmark

import (
	"context"

	"github.com/squarefactory/benchmark-api/scheduler"
)

type SlurmScheduler interface {
	Submit(ctx context.Context, req *scheduler.SubmitRequest) (string, error)
	CancelJob(ctx context.Context, req *scheduler.CancelRequest) error
	HealthCheck(ctx context.Context) error
	FindRunningJobByName(
		ctx context.Context,
		req *scheduler.FindRunningJobByNameRequest,
	) (int, error)
	FindMemPerNode(ctx context.Context) (int, error)
	FindGPUPerNode(ctx context.Context) (int, error)
	FindCPUPerNode(ctx context.Context) (int, error)
}

type Benchmark struct {
	dat         DATParams
	sbatch      SBATCHParams
	SlurmClient SlurmScheduler
}

type BenchmarkFile struct {
	DatFile    string
	SbatchFile string
}

type DATParams struct {
	ProblemSize int
	P           int
	Q           int
}

type SBATCHParams struct {
	Node          int
	NtasksPerNode int
	GpusPerNode   int
	CpusPerTasks  int
}
