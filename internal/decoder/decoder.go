package decoder

import (
	"fmt"

	"github.com/8086-simulator/internal/instruction"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(data []byte) ([]*instruction.Instruction, error) {
	i := 0
	var instructions []*instruction.Instruction
	for i < len(data) {
		instrFound := false
		for _, p := range instruction.Table {
			parsedOpCode := p.GetOpCode(data, i)
			if parsedOpCode != p.OpCode {
				continue
			}

			instrFound = true
			ins := instruction.NewInstruction(data, i, p)
			ins.DBit = p.GetDBit(data, i)
			ins.WBit = p.GetWBit(data, i)
			ins.Reg = p.GetReg(data, i)
			ins.RM = p.GetRM(data, i)
			ins.Mod = p.GetMod(data, i)
			ins.SBit = p.GetSBit(data, i)
			ins.DestRegister = p.GetDestRegister(ins)
			ins.SourceRegister = p.GetSourceRegister(ins)
			ins.SourceDisplacement = p.GetSourceDisplacement(data, i, ins)
			ins.DestDisplacement = p.GetDestDisplacement(data, i, ins)
			ins.Immediate = p.GetImmediate(data, i, ins)
			ins.SourceAddr = p.GetSorceAddr(data, i, ins)
			ins.DestAddr = p.GetDestAddr(data, i, ins)
			ins.Text = p.GetText(p, ins)
			instructions = append(instructions, ins)
			i += p.GetBytesCount(p, ins)
			break
		}

		if !instrFound {
			return nil, fmt.Errorf("instruction not found at index %d with opcode %d", i, data[i])
		}
	}

	return instructions, nil
}
