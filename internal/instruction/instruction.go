package instruction

import (
	"fmt"

	"github.com/8086-simulator/internal/bits"
)

type Op string

const (
	MOV    Op = "mov"
	ADD    Op = "add"
	SUB    Op = "sub"
	CMP    Op = "cmp"
	JNZ    Op = "jnz"
	JE     Op = "je"
	JL     Op = "jl"
	JLE    Op = "jle"
	JB     Op = "jb"
	JBE    Op = "jbe"
	JP     Op = "jp"
	JO     Op = "jo"
	JS     Op = "js"
	JNE    Op = "jne"
	JNL    Op = "jnl"
	JG     Op = "jg"
	JNB    Op = "jnb"
	JA     Op = "ja"
	JNP    Op = "jnp"
	JNO    Op = "jno"
	JNS    Op = "jns"
	LOOP   Op = "loop"
	LOOPZ  Op = "loopz"
	LOOPNZ Op = "loopnz"
	JCXZ   Op = "jcxz"
)

type OperandType int

var regOpCodeEnc = map[byte]Op{
	0b000: ADD,
	0b101: SUB,
	0b111: CMP,
}

const (
	OpTypeRegMemToFromReg OperandType = iota
	OpTypeImmToReg
	OpTypeImmToAcc
	OpTypeJump
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
	SBit               bool
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
func NewInstruction(instructions []byte, i int, p *Pattern) *Instruction {
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

func (ins *Instruction) GetSourceReg() string {
	var sourceReg byte
	if ins.DBit {
		sourceReg = ins.RM
	} else {
		sourceReg = ins.Reg
	}
	return regFieldEnc[sourceReg][ins.WBit]
}

func (ins *Instruction) GetDestReg() string {
	var destReg byte
	if ins.DBit {
		destReg = ins.Reg
	} else {
		destReg = ins.RM
	}
	return regFieldEnc[destReg][ins.WBit]
}

func (ins *Instruction) GetFromToRegMemInstrByteCount() int {
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

func (ins *Instruction) GetImmToRegInstrByteCount() int {
	defaultInc := 2
	if ins.WBit {
		defaultInc++
	}
	return defaultInc
}

type Pattern struct {
	OpCode                byte
	Op                    Op
	GetOpCode             func(instructions []byte, i int) byte
	OperandType           OperandType
	GetBytesCount         func(p *Pattern, ins *Instruction) int
	GetDBit               func(instructions []byte, i int) bool
	GetWBit               func(instructions []byte, i int) bool
	GetSBit               func(instructions []byte, i int) bool
	GetMod                func(instructions []byte, i int) byte
	GetReg                func(instructions []byte, i int) byte
	GetRM                 func(instructions []byte, i int) byte
	GetDestRegister       func(ins *Instruction) string
	GetSourceRegister     func(ins *Instruction) string
	GetText               func(p *Pattern, ins *Instruction) string
	GetImmediate          func(instructions []byte, i int, ins *Instruction) int
	GetSorceAddr          func(instructions []byte, i int, ins *Instruction) string
	GetDestAddr           func(instructions []byte, i int, ins *Instruction) string
	GetSourceDisplacement func(instructions []byte, i int, ins *Instruction) []byte
	GetDestDisplacement   func(instructions []byte, i int, ins *Instruction) []byte
}

// NewPattern creates a new Pattern with all default functions
func NewPattern() *Pattern {
	return &Pattern{
		GetOpCode:         func(instructions []byte, i int) byte { return 0 },
		GetBytesCount:     func(_ *Pattern, _ *Instruction) int { return 2 },
		GetDBit:           func(instructions []byte, i int) bool { return false },
		GetWBit:           func(instructions []byte, i int) bool { return false },
		GetSBit:           func(instructions []byte, i int) bool { return false },
		GetMod:            func(instructions []byte, i int) byte { return 0 },
		GetReg:            func(instructions []byte, i int) byte { return 0 },
		GetRM:             func(instructions []byte, i int) byte { return 0 },
		GetDestRegister:   func(ins *Instruction) string { return ins.GetDestReg() },
		GetSourceRegister: func(ins *Instruction) string { return ins.GetSourceReg() },
		GetText:           func(p *Pattern, ins *Instruction) string { return ins.GetText(p) },
		GetImmediate: func(instructions []byte, i int, ins *Instruction) int {
			if ins.WBit {
				return bits.ToSigned16(instructions[i+1], instructions[i+2])
			}
			return bits.ToSigned8(instructions[i+1])
		},
		GetSorceAddr: func(_ []byte, _ int, ins *Instruction) string {
			if ins.Mod == 0b11 {
				return ""
			}
			return effectiveAddrEnc[ins.Mod][ins.RM]
		},
		GetDestAddr: func(_ []byte, _ int, _ *Instruction) string { return "" },
		GetSourceDisplacement: func(instructions []byte, i int, ins *Instruction) []byte {
			if ins.Mod == 0b00 && ins.RM == 0b110 {
				return instructions[i+2 : i+4]
			} else if ins.Mod == 0b01 {
				return instructions[i+2 : i+3]
			} else if ins.Mod == 0b10 {
				return instructions[i+2 : i+4]
			} else {
				return nil
			}
		},
		GetDestDisplacement: func(instructions []byte, i int, ins *Instruction) []byte {
			if ins.Mod == 0b00 && ins.RM == 0b110 {
				return instructions[i+2 : i+4]
			} else if ins.Mod == 0b01 {
				return instructions[i+2 : i+3]
			} else if ins.Mod == 0b10 {
				return instructions[i+2 : i+4]
			} else {
				return nil
			}
		},
	}
}

var Table = []*Pattern{
	// MOV
	// MOV - Register/memory to/from register
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b100010
		p.Op = MOV
		p.OperandType = OpTypeRegMemToFromReg
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return ins.GetFromToRegMemInstrByteCount()
		}
		p.GetImmediate = func(_ []byte, _ int, _ *Instruction) int { return 0 }
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 2, 6) }
		p.GetDBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 1) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 0) }
		p.GetMod = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 6, 2) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 3, 3) }
		p.GetRM = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 0, 3) }
		return p
	}(),
	// MOV - Immediate to register
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b1011
		p.Op = MOV
		p.OperandType = OpTypeImmToReg
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return ins.GetImmToRegInstrByteCount()
		}
		p.GetDestRegister = func(ins *Instruction) string {
			return regFieldEnc[ins.Reg][ins.WBit]
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 4, 4) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 3) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 3) }
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %s, %d", p.Op, ins.DestRegister, ins.Immediate)
		}
		return p
	}(),

	// ADD
	// ADD - Register/memory to/from register
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b000000
		p.Op = ADD
		p.OperandType = OpTypeRegMemToFromReg
		p.GetImmediate = func(_ []byte, _ int, _ *Instruction) int { return 0 }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return ins.GetFromToRegMemInstrByteCount()
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 2, 6) }
		p.GetDBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 1) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 0) }
		p.GetMod = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 6, 2) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 3, 3) }
		p.GetRM = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 0, 3) }
		return p
	}(),

	// ADD - Immediate to accumulator
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b0000010
		p.Op = ADD
		p.OperandType = OpTypeImmToAcc
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return ins.GetImmToRegInstrByteCount()
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 1, 7) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 0) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 3) }
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %s, %d", p.Op, ins.DestRegister, ins.Immediate)
		}
		return p
	}(),

	// SUB
	// SUB - Register/memory to/from register
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b001010
		p.Op = SUB
		p.OperandType = OpTypeRegMemToFromReg
		p.GetImmediate = func(_ []byte, _ int, _ *Instruction) int { return 0 }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return ins.GetFromToRegMemInstrByteCount()
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 2, 6) }
		p.GetDBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 1) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 0) }
		p.GetMod = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 6, 2) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 3, 3) }
		p.GetRM = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 0, 3) }
		return p
	}(),
	// SUB - Immediate to accumulator
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b0010110
		p.Op = SUB
		p.OperandType = OpTypeImmToAcc
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return ins.GetImmToRegInstrByteCount()
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 1, 7) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 0) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 3) }
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %s, %d", p.Op, ins.DestRegister, ins.Immediate)
		}
		return p
	}(),
	// ADD, SUB, CMP - Immediate to register/memory
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b100000
		p.OperandType = OpTypeImmToReg
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			inc := 3
			if ins.WBit && !ins.SBit {
				inc += 1
			}
			switch ins.Mod {
			case 0b00:
				if ins.RM == 0b110 {
					inc += 2
				}

			case 0b11:
				inc += 0
			case 0b01:
				inc += 3
			case 0b10:
				inc += 2
			default:
				return 0
			}
			return inc
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 2, 6) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 0) }
		p.GetSBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 1) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 3, 3) }
		p.GetRM = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 0, 3) }
		p.GetMod = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 6, 2) }
		p.GetText = func(p *Pattern, ins *Instruction) string {
			if !ins.DBit {
				tmpAddr := ins.SourceAddr
				tmpDisp := ins.SourceDisplacement
				ins.SourceAddr = ins.DestAddr
				ins.SourceDisplacement = ins.DestDisplacement
				ins.DestAddr = tmpAddr
				ins.DestDisplacement = tmpDisp
			}

			source := fmt.Sprintf("%d", ins.Immediate)
			dest := ins.formatOperand(ins.DestAddr, ins.DestDisplacement, ins.DestRegister)
			insType := "byte"
			if ins.WBit {
				insType = "word"
			}

			op := regOpCodeEnc[ins.Reg]
			if ins.Mod == 0b11 {
				return fmt.Sprintf("%s %s, %s", op, dest, source)
			}

			return fmt.Sprintf("%s %s %s, %s", op, insType, dest, source)
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			idx := 0
			switch ins.Mod {
			case 0b00:
				if ins.RM == 0b110 {
					idx = i + 4
				} else {
					idx = i + 2
				}
			case 0b01:
				idx = i + 3
			case 0b10:
				idx = i + 4
			case 0b11:
				idx = i + 2
			}

			return bits.ToSigned8(instructions[idx])
		}

		return p
	}(),
	// CMP
	// CMP - Register/memory to/from register
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b001110
		p.Op = CMP
		p.OperandType = OpTypeRegMemToFromReg
		p.GetImmediate = func(_ []byte, _ int, _ *Instruction) int { return 0 }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return ins.GetFromToRegMemInstrByteCount()
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 2, 6) }
		p.GetDBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 1) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 0) }
		p.GetMod = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 6, 2) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 3, 3) }
		p.GetRM = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i+1], 0, 3) }
		return p
	}(),
	// CMP - Immediate to accumulator
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b0011110
		p.Op = CMP
		p.OperandType = OpTypeImmToAcc
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return ins.GetImmToRegInstrByteCount()
		}
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 1, 7) }
		p.GetWBit = func(instructions []byte, i int) bool { return bits.GetBit(instructions[i], 0) }
		p.GetReg = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 3) }
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %s, %d", p.Op, ins.DestRegister, ins.Immediate)
		}
		return p
	}(),
	// Jumps
	// JNZ
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01110101
		p.Op = JNZ
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JE
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01110100
		p.Op = JE
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JL
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01111100
		p.Op = JL
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JLE
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01111110
		p.Op = JLE
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JB
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01110010
		p.Op = JB
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JBE
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01110110
		p.Op = JBE
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JP
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01111010
		p.Op = JP
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JO
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01110000
		p.Op = JO
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JS
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01111000
		p.Op = JS
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JNE
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01110101
		p.Op = JNE
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JNL
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01111101
		p.Op = JNL
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JG
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01111111
		p.Op = JG
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JNB
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01110011
		p.Op = JNB
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JA
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01110111
		p.Op = JA
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JNP
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01111011
		p.Op = JNP
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JNO
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01110001
		p.Op = JNO
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JNS
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b01111001
		p.Op = JNS
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// LOOP
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b11100010
		p.Op = LOOP
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// LOOPZ
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b11100001
		p.Op = LOOPZ
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// LOOPNZ
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b11100000
		p.Op = LOOPNZ
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
	// JCXZ
	func() *Pattern {
		p := NewPattern()
		p.OpCode = 0b11100011
		p.Op = JCXZ
		p.OperandType = OpTypeJump
		p.GetOpCode = func(instructions []byte, i int) byte { return bits.GetBits(instructions[i], 0, 8) }
		p.GetBytesCount = func(_ *Pattern, ins *Instruction) int {
			return 2
		}
		p.GetImmediate = func(instructions []byte, i int, ins *Instruction) int {
			return bits.ToSigned8(instructions[i+1])
		}
		p.GetText = func(p *Pattern, ins *Instruction) string {
			return fmt.Sprintf("%s %d", p.Op, ins.Immediate)
		}
		return p
	}(),
}
