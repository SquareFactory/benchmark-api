package resultparser

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strings"
)

var outputFile = "benchmark.csv"

func WriteResultsToCSV(inputFile string) error {

	// Open the input file
	input, err := os.Open(inputFile)
	if err != nil {
		log.Printf("Failed to open input file: %s", err)
		return err
	}
	defer input.Close()

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

	// Write the CSV header
	header := []string{
		"Identifier",
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

	// Read the input file line by line
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line starts with 'HPL_AI'
		if strings.HasPrefix(line, "HPL_AI") {
			// Split the line into fields
			fields := strings.Fields(line)

			// Extract the required values
			problemsize := fields[2]
			nb := fields[3]
			p := fields[4]
			q := fields[5]
			time := fields[6]
			gflops := fields[7]
			refine := fields[8]
			iter := fields[9]
			gflops_wrefinement := fields[10]

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

	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read input file: %s", err)
		return err
	}

	log.Printf("Data has been successfully written to %s", outputFile)
	return nil
}
