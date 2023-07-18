package resultparser

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

var outputFile = "benchmark.csv"

func WriteResultsToCSV(inputFile string) error {

	// Read the input file contents
	inputBytes, err := os.ReadFile(inputFile)
	if err != nil {
		log.Printf("Failed to read input file: %s", err)
		return err
	}

	// Convert the input file contents to string
	inputData := string(inputBytes)

	// Split the input data into lines
	lines := strings.Split(inputData, "\n")

	// Create the output file
	output, err := os.Create(outputFile)
	if err != nil {
		log.Printf("Failed to create output file: %s", err)
		return err
	}
	defer output.Close()

	// Create a CSV writer
	writer := csv.NewWriter(output)
	defer writer.Flush()

	header := []string{
		"ProblemSize",
		"NB",
		"P",
		"Q",
		"Time",
		"Gflops",
		"Refine",
		"Iter",
		"Gflops_wrefinement",
	}
	err = writer.Write(header)
	if err != nil {
		log.Printf("Failed to write CSV header: %s", err)
		return err
	}

	// Process each line and extract the required values
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "HPL_AI") {
			// Split the line into fields
			fields := strings.Fields(line)

			// Extract the required values
			problemsize := fields[1]
			nb := fields[2]
			p := fields[3]
			q := fields[4]
			time := fields[5]
			gflops := fields[6]
			refine := fields[7]
			iter := fields[8]
			gflops_wrefinement := fields[9]

			// Write the extracted values to the CSV file
			record := []string{
				problemsize,
				nb,
				p,
				q,
				time,
				gflops,
				refine,
				iter,
				gflops_wrefinement,
			}

			err = writer.Write(record)
			if err != nil {
				log.Printf("Failed to write CSV record: %s", err)
				return err
			}
		}
	}

	log.Printf("Data has been successfully written to %s", outputFile)
	return nil
}
