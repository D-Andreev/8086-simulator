package simulator

import "testing"

func TestSimulator(t *testing.T) {
	simulator := NewSimulator()
	simulator.Run()
}
