package main

import (
	"log"
	"os"

	"github.com/8086-simulator/internal/decoder"
	"github.com/8086-simulator/internal/simulator"
)

const (
	ExecMode = "exec"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		log.Fatal("No file provided")
	}

	content, err := os.ReadFile(argsWithoutProg[0])
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	dec := decoder.NewDecoder()
	_, err = dec.Decode(content)
	if err != nil {
		log.Fatalf("Error decoding data: %v", err)
	}

	if len(argsWithoutProg) > 1 && argsWithoutProg[1] == ExecMode {
		sim := simulator.NewSimulator()
		sim.Run()
	}
}
