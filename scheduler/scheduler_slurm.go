package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/squarefactory/benchmark-api/utils"
)

const QosName = "mining"

type Slurm struct {
	executor  Executor
	adminUser string
}

func NewSlurm(
	executor Executor,
	adminUser string,
) *Slurm {
	return &Slurm{
		executor:  executor,
		adminUser: adminUser,
	}
}

// CancelJob kills a job using scancel command.
func (s *Slurm) CancelJob(ctx context.Context, req *CancelRequest) error {
	cmd := fmt.Sprintf("scancel --name=%s --me", req.Name)
	_, err := s.executor.ExecAs(ctx, req.User, cmd)
	if err != nil {
		log.Printf("cancel failed: %s", err)
	}
	return err
}

// Submit a sbatch definition script to the SLURM controller using the sbatch command.
func (s *Slurm) Submit(ctx context.Context, req *SubmitRequest) (string, error) {
	eof := utils.GenerateRandomString(10)

	cmd := fmt.Sprintf(`sbatch \
  --job-name=%s \
  --qos=%s \
  --output=/tmp/benchmark-%%j_%%a.log \
  --parsable << '%s'
%s
%s`,
		req.Name,
		QosName,
		eof,
		req.Body,
		eof,
	)
	out, err := s.executor.ExecAs(ctx, req.User, cmd)
	if err != nil {
		log.Printf("submit failed: %s", err)
		return strings.TrimSpace(strings.TrimRight(string(out), "\n")), err
	}

	return strings.TrimSpace(strings.TrimRight(string(out), "\n")), nil
}

// HealthCheck runs squeue to check if the queue is running
func (s *Slurm) HealthCheck(ctx context.Context) error {
	_, err := s.executor.ExecAs(ctx, s.adminUser, "squeue")
	if err != nil {
		log.Printf("healthcheck failed: %s", err)
	}
	return err
}

// FindRunningJobByName find a running job using squeue.
func (s *Slurm) FindRunningJobByName(
	ctx context.Context,
	req *FindRunningJobByNameRequest,
) (int, error) {
	cmd := fmt.Sprintf("squeue --name %s -O ArrayJobId:256 --noheader", req.Name)
	out, err := s.executor.ExecAs(ctx, req.User, cmd)
	if err != nil {
		log.Printf("FindRunningJobByName failed: %s", err)
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		log.Println("no jobs currently running")
		return 0, errors.New("no running jobs found")
	}

	jobID, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		log.Printf("Failed to parse JobId: %s", err)
		return 0, err
	}

	return jobID, nil
}

func (s *Slurm) FindMemPerNode(ctx context.Context) (int, error) {
	cmd := "scontrol show nodes | grep Mem | sed -E 's|.*RealMemory=([^,]*)|\\1|g' | awk '{print $1}'"
	out, err := s.executor.ExecAs(ctx, s.adminUser, cmd)
	if err != nil {
		log.Printf("FindMemPerNode failed: %s", err)
		return 0, err
	}

	out = strings.TrimSpace(string(out))
	lines := strings.Split(out, "\n")
	mem, err := strconv.Atoi(lines[0])
	if err != nil {
		log.Printf("Failed to convert %q to integer: %s", lines[0], err)
		return 0, err
	}

	return mem, nil
}

func (s *Slurm) FindGPUPerNode(ctx context.Context) (int, error) {
	cmd := "scontrol show nodes | grep CfgTRES | sed -E 's|.*gres/gpu=([^,]*)|\\1|g'"
	out, err := s.executor.ExecAs(ctx, s.adminUser, cmd)
	if err != nil {
		log.Printf("FindGPUPerNode failed: %s", err)
		return 0, err
	}

	out = strings.TrimSpace(string(out))
	lines := strings.Split(out, "\n")
	gpu, err := strconv.Atoi(lines[0])
	if err != nil {
		log.Printf("Failed to convert %q to integer: %s", lines[0], err)
		return 0, err
	}

	return gpu, nil
}
