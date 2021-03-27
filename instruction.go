package main

// See https://en.wikipedia.org/wiki/MIPS_architecture#Instruction_formats
type Instruction struct {
	Opcode uint32

	Rs, Rt, Rd  uint32
	ShiftAmount uint32
	Function    uint32

	Imm16, Imm16sx uint32

	Address uint32
}

func NewInstruction(value uint32) Instruction {
	return Instruction{
		Opcode: value >> 26,
		Rs:     (value >> 21) & 0x1F,
		Rt:     (value >> 16) & 0x1F,
		Rd:     (value >> 11) & 0x1F,

		Imm16:   value & 0xFFFF,
		Imm16sx: uint32(int16(value & 0xFFFF)),

		Address: value & 0x03FFFFFF,
	}
}

func (insn *Instruction) String() string {
	return "unimplemented"
}
