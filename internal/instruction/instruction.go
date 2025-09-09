package instruction

import (
	"fmt"

	"github.com/8086-simulator/internal/bits"
)

type Op string

const (
	MOV Op = "mov"
)

type OperandType int

const (
	OpTypeRegToReg OperandType = iota
	OpTypeImmToReg
)

var regFieldEnc = map[byte]map[bool]string{
	0b000: {false: "al", true: "ax"},
	0b001: {false: "cl", true: "cx"},
	0b010: {false: "dl", true: "dx"},
	0b011: {false: "bl", true: "bx"},
	0b100: {false: "ah", true: "sp"},
	0b101: {false: "ch", true: "bp"},
	0b110: {false: "dh", true: "si"},
	0b111: {false: "bh", true: "di"},
}

var effectiveAddrEnc = map[byte]map[byte]string{
	0b00: {
		0b000: "bx + si",
		0b001: "bx + di",
		0b010: "bp + si",
		0b011: "bp + di",
		0b100: "si",
		0b101: "di",
		0b110: "bp", // direct access
		0b111: "bx",
	},
	0b01: {
		0b000: "bx + si",
		0b001: "bx + di",
		0b010: "bp + si",
		0b011: "bp + di",
		0b100: "si",
		0b101: "di",
		0b110: "bp",
		0b111: "bx",
	},
	0b10: {
		0b000: "bx + si",
		0b001: "bx + di",
		0b010: "bp + si",
		0b011: "bp + di",
		0b100: "si",
		0b101: "di",
		0b110: "bp",
		0b111: "bx",
	},
}

type Instruction struct {
	Op                 Op
	OperandType        OperandType
	DBit               bool
	WBit               bool
	Mod                byte
	Reg                byte
	RM                 byte
	DestRegister       string
	SourceRegister     string
	Text               string
	Immediate          int
	SourceAddr         string
	DestAddr           string
	SourceDisplacement []byte
	DestDisplacement   []byte
}

// NewInstruction creates a new Instruction with all default values
func NewInstruction(instructions []byte, i int, p Pattern) *Instruction {
	return &Instruction{
		Op:          p.Op,
		OperandType: p.OperandType,
	}
}

// formatOperand formats an operand as either a register or memory address with displacement
func (ins *Instruction) formatOperand(addr string, displacement []byte, register string) string {
	if addr != "" {
		result := fmt.Sprintf("[%s", addr)
		if len(displacement) > 0 {
			if len(displacement) == 1 {
				result += fmt.Sprintf(" + %d", bits.ToSigned8(displacement[0]))
			} else {
				result += fmt.Sprintf(" + %d", bits.ToSigned16(displacement[0], displacement[1]))
			}
		}
		result += "]"
		return result
	}
	return register
}

// GetText formats the instruction as a string
func (ins *Instruction) GetText(p *Pattern) string {
	if !ins.DBit {
		tmpAddr := ins.SourceAddr
		tmpDisp := ins.SourceDisplacement
		ins.SourceAddr = ins.DestAddr
		ins.SourceDisplacement = ins.DestDisplacement
		ins.DestAddr = tmpAddr
		ins.DestDisplacement = tmpDisp
	}

	source := ins.formatOperand(ins.SourceAddr, ins.SourceDisplacement, ins.SourceRegister)
	dest := ins.formatOperand(ins.DestAddr, ins.DestDisplacement, ins.DestRegister)

	return fmt.Sprintf("%s %s, %s", p.Op, dest, source)
}

type Pattern struct {
	OpCode                byte
	Op                    Op
	GetOpCode             func(instructions []byte, i int) byte
	OperandType           OperandType
	GetBytesCount         func(p *Pattern, ins *Instruction) int
	GetDBit               func(instructions []byte, i int) bool
	GetWBit               func(instructions []byte, i int) bool
	GetMod                func(instructions []byte, i int) byte
	GetReg                func(instructions []byte, i int) byte
	GetRM                 func(instructions []byte, i int) byte
	GetDestRegister       func(dBit bool, reg byte, rm byte, wBit bool) string
	GetSourceRegister     func(dBit bool, reg byte, rm byte, wBit bool) string
	GetText               func(p *Pattern, ins *Instruction) string
	GetImmediate          func(instructions []byte, i int, ins *Instruction) int
	GetSorceAddr          func(instructions []byte, i int, ins *Instruction) string
	GetDestAddr           func(instructions []byte, i int, ins *Instruction) string
	GetSourceDisplacement func(instructions []byte, i int, ins *Instruction) []byte
	GetDestDisplacement   func(instructions []byte, i int, ins *Instruction) []byte
}

