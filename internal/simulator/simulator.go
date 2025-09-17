package simulator

import (
	"fmt"

	"github.com/8086-simulator/internal/instruction"
)

type Result struct {
	Text string
}

type Simulator struct {
	Registers     map[string]int
	registerOrder []string
	flags         map[string]bool
	flagOrder     []string
}

func NewSimulator() *Simulator {
	s := &Simulator{
		Registers: make(map[string]int),
	}
	s.Init()
	return s
}

func (s *Simulator) Init() {
	s.registerOrder = []string{"ax", "bx", "cx", "dx", "sp", "bp", "si", "di"}
	s.Registers = map[string]int{
		"ax": 0,
		"bx": 0,
		"cx": 0,
		"dx": 0,
		"sp": 0,
		"bp": 0,
		"si": 0,
		"di": 0,
	}
	s.flags = map[string]bool{
		"Z": false,
		"S": false,
	}
	s.flagOrder = []string{"Z", "S"}
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
		case instruction.ADD, instruction.SUB, instruction.CMP:
			switch ins.OperandType {
			case instruction.OpTypeImmToReg:
				destPrevVal := s.Registers[ins.DestRegister]
				s.Registers[ins.DestRegister] = ins.Immediate
				flagsPrevVal := s.printFlags()
				s.flags["Z"] = s.Registers[ins.DestRegister] == 0
				s.flags["S"] = s.Registers[ins.DestRegister] < 0
				flagsNewVal := s.printFlags()
				if flagsPrevVal != flagsNewVal {
					results = append(
						results,
						&Result{Text: fmt.Sprintf("%s ; %s:0x%x->0x%x %s", ins.Text, ins.DestRegister, destPrevVal, ins.Immediate, flagsNewVal)},
					)
					break
				}
				results = append(
					results,
					&Result{Text: fmt.Sprintf("%s ; %s:0x%x->0x%x", ins.Text, ins.DestRegister, destPrevVal, ins.Immediate)},
				)
			case instruction.OpTypeRegMemToFromReg:
				destPrevVal := s.Registers[ins.DestRegister]
				sourceVal := s.Registers[ins.SourceRegister]
				s.Registers[ins.DestRegister] = sourceVal
				flagsPrevVal := s.printFlags()
				s.flags["Z"] = s.Registers[ins.DestRegister] == 0
				s.flags["S"] = s.Registers[ins.DestRegister] < 0
				flagsNewVal := s.printFlags()
				if flagsPrevVal != flagsNewVal {
					results = append(
						results,
						&Result{Text: fmt.Sprintf("%s ; %s:0x%x->0x%x %s", ins.Text, ins.DestRegister, destPrevVal, ins.Immediate, flagsNewVal)},
					)
					break
				}
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

func (s *Simulator) printFlags() string {
	if len(s.flags) == 0 {
		return ""
	}

	flags := "flags:->"
	for _, flag := range s.flagOrder {
		if s.flags[flag] {
			flags += fmt.Sprintf("%s ", flag)
		}
	}
	flags += "\n"
	return flags
}
