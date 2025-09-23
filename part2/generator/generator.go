package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

// Square calculates the square of a number
func Square(A float64) float64 {
	Result := A * A
	return Result
}

// RadiansFromDegrees converts degrees to radians
func RadiansFromDegrees(Degrees float64) float64 {
	Result := 0.01745329251994329577 * Degrees
	return Result
}

// ReferenceHaversine calculates the Haversine distance between two points on Earth
// NOTE: EarthRadius is generally expected to be 6372.8
func ReferenceHaversine(X0, Y0, X1, Y1, EarthRadius float64) float64 {
	/* NOTE: This is not meant to be a "good" way to calculate the Haversine distance.
	   Instead, it attempts to follow, as closely as possible, the formula used in the real-world
	   question on which these homework exercises are loosely based.
	*/

	lat1 := Y0
	lat2 := Y1
	lon1 := X0
	lon2 := X1

	dLat := RadiansFromDegrees(lat2 - lat1)
	dLon := RadiansFromDegrees(lon2 - lon1)
	lat1 = RadiansFromDegrees(lat1)
	lat2 = RadiansFromDegrees(lat2)

	a := Square(math.Sin(dLat/2.0)) + math.Cos(lat1)*math.Cos(lat2)*Square(math.Sin(dLon/2))
	c := 2.0 * math.Asin(math.Sqrt(a))

	Result := EarthRadius * c

	return Result
}

// CoordinatePair represents a pair of coordinates (x0, y0, x1, y1)
type CoordinatePair struct {
	X0 float64 `json:"x0"`
	Y0 float64 `json:"y0"`
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
}

// CoordinateData represents the JSON structure with pairs array
type CoordinateData struct {
	Pairs []CoordinatePair `json:"pairs"`
}

// RandomFloatInRange generates a random float64 number within the specified range [min, max)
func RandomFloatInRange(min, max float64, seed int64) float64 {
	NewRandom := rand.New(rand.NewSource(seed))

	return min + NewRandom.Float64()*(max-min)
}

// WriteDistancesToBinaryFile writes distances to a binary file with newlines and total sum
func WriteDistancesToBinaryFile(filename string, distances []float64) (float64, error) {
	file, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var totalSum float64

	for _, distance := range distances {
		totalSum += distance
		err := binary.Write(file, binary.LittleEndian, distance)
		if err != nil {
			return 0, err
		}

		_, err = file.Write([]byte{'\n'})
		if err != nil {
			return 0, err
		}
	}

	err = binary.Write(file, binary.LittleEndian, totalSum)
	if err != nil {
		return 0, err
	}

	return totalSum, nil
}