// NewPattern creates a new Pattern with all default functions
func NewPattern() Pattern {
	return Pattern{
		GetOpCode:         func(instructions []byte, i int) byte { return 0 },
		GetBytesCount:     func(_ *Pattern, _ *Instruction) int { return 2 },
		GetDBit:           func(instructions []byte, i int) bool { return false },
		GetWBit:           func(instructions []byte, i int) bool { return false },
		GetMod:            func(instructions []byte, i int) byte { return 0 },
		GetReg:            func(instructions []byte, i int) byte { return 0 },
		GetRM:             func(instructions []byte, i int) byte { return 0 },
		GetDestRegister:   func(_ bool, _ byte, _ byte, _ bool) string { return "" },
		GetSourceRegister: func(_ bool, _ byte, _ byte, _ bool) string { return "" },
		GetText: func(p *Pattern, ins *Instruction) string {
			return ins.GetText(p)
		},
		GetImmediate:          func(_ []byte, _ int, _ *Instruction) int { return 0 },
		GetSorceAddr:          func(_ []byte, _ int, _ *Instruction) string { return "" },
		GetDestAddr:           func(_ []byte, _ int, _ *Instruction) string { return "" },
		GetSourceDisplacement: func(_ []byte, _ int, _ *Instruction) []byte { return nil },
		GetDestDisplacement:   func(_ []byte, _ int, _ *Instruction) []byte { return nil },
	}
}

var Table = []Pattern{
	// MOV
	// MOV - Register/memory to/from register
	func() Pattern {
		p := NewPattern()
		p.OpCode = 0b100010
		p.Op = MOV
		p.OperandType = OpTypeRegToReg
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			defaultInc := 2
			if ins.Mod == 0b00 && ins.RM == 0b110 {
				if ins.SourceAddr == "bp" { // bp is speacial case for direct address
					defaultInc += 1
				} else {
					defaultInc += 2
				}
			} else if ins.Mod == 0b01 {
				defaultInc += 1
			} else if ins.Mod == 0b10 {
				defaultInc += 2
			}
			return defaultInc
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 2, 6) }
		p.GetDBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 1) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 0) }
		p.GetMod = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 6, 2) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 3, 3) }
		p.GetRM = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 0, 3) }
		p.GetSourceRegister = func(dBit bool, reg byte, rm byte, wBit bool) string {
			var sourceReg byte
			if dBit {
				sourceReg = rm
			} else {
				sourceReg = reg
			}
			return regFieldEnc[sourceReg][wBit]
		}
		p.GetDestRegister = func(dBit bool, reg byte, rm byte, wBit bool) string {
			var destReg byte
			if dBit {
				destReg = reg
			} else {
				destReg = rm
			}
			return regFieldEnc[destReg][wBit]
		}
		p.GetSorceAddr = func(_ []byte, _ int, ins *Instruction) string {
			if ins.Mod == 0b11 {
				return ""
			}
			return effectiveAddrEnc[ins.Mod][ins.RM]
		}
		p.GetSourceDisplacement = func(instructions []byte, i int, ins *Instruction) []byte {
			if ins.Mod == 0b00 && ins.RM == 0b110 {
				return instructions[i+2 : i+4]
			} else if ins.Mod == 0b01 {
				return instructions[i+2 : i+3]
			} else if ins.Mod == 0b10 {
				return instructions[i+2 : i+4]
			} else {
				return nil
			}
		}
		p.GetDestDisplacement = func(instructions []byte, i int, ins *Instruction) []byte {
			if ins.Mod == 0b00 && ins.RM == 0b110 {
				return instructions[i+2 : i+4]
			} else if ins.Mod == 0b01 {
				return instructions[i+2 : i+3]
			} else if ins.Mod == 0b10 {
				return instructions[i+2 : i+4]
			} else {
				return nil
			}
		}
		return p
	}(),
	// MOV - Immediate to register
	func() Pattern {
		p := NewPattern()
		p.OpCode = 0b1011
		p.Op = MOV
		p.OperandType = OpTypeImmToReg
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			defaultInc := 2
			if ins.WBit {
				defaultInc++
			}
			return defaultInc
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 4, 4) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 3) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 3) }
		p.GetDestRegister = func(_ bool, reg byte, _ byte, wBit bool) string {
			return regFieldEnc[reg][wBit]
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %s, %d", p.Op, ins.DestRegister, ins.Immediate)
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			if ins.WBit {
				return bits.ToSigned16(instructions[i+1], instructions[i+2])
			}
			return bits.ToSigned8(instructions[i+1])
		}
		return p
	}(),
	// ADD
	// SUB
	// CMP
	// Jumps
}
