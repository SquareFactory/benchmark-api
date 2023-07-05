package benchmark

import (
	"bytes"
	"context"
	"log"
	"math"
	"os"
	"text/template"

	"github.com/squarefactory/benchmark-api/executor"
	"github.com/squarefactory/benchmark-api/scheduler"
)

const (
	user                         = "root"
	benchmarkMemoryUsePercentage = 0.75
	JobName                      = "HPL-Benchmark"
	DatFilePath                  = "hpl.dat"
)

func (b *BenchmarkFile) Run(ctx context.Context) error {
	slurm := scheduler.NewSlurm(&executor.Shell{}, user)

	if err := os.WriteFile(DatFilePath, []byte(b.DatFile), 0644); err != nil {
		return err
	}

	out, err := slurm.Submit(ctx, &scheduler.SubmitRequest{
		Name: JobName,
		User: user,
		Body: b.SbatchFile,
	})
	if err != nil {
		log.Printf("Failed to run benchmark: %s", err)
		return err
	}

	log.Printf("Successfully started benchmark: %s", out)
	return nil
}

func GenerateFiles(ctx context.Context, node int) (BenchmarkFile, error) {
	slurm := scheduler.NewSlurm(&executor.Shell{}, user)

	b, err := CalculateBenchmarkParams(slurm, ctx)
	if err != nil {
		log.Printf("Failed to generate benchmark parameters: %s", err)
		return BenchmarkFile{}, err
	}

	DatFile, err := b.GenerateDAT()
	if err != nil {
		log.Printf("Failed to generate DAT file: %s", err)
		return BenchmarkFile{}, err
	}
	SbatchFile, err := b.GenerateSBATCH(node)
	if err != nil {
		log.Printf("Failed to generate SBATCH file: %s", err)
		return BenchmarkFile{}, err
	}

	return BenchmarkFile{
		SbatchFile: SbatchFile,
		DatFile:    DatFile,
	}, nil
}

func (b *Benchmark) GenerateDAT() (string, error) {

	// Templating gpu mining job
	DATTmpl := template.Must(template.New("jobTemplate").Parse(DatTmpl))
	var DatFile bytes.Buffer
	if err := DATTmpl.Execute(&DatFile, struct {
		ProblemSize int
		P           int
		Q           int
	}{
		ProblemSize: b.dat.ProblemSize,
		P:           b.dat.P,
		Q:           b.dat.Q,
	}); err != nil {
		log.Printf("dat templating failed: %s", err)
		return "", err
	}

	return DatFile.String(), nil
}

func (b *Benchmark) GenerateSBATCH(node int) (string, error) {

	// Templating gpu mining job
	SbatchTmpl := template.Must(template.New("jobTemplate").Parse(SbatchTmpl))
	var SbatchFile bytes.Buffer
	if err := SbatchTmpl.Execute(&SbatchFile, struct {
		Node          int
		CpusPerTasks  int
		GpusPerNode   int
		NtasksPerNode int
	}{
		Node:          node,
		CpusPerTasks:  b.sbatch.CpusPerTasks,
		GpusPerNode:   b.sbatch.GpusPerNode,
		NtasksPerNode: b.sbatch.NtasksPerNode,
	}); err != nil {
		log.Printf("sbatch templating failed: %s", err)
		return "", err
	}

	return SbatchFile.String(), nil

}

// Returns a benchmark and all its parameters
func CalculateBenchmarkParams(slurm *scheduler.Slurm, ctx context.Context) (*Benchmark, error) {

	ProblemSize, err := CalculateProblemSize(slurm, ctx)
	if err != nil {
		return nil, err
	}
	P, Q, err := CalculateProcessGrid(slurm, ctx)
	if err != nil {
		return nil, err
	}

	NtasksPerNode := P * Q
	CpusPerNode, err := slurm.FindCPUPerNode(ctx)
	if err != nil {
		return nil, err
	}
	CpusPerTasks := CpusPerNode / NtasksPerNode

	GpusPerNode, err := slurm.FindGPUPerNode(ctx)
	if err != nil {
		return nil, err
	}

	b := Benchmark{
		dat: DATParams{
			ProblemSize: ProblemSize,
			P:           P,
			Q:           Q,
		},
		sbatch: SBATCHParams{
			Node:          1,
			NtasksPerNode: NtasksPerNode,
			CpusPerTasks:  CpusPerTasks,
			GpusPerNode:   GpusPerNode,
		},
	}
	return &b, nil
}

// Calculates the optimal values of P and Q based on the number of GPUs available per snodes
func CalculateProcessGrid(slurm *scheduler.Slurm, ctx context.Context) (int, int, error) {

	numGPUs, err := slurm.FindGPUPerNode(ctx)
	if err != nil {
		log.Printf("failed to calculate gpus per node : %s", err)
		return 0, 0, err
	}

	if numGPUs == 1 {
		return 1, 1, nil
	}

	sqrtNumGPUs := int(math.Sqrt(float64(numGPUs)))

	for i := sqrtNumGPUs; i > 0; i-- {
		if numGPUs%i == 0 && i != 1 {
			return i, numGPUs / i, nil
		}
	}

	return 2, numGPUs, nil // If no other valid P is found, default to 2
}

// Calculates the problem size from the ram available
func CalculateProblemSize(slurm *scheduler.Slurm, ctx context.Context) (int, error) {
	mem, err := slurm.FindMemPerNode(ctx)
	if err != nil {
		log.Printf("failed to calculate problem size: %s", err)
		return 0, err
	}

	problemSize := math.Sqrt(float64(mem)/8) * benchmarkMemoryUsePercentage

	return int(problemSize), nil
}
