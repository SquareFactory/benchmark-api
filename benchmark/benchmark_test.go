package benchmark_test

import (
	"context"
	"testing"

	"github.com/squarefactory/benchmark-api/benchmark"
	"github.com/squarefactory/benchmark-api/mocks"
	"github.com/squarefactory/benchmark-api/scheduler"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	JobName = "HPL-Benchmark"
	admin   = "root"
)

type ServiceTestSuite struct {
	suite.Suite
	scheduler *mocks.Scheduler
	impl      *benchmark.Benchmark
}

func (suite *ServiceTestSuite) BeforeTest(suiteName, testName string) {
	suite.scheduler = mocks.NewScheduler(suite.T())
	suite.impl = &benchmark.Benchmark{
		SlurmClient: suite.scheduler,
	}
}

func (suite *ServiceTestSuite) TestRun() {

	// Arrange
	files := benchmark.BenchmarkFile{
		DatFile:    "testdatfile",
		SbatchFile: "testsbatchfile",
	}

	expectedSubmitRequest := &scheduler.SubmitRequest{
		Name: JobName,
		User: admin,
		Body: "testsbatchfile",
	}

	suite.scheduler.On(
		"Submit",
		mock.Anything,
		expectedSubmitRequest,
	).Return("test submit response", nil)

	// Act
	err := suite.impl.Run(context.Background(), &files)

	// Assert
	suite.NoError(err)
	suite.scheduler.AssertExpectations(suite.T())
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, &ServiceTestSuite{})
}
