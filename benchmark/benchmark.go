package benchmark

import (
	"bytes"
	"context"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/squarefactory/benchmark-api/scheduler"
)

const (
	User        = "root"
	GBtoMB      = 1000
	JobName     = "HPL-Benchmark"
	DatFilePath = "hpl.dat"
)

var benchmarkMemoryUsePercentage = []float64{
	0.75,
	0.76,
	0.77,
	0.78,
	0.79,
	0.80,
	0.81,
	0.82,
	0.83,
	0.84,
}

func NewBenchmark(
	dat DATParams,
	sbatch SBATCHParams,
	slurm SlurmScheduler,
) *Benchmark {
	return &Benchmark{
		Dat:         dat,
		Sbatch:      sbatch,
		SlurmClient: slurm,
	}
}

func (b *Benchmark) Run(ctx context.Context, files *BenchmarkFile) error {

	if err := os.WriteFile(DatFilePath, []byte(files.DatFile), 0644); err != nil {
		return err
	}

	out, err := b.SlurmClient.Submit(ctx, &scheduler.SubmitRequest{
		Name: JobName,
		User: User,
		Body: files.SbatchFile,
	})
	if err != nil {
		log.Printf("Failed to run benchmark: %s", err)
		return err
	}

	log.Printf("Successfully started benchmark: %s", out)
	return nil
}

func (b *Benchmark) GenerateFiles(ctx context.Context) (BenchmarkFile, error) {

	DatFile, err := b.GenerateDAT()
	if err != nil {
		log.Printf("Failed to generate DAT file: %s", err)
		return BenchmarkFile{}, err

	}
	var SbatchFile string

	if b.Sbatch.Node != 1 {
		SbatchFile, err = b.GenerateMultiNodeSBATCH()
		if err != nil {
			log.Printf("Failed to generate SBATCH file: %s", err)
			return BenchmarkFile{}, err
		}
	} else {
		SbatchFile, err = b.GenerateSingleNodeSBATCH()
		if err != nil {
			log.Printf("Failed to generate SBATCH file: %s", err)
			return BenchmarkFile{}, err
		}
	}

	return BenchmarkFile{
		SbatchFile: SbatchFile,
		DatFile:    DatFile,
	}, nil
}

func (b *Benchmark) GenerateDAT() (string, error) {

	DATTmpl := template.Must(template.New("jobTemplate").Parse(DatTmpl))
	var DatFile bytes.Buffer
	if err := DATTmpl.Execute(&DatFile, struct {
		NProblemSize int
		ProblemSize  string
		NBlockSize   int
		BlockSize    string
		P            int
		Q            int
	}{
		NProblemSize: b.Dat.NProblemSize,
		ProblemSize:  b.Dat.ProblemSize,
		NBlockSize:   b.Dat.NBlockSize,
		BlockSize:    b.Dat.BlockSize,
		P:            b.Dat.P,
		Q:            b.Dat.Q,
	}); err != nil {
		log.Printf("dat templating failed: %s", err)
		return "", err
	}

	return DatFile.String(), nil
}

func (b *Benchmark) GenerateMultiNodeSBATCH() (string, error) {

	SbatchTmpl := template.Must(template.New("jobTemplate").Parse(MultiNodeTmpl))
	var SbatchFile bytes.Buffer
	if err := SbatchTmpl.Execute(&SbatchFile, struct {
		ContainerPath string
		Workspace     string
		Node          int
		CpusPerTasks  int
		GpusPerNode   int
		NtasksPerNode int
		GpuAffinity   string
		CpuAffinity   string
	}{
		ContainerPath: b.Sbatch.ContainerPath,
		Workspace:     b.Sbatch.Workspace,
		Node:          b.Sbatch.Node,
		CpusPerTasks:  b.Sbatch.CpusPerTasks,
		GpusPerNode:   b.Sbatch.GpusPerNode,
		NtasksPerNode: b.Sbatch.NtasksPerNode,
		GpuAffinity:   b.Sbatch.GpuAffinity,
		CpuAffinity:   b.Sbatch.CpuAffinity,
	}); err != nil {
		log.Printf("sbatch templating failed: %s", err)
		return "", err
	}

	return SbatchFile.String(), nil

}

func (b *Benchmark) GenerateSingleNodeSBATCH() (string, error) {
	SbatchTmpl := template.Must(template.New("jobTemplate").Parse(SingleNodeTmpl))
	var SbatchFile bytes.Buffer
	if err := SbatchTmpl.Execute(&SbatchFile, struct {
		ContainerPath string
		Workspace     string
		Node          int
		CpusPerTasks  int
		GpusPerNode   int
		NtasksPerNode int
		GpuAffinity   string
		CpuAffinity   string
	}{
		ContainerPath: b.Sbatch.ContainerPath,
		Workspace:     b.Sbatch.Workspace,
		Node:          b.Sbatch.Node,
		CpusPerTasks:  b.Sbatch.CpusPerTasks,
		GpusPerNode:   b.Sbatch.GpusPerNode,
		NtasksPerNode: b.Sbatch.NtasksPerNode,
		GpuAffinity:   b.Sbatch.GpuAffinity,
		CpuAffinity:   b.Sbatch.CpuAffinity,
	}); err != nil {
		log.Printf("sbatch templating failed: %s", err)
		return "", err
	}

	return SbatchFile.String(), nil
}

