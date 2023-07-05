package benchmark

type Benchmark struct {
	dat    DATParams
	sbatch SBATCHParams
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
