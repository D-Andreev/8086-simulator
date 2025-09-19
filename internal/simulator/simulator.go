package simulator

import (
	"fmt"

	"github.com/8086-simulator/internal/bits"
	"github.com/8086-simulator/internal/instruction"
)

type Result struct {
	Text string
}

type Simulator struct {
	Registers     map[string][]byte
	registerOrder []string
	flags         map[string]bool
	flagOrder     []string
}

func NewSimulator() *Simulator {
	s := &Simulator{
		Registers: make(map[string][]byte),
	}
	s.Init()
	return s
}

func (s *Simulator) Init() {
	s.registerOrder = []string{"ax", "bx", "cx", "dx", "sp", "bp", "si", "di"}
	s.Registers = map[string][]byte{
		"ax": {0, 0},
		"bx": {0, 0},
		"cx": {0, 0},
		"dx": {0, 0},
		"sp": {0, 0},
		"bp": {0, 0},
		"si": {0, 0},
		"di": {0, 0},
	}
	s.flags = map[string]bool{
		"Z": false,
		"S": false,
	}
	s.flagOrder = []string{"Z", "S"}
}

func (s *Simulator) printImmediateValue(rawData []byte) uint16 {
	if len(rawData) == 2 {

		return bits.ToUnsigned16(rawData[0], rawData[1])
	}

	return bits.ToUnsigned8(rawData[0])
}

func (s *Simulator) Run(instructions []*instruction.Instruction) ([]*Result, error) {
	results := []*Result{}
	for _, ins := range instructions {
		switch ins.Op {
		case instruction.MOV:
			switch ins.OperandType {
			case instruction.OpTypeImmToReg:
				destPrevVal := s.Registers[ins.DestRegister]
				s.Registers[ins.DestRegister] = ins.Immediate.Raw
				results = append(
					results,
					&Result{
						Text: fmt.Sprintf(
							"%s %s, %d ; %s:0x%x->0x%x",
							ins.Op,
							ins.DestRegister,
							s.printImmediateValue(ins.Immediate.Raw),
							ins.DestRegister,
							destPrevVal[0],
							s.printImmediateValue(ins.Immediate.Raw),
						),
					},
				)
			case instruction.OpTypeRegMemToFromReg:
				destPrevVal := s.Registers[ins.DestRegister]
				sourceVal := s.Registers[ins.SourceRegister]
				s.Registers[ins.DestRegister] = sourceVal
				results = append(
					results,
					&Result{Text: fmt.Sprintf("%s ; %s:0x%x->0x%x", ins.Text, ins.DestRegister, destPrevVal[0], sourceVal[0])},
				)
			default:
				return nil, fmt.Errorf("unsupported operand type: %d", ins.OperandType)
			}
		case instruction.ADD, instruction.SUB, instruction.CMP:
			switch ins.OperandType {
			case instruction.OpTypeImmToReg:
				destPrevVal := s.Registers[ins.DestRegister]
				s.Registers[ins.DestRegister] = s.doArithmeticOp(ins, destPrevVal, ins.Immediate.Raw, ins.DestRegister)
				flagsNewVal := s.printFlags()
				if flagsNewVal != "" {
					results = append(
						results,
						&Result{
							Text: fmt.Sprintf(
								"%s ; %s:0x%x->0x%x flags:->%s",
								ins.Text,
								ins.DestRegister,
								s.printImmediateValue(destPrevVal),
								s.printImmediateValue(s.Registers[ins.DestRegister]),
								flagsNewVal,
							),
						},
					)
				} else {
					results = append(
						results,
						&Result{
							Text: fmt.Sprintf(
								"%s ; %s:0x%x->0x%x",
								ins.Text,
								ins.DestRegister,
								s.printImmediateValue(destPrevVal),
								s.printImmediateValue(s.Registers[ins.DestRegister]),
							),
						},
					)
				}
			case instruction.OpTypeRegMemToFromReg:
				destPrevVal := s.Registers[ins.DestRegister]
				sourceVal := s.Registers[ins.SourceRegister]
				flagsPrevVal := s.printFlags()
				result := s.doArithmeticOp(ins, destPrevVal, sourceVal, ins.DestRegister)
				flagsNewVal := s.printFlags()
				if ins.Op == instruction.CMP {
					// CMP doesn't modify the destination register
					if flagsPrevVal != flagsNewVal {
						results = append(
							results,
							&Result{Text: fmt.Sprintf("%s ; flags:%s->%s", ins.Text, flagsPrevVal, flagsNewVal)},
						)
					} else {
						results = append(
							results,
							&Result{Text: fmt.Sprintf("%s ; flags:%s->%s", ins.Text, flagsPrevVal, flagsNewVal)},
						)
					}
				} else {
					if flagsNewVal != "" {
						results = append(
							results,
							&Result{
								Text: fmt.Sprintf(
									"%s ; %s:0x%x->0x%x flags:->%s",
									ins.Text,
									ins.DestRegister,
									s.printImmediateValue(destPrevVal),
									s.printImmediateValue(result),
									flagsNewVal,
								),
							},
						)
					} else {
						results = append(
							results,
							&Result{
								Text: fmt.Sprintf(
									"%s ; %s:0x%x->0x%x",
									ins.Text,
									ins.DestRegister,
									s.printImmediateValue(destPrevVal),
									s.printImmediateValue(result),
								),
							},
						)
					}
				}
			default:
				return nil, fmt.Errorf("unsupported operand type: %d", ins.OperandType)
			}
		default:
			return nil, fmt.Errorf("unsupported instruction: %s", ins.Op)
		}
	}

	return results, nil
}

func (s *Simulator) doArithmeticOp(ins *instruction.Instruction, destPrevVal []byte, sourceVal []byte, destRegister string) []byte {
	switch ins.Op {
	case instruction.ADD:
		s.Registers[destRegister] = bits.Add(destPrevVal, sourceVal)
		s.flags["Z"] = bits.IsZero(s.Registers[destRegister])
		s.flags["S"] = bits.IsNegative(s.Registers[destRegister])
		return s.Registers[destRegister]
	case instruction.SUB:
		s.Registers[destRegister] = bits.Sub(destPrevVal, sourceVal)
		s.flags["Z"] = bits.IsZero(s.Registers[destRegister])
		s.flags["S"] = bits.IsNegative(s.Registers[destRegister])
		return s.Registers[destRegister]
	case instruction.CMP:
		res := bits.Sub(destPrevVal, sourceVal)
		s.flags["Z"] = bits.IsZero(res)
		s.flags["S"] = bits.IsNegative(res)
		return destPrevVal // CMP doesn't modify the destination register
	}
	return nil
}

func (s *Simulator) printFlags() string {
	flags := ""
	for _, flag := range s.flagOrder {
		if s.flags[flag] {
			flags += flag
		}
	}
	return flags
}
