package resultparser

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func WriteResultsToCSV(resultFile, csvFile string) error {

	// Read the input file contents
	inputBytes, err := os.ReadFile(resultFile)
	if err != nil {
		log.Printf("Failed to read input file: %s", err)
		return err
	}

	// Convert the input file contents to string
	inputData := string(inputBytes)

	// Split the input data into lines
	lines := strings.Split(inputData, "\n")

	// Create the output file
	output, err := os.Create(csvFile)
	if err != nil {
		log.Printf("Failed to create output file: %s", err)
		return err
	}
	defer output.Close()

	if err := WriteDataAsCsvRecord(output, lines); err != nil {
		log.Printf("failed to write data in csv format: %s", err)
		return err
	}

	log.Printf("Data has been successfully written to %s", csvFile)
	return nil
}

func AppendResultsToCsv(resultFile, csvFile string) error {
	// Read the input file contents
	inputBytes, err := os.ReadFile(resultFile)
	if err != nil {
		log.Printf("Failed to read input file: %s", err)
		return err
	}

	inputData := string(inputBytes)
	lines := strings.Split(inputData, "\n")

	output, err := os.OpenFile(csvFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open CSV file: %s", err)
		return err
	}
	defer output.Close()

	if err := WriteDataAsCsvRecord(output, lines); err != nil {
		log.Printf("failed to write data in csv format: %s", err)
		return err
	}

	log.Printf("Data has been successfully appended to %s", csvFile)
	return nil
}

func WriteDataAsCsvRecord(file *os.File, lines []string) error {

	writer := csv.NewWriter(file)
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
	err := writer.Write(header)
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

			err := writer.Write(record)
			if err != nil {
				log.Printf("Failed to write CSV record: %s", err)
				return err
			}
		}
	}
	return nil

}

func FindMaxGflopsRow(csvFile string) ([]string, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println("Error opening the CSV file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
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
