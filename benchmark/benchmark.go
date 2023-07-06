package benchmark

import (
	"bytes"
	"context"
	"log"
	"math"
	"os"
	"text/template"

	"github.com/squarefactory/benchmark-api/scheduler"
)

const (
	user                         = "root"
	benchmarkMemoryUsePercentage = 0.75
	JobName                      = "HPL-Benchmark"
	DatFilePath                  = "hpl.dat"
)

func (b *Benchmark) Run(ctx context.Context, files *BenchmarkFile) error {

	if err := os.WriteFile(DatFilePath, []byte(files.DatFile), 0644); err != nil {
		return err
	}

	out, err := b.SlurmClient.Submit(ctx, &scheduler.SubmitRequest{
		Name: JobName,
		User: user,
		Body: files.SbatchFile,
	})
	if err != nil {
		log.Printf("Failed to run benchmark: %s", err)
		return err
	}

	log.Printf("Successfully started benchmark: %s", out)
	return nil
}

func (b *Benchmark) GenerateFiles(ctx context.Context, node int) (BenchmarkFile, error) {

	if err := b.CalculateBenchmarkParams(ctx); err != nil {
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
		ProblemSize: b.Dat.ProblemSize,
		P:           b.Dat.P,
		Q:           b.Dat.Q,
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
		CpusPerTasks:  b.Sbatch.CpusPerTasks,
		GpusPerNode:   b.Sbatch.GpusPerNode,
		NtasksPerNode: b.Sbatch.NtasksPerNode,
	}); err != nil {
		log.Printf("sbatch templating failed: %s", err)
		return "", err
	}

	return SbatchFile.String(), nil

}

// Returns a benchmark and all its parameters
func (b *Benchmark) CalculateBenchmarkParams(ctx context.Context) error {
	if err := b.CalculateProblemSize(ctx); err != nil {
		return err
	}

	if err := b.CalculateProcessGrid(ctx); err != nil {
		return err
	}

	b.Sbatch.NtasksPerNode = b.Dat.P * b.Dat.Q
	CpusPerNode, err := b.SlurmClient.FindCPUPerNode(ctx)
	if err != nil {
		return err
	}
	b.Sbatch.CpusPerTasks = CpusPerNode / b.Sbatch.NtasksPerNode

	b.Sbatch.GpusPerNode, err = b.SlurmClient.FindGPUPerNode(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Calculates the optimal values of P and Q based on the number of GPUs available per snodes
func (b *Benchmark) CalculateProcessGrid(ctx context.Context) error {

	numGPUs, err := b.SlurmClient.FindGPUPerNode(ctx)
	if err != nil {
		log.Printf("failed to calculate gpus per node : %s", err)
		return err
	}

	if numGPUs == 1 {
		b.Dat.P = 1
		b.Dat.Q = 1
		return nil
	}

	sqrtNumGPUs := int(math.Sqrt(float64(numGPUs)))

	for i := sqrtNumGPUs; i > 0; i-- {
		if numGPUs%i == 0 && i != 1 {
			b.Dat.P = i
			b.Dat.Q = numGPUs / i
			return nil
		}
	}

	b.Dat.P = 2
	b.Dat.Q = numGPUs
	return nil // If no other valid P is found, default to 2
}

// Calculates the problem size from the ram available
func (b *Benchmark) CalculateProblemSize(ctx context.Context) error {
	mem, err := b.SlurmClient.FindMemPerNode(ctx)
	if err != nil {
		log.Printf("failed to calculate problem size: %s", err)
		return err
	}

	problemSize := math.Sqrt(float64(mem)/8) * benchmarkMemoryUsePercentage

	b.Dat.ProblemSize = int(problemSize)
	return nil
}
