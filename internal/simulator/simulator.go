package simulator

import (
	"fmt"

	"github.com/8086-simulator/internal/instruction"
)

type Result struct {
	Text string
}

type Simulator struct {
	Registers     map[string]int16
	registerOrder []string
}

func NewSimulator() *Simulator {
	s := &Simulator{
		Registers: make(map[string]int16),
	}
	s.Init()
	return s
}

func (s *Simulator) Init() {
	s.registerOrder = []string{"ax", "bx", "cx", "dx", "sp", "bp", "si", "di"}
	s.Registers = map[string]int16{
		"ax": 0,
		"bx": 0,
		"cx": 0,
		"dx": 0,
		"sp": 0,
		"bp": 0,
		"si": 0,
		"di": 0,
	}
}

func (s *Simulator) Run(instructions []*instruction.Instruction) ([]*Result, error) {
	results := []*Result{}
	for _, ins := range instructions {
		switch ins.Op {
		case instruction.MOV:
			switch ins.OperandType {
			case instruction.OpTypeImmToReg:
				destPrevVal := s.Registers[ins.DestRegister]
				s.Registers[ins.DestRegister] = ins.Immediate
				results = append(
					results,
					&Result{Text: fmt.Sprintf("%s ; %s:0x%x->0x%x", ins.Text, ins.DestRegister, destPrevVal, ins.Immediate)},
				)
			case instruction.OpTypeRegMemToFromReg:
				destPrevVal := s.Registers[ins.DestRegister]
				sourceVal := s.Registers[ins.SourceRegister]
				s.Registers[ins.DestRegister] = sourceVal
				results = append(
					results,
					&Result{Text: fmt.Sprintf("%s ; %s:0x%x->0x%x", ins.Text, ins.DestRegister, destPrevVal, sourceVal)},
				)
			default:
				return nil, fmt.Errorf("unsupported operand type: %d", ins.OperandType)
			}
		default:
			return nil, fmt.Errorf("unsupported instruction: %s", ins.Op)
		}
	}

	return results, nil
}