func (b *Benchmark) CalculateBenchmarkParams(ctx context.Context) error {
	if err := b.CalculateDATParams(ctx); err != nil {
		log.Printf("Failed to calculate dat params: %s", err)
		return err
	}

	if err := b.CalculateSBATCHParams(ctx); err != nil {
		log.Printf("Failed to calculate sbatch params: %s", err)
		return err
	}

	return nil
}

func (b *Benchmark) CalculateDATParams(ctx context.Context) error {
	if err := b.CalculateProblemSize(ctx); err != nil {
		return err
	}

	if err := b.CalculateProcessGrid(ctx); err != nil {
		return err
	}

	b.Dat.NBlockSize = 10
	b.Dat.BlockSize = "64 128 224 256 384 512 640 768 896 1024"

	return nil
}

func (b *Benchmark) CalculateSBATCHParams(ctx context.Context) error {
	b.Sbatch.NtasksPerNode = b.Dat.P * b.Dat.Q / b.Sbatch.Node
	CpusPerNode, err := b.SlurmClient.FindCPUPerNode(ctx)
	if err != nil {
		return err
	}
	b.Sbatch.CpusPerTasks = CpusPerNode / b.Sbatch.NtasksPerNode

	b.Sbatch.GpusPerNode, err = b.SlurmClient.FindGPUPerNode(ctx)
	if err != nil {
		return err
	}

	if err := b.CalculateAffinity(ctx); err != nil {
		return err
	}

	return nil
}

// Calculates the optimal values of P and Q based on the number of GPUs available per nodes
func (b *Benchmark) CalculateProcessGrid(ctx context.Context) error {

	numGPUs, err := b.SlurmClient.FindGPUPerNode(ctx)
	if err != nil {
		log.Printf("failed to calculate gpus per node : %s", err)
		return err
	}
	totalGPUS := numGPUs * b.Sbatch.Node

	if totalGPUS == 1 {
		b.Dat.P = 1
		b.Dat.Q = 1
		return nil
	}

	sqrttotalGPUS := int(math.Sqrt(float64(totalGPUS)))

	for i := sqrttotalGPUS; i > 0; i-- {
		if totalGPUS%i == 0 && i != 1 {
			b.Dat.P = i
			b.Dat.Q = totalGPUS / i
			return nil
		}
	}

	b.Dat.P = 2
	b.Dat.Q = totalGPUS
	return nil // If no other valid P is found, default to 2
}

// Calculates the problem size from the ram available
func (b *Benchmark) CalculateProblemSize(ctx context.Context) error {

	mem, err := b.SlurmClient.FindMemPerNode(ctx)
	if err != nil {
		log.Printf("failed to calculate problem size: %s", err)
		return err
	}

	b.Dat.NProblemSize = len(benchmarkMemoryUsePercentage)
	for _, values := range benchmarkMemoryUsePercentage {
		problemSize := int(
			math.Sqrt(float64(mem*b.Sbatch.Node)/8)*values,
		) * GBtoMB

		b.Dat.ProblemSize += strconv.Itoa(problemSize) + " "
	}

	return nil
}

func (b *Benchmark) CalculateAffinity(ctx context.Context) error {

	out, err := b.SlurmClient.FindCPUAffinity(ctx)
	if err != nil {
		log.Printf("failed to calculate cpu affinity: %s", err)
		return err
	}
	gpusPerTasks := b.Sbatch.NtasksPerNode / b.Sbatch.GpusPerNode

	pattern := `(\d+)\s+(\d+-\d+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(out, -1)

	// Process the matches
	var cpuAffinityValues []string
	var gpuAffinityValues []string

	for _, match := range matches {
		cpuAffinity := match[2]
		gpuAffinity := match[1]

		// Generate the CPU affinity value by repeating the CPU affinity value for the given number of tasks per GPU
		cpu := strings.Repeat(cpuAffinity+":", gpusPerTasks)
		cpuAffinityValues = append(
			cpuAffinityValues,
			cpu[:len(cpu)-1],
		) // Remove the trailing colon

		// Generate the GPU affinity value by repeating the CPU affinity value for the given number of tasks per GPU
		gpu := strings.Repeat(gpuAffinity+":", gpusPerTasks)
		gpuAffinityValues = append(
			gpuAffinityValues,
			gpu[:len(gpu)-1],
		) // Remove the trailing colon

	}

	// Join the GPU affinity values with a colon
	b.Sbatch.CpuAffinity = strings.Join(cpuAffinityValues, ":")
	b.Sbatch.GpuAffinity = strings.Join(gpuAffinityValues, ":")

	return nil
}
