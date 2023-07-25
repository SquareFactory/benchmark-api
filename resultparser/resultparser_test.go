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

	cleanData := `141000,128,2,2,25.05,7.461e+04,6.89523,2,5.851e+04
	141000,224,2,2,26.48,7.057e+04,6.98278,2,5.584e+04
	141000,256,2,2,29.61,6.312e+04,6.92718,2,5.115e+04
	141000,384,2,2,33.69,5.547e+04,7.00677,2,4.592e+04`

	err := os.WriteFile(tempInputFile, []byte(cleanData), 0644)
	require.NoError(t, err)

	expected := []string([]string{"141000", "128", "2", "2", "25.05", "7.461e+04", "6.89523", "2", "5.851e+04"})

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
