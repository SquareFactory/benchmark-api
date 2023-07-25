package resultparser

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const CsvFile = "benchmark.csv"

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
	output, err := os.Create(CsvFile)
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

	log.Printf("Data has been successfully written to %s", CsvFile)
	return nil
}

func FindMaxGflopsRow(inputFile string) ([]string, error) {
	csvFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening the CSV file:", err)
		return nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV records:", err)
		return nil, err
	}

	var maxGflops float64 = -1
	var maxGflopsRow []string

	for _, row := range records {
		gflops, err := strconv.ParseFloat(row[5], 64) // Gflops is in the 6th column (index 5)
		if err != nil {
			fmt.Println("Error converting Gflops to float:", err)
			continue
		}

		if gflops > maxGflops {
			maxGflops = gflops
			maxGflopsRow = row
		}
	}

	return maxGflopsRow, nil
}