func main() {
	exampleUsage := "Usage: generator <uniform|cluster> <random seed> <number of coordinates paris to generate>"
	if len(os.Args) != 4 {
		fmt.Println(exampleUsage)
		return
	}

	mode := os.Args[1]
	seed := os.Args[2]
	numberOfCoordinates := os.Args[3]

	if mode != "uniform" && mode != "cluster" {
		fmt.Println(exampleUsage)
		return
	}

	seedInt, err := strconv.Atoi(seed)
	if err != nil {
		fmt.Println(exampleUsage)
		return
	}

	if seedInt <= 0 {
		fmt.Println(exampleUsage)
		return
	}

	numberOfCoordinatesInt, err := strconv.Atoi(numberOfCoordinates)
	if err != nil {
		fmt.Println(exampleUsage)
		return
	}

	if numberOfCoordinatesInt < 0 {
		fmt.Println(exampleUsage)
		return
	}

	inputDir := "input"
	err = os.Mkdir(inputDir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating input directory: %v\n", err)
		return
	}

	// Create coordinate data structure
	coordinateData := CoordinateData{
		Pairs: make([]CoordinatePair, numberOfCoordinatesInt),
	}

	distances := make([]float64, numberOfCoordinatesInt)
	earthRadius := 6372.8 // Standard Earth radius in km

	switch mode {
	case "uniform":
		// Generate uniform random coordinates
		for i := 0; i < numberOfCoordinatesInt; i++ {
			// Use different seeds for each coordinate to ensure randomness
			x0 := RandomFloatInRange(-180.0, 180.0, int64(seedInt+i*2))
			y0 := RandomFloatInRange(-90.0, 90.0, int64(seedInt+i*2+1))
			x1 := RandomFloatInRange(-180.0, 180.0, int64(seedInt+i*2+1000))
			y1 := RandomFloatInRange(-90.0, 90.0, int64(seedInt+i*2+1001))

			coordinateData.Pairs[i] = CoordinatePair{
				X0: x0,
				Y0: y0,
				X1: x1,
				Y1: y1,
			}

			// Calculate Haversine distance for this pair
			distances[i] = ReferenceHaversine(x0, y0, x1, y1, earthRadius)
		}
	case "cluster":
		// Generate clustered coordinates (points closer together)
		// Create a few cluster centers
		numClusters := 3
		clusterCenters := make([]CoordinatePair, numClusters)

		// Generate cluster centers
		for i := 0; i < numClusters; i++ {
			clusterCenters[i] = CoordinatePair{
				X0: RandomFloatInRange(-180.0, 180.0, int64(seedInt+i*10)),
				Y0: RandomFloatInRange(-90.0, 90.0, int64(seedInt+i*10+1)),
				X1: RandomFloatInRange(-180.0, 180.0, int64(seedInt+i*10+2)),
				Y1: RandomFloatInRange(-90.0, 90.0, int64(seedInt+i*10+3)),
			}
		}

		// Generate points around cluster centers
		for i := 0; i < numberOfCoordinatesInt; i++ {
			clusterIndex := i % numClusters
			center := clusterCenters[clusterIndex]

			// Add some random variation around the cluster center (±10 degrees)
			x0 := center.X0 + RandomFloatInRange(-10.0, 10.0, int64(seedInt+i*4))
			y0 := center.Y0 + RandomFloatInRange(-10.0, 10.0, int64(seedInt+i*4+1))
			x1 := center.X1 + RandomFloatInRange(-10.0, 10.0, int64(seedInt+i*4+2))
			y1 := center.Y1 + RandomFloatInRange(-10.0, 10.0, int64(seedInt+i*4+3))

			// Clamp to valid coordinate ranges
			x0 = math.Max(-180.0, math.Min(180.0, x0))
			y0 = math.Max(-90.0, math.Min(90.0, y0))
			x1 = math.Max(-180.0, math.Min(180.0, x1))
			y1 = math.Max(-90.0, math.Min(90.0, y1))

			coordinateData.Pairs[i] = CoordinatePair{
				X0: x0,
				Y0: y0,
				X1: x1,
				Y1: y1,
			}

			distances[i] = ReferenceHaversine(x0, y0, x1, y1, earthRadius)
		}
	}

	// Create JSON file
	jsonData, err := json.MarshalIndent(coordinateData, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Write JSON file
	jsonFilename := fmt.Sprintf("input/coordinates_%s_%d_%d.json", mode, seedInt, numberOfCoordinatesInt)
	err = os.WriteFile(jsonFilename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		return
	}

	binaryFilename := fmt.Sprintf("input/distances_%s_%d_%d", mode, seedInt, numberOfCoordinatesInt)
	totalSum, err := WriteDistancesToBinaryFile(binaryFilename, distances)
	if err != nil {
		fmt.Printf("Error writing binary file: %v\n", err)
		return
	}

	fmt.Printf("Generated %d coordinate pairs in %s mode\n", numberOfCoordinatesInt, mode)
	fmt.Printf("Total sum of all distances: %.6f km\n", totalSum)
	fmt.Printf("JSON file saved as: %s\n", jsonFilename)
	fmt.Printf("Binary file saved as: %s\n", binaryFilename)
}
