package main

import (
	"encoding/binary"
	"fmt"
	"os"

	harvestine "github.com/8086-simulator/part2/harvestine"
	jsonparser "github.com/8086-simulator/part2/json-parser"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input_file> <results_file>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	resultsPath := os.Args[2]

	input, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error reading input file: ", err)
		os.Exit(1)
	}
	results, totalSum, err := readDistancesFromBinaryFile(resultsPath)
	if err != nil {
		fmt.Println("Error reading results file: ", err)
		os.Exit(1)
	}

	lexer := jsonparser.NewLexer(string(input))
	tokens := lexer.Tokenize()
	parser := jsonparser.NewParser(tokens)
	ast := parser.Parse()

	diff := 0.0
	pairs := ast.Children[0].Val.Children
	harvestineSum := 0.0

	for i := range pairs {
		distance := harvestine.ReferenceHaversine(
			pairs[i].Children[0].Val.Value.(float64), pairs[i].Children[1].Val.Value.(float64),
			pairs[i].Children[2].Val.Value.(float64), pairs[i].Children[3].Val.Value.(float64),
			6372.8,
		)
		harvestineSum += distance
		diff += distance - results[i]
	}

	fmt.Println("Input size: ", len(input))
	fmt.Println("Pair count: ", len(ast.Children[0].Val.Children))
	fmt.Println("Harvestine sum: ", harvestineSum)
	fmt.Println()
	fmt.Println("Reference sum: ", totalSum)
	fmt.Println("Difference: ", diff)
}

// readDistancesFromBinaryFile reads distances and total sum from a binary file
func readDistancesFromBinaryFile(filename string) ([]float64, float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	// Get file size to determine how many distances we have
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	// Each float64 is 8 bytes, and we have one extra float64 for the total sum
	// So: (fileSize - 8) / 8 = number of distances
	fileSize := fileInfo.Size()
	if fileSize < 8 {
		return nil, 0, fmt.Errorf("file too small to contain valid data")
	}

	numDistances := (fileSize - 8) / 8
	distances := make([]float64, numDistances)

	// Read all distances
	for i := int64(0); i < numDistances; i++ {
		var distance float64
		err := binary.Read(file, binary.LittleEndian, &distance)
		if err != nil {
			return nil, 0, fmt.Errorf("error reading distance %d: %v", i, err)
		}
		distances[i] = distance
	}

	// Read the total sum (last float64 in the file)
	var totalSum float64
	err = binary.Read(file, binary.LittleEndian, &totalSum)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading total sum: %v", err)
	}

	return distances, totalSum, nil
}
