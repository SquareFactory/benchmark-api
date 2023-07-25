package resultparser_test

import (
	"os"
	"testing"

	"github.com/squarefactory/benchmark-api/resultparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteResultsToCSV(t *testing.T) {

	tempInputFile := "/tmp/benchmark.log"
	defer os.Remove(tempInputFile)

	cleanData := `HPL_AI WRC01 1 1 1 1 0.001 10.0 1 1 9.5
	HPL_AI WRC01 2 2 2 2 0.002 20.0 1 1 19.0`

	err := os.WriteFile(tempInputFile, []byte(cleanData), 0644)
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Positive test",
			input:   "/tmp/benchmark.log",
			wantErr: false,
		},

		{
			name:    "File does not exist",
			input:   "/tmp/non_existing_file.txt",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := resultparser.WriteResultsToCSV(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("WriteResultsToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindMaxGflopsRow(t *testing.T) {
	tempInputFile := "/tmp/benchmark.log"
	defer os.Remove(tempInputFile)

	cleanData := `ProblemSize,NB,P,Q,Time,Gflops,Refine,Iter,Gflops_wrefinement
95000,64,2,2,29.67,1.927e+04,5.71402,2,1.616e+04
95000,128,2,2,15.67,3.647e+04,5.72738,2,2.671e+04
95000,224,2,2,19.00,3.009e+04,5.74521,2,2.310e+04
95000,256,2,2,17.30,3.304e+04,5.71711,2,2.483e+04
95000,384,2,2,14.75,3.876e+04,5.77248,2,2.785e+04
95000,512,2,2,14.93,3.828e+04,5.76942,2,2.761e+04`

	err := os.WriteFile(tempInputFile, []byte(cleanData), 0644)
	require.NoError(t, err)

	expected := []string{"95000", "384", "2", "2", "14.75", "3.876e+04", "5.77248", "2", "2.785e+04"}

	tests := []struct {
		name      string
		inputFile string
		wantErr   bool
	}{
		{
			name:      "Valid test",
			inputFile: tempInputFile,
			wantErr:   false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resultparser.FindMaxGflopsRow(tt.inputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMaxGflopsRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, expected, got)
		})
	}
}
