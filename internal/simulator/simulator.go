package simulator

import "fmt"

type Simulator struct{}

func NewSimulator() *Simulator {
	return &Simulator{}
}

func (s *Simulator) Run() {
	fmt.Println("Simulator running")
}
