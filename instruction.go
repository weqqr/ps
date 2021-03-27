package main

// See https://en.wikipedia.org/wiki/MIPS_architecture#Instruction_formats
type Instruction struct {
	opcode uint32

	rs, rt, rd uint32
	shamt      uint32
	funct      uint32

	imm16, imm16sx uint32

	address uint32
}

func NewInstruction(value uint32) Instruction {
	return Instruction{
		opcode: value >> 26,
		rs:     (value >> 21) & 0x1F,
		rt:     (value >> 16) & 0x1F,
		rd:     (value >> 11) & 0x1F,

		imm16:   value & 0xFFFF,
		imm16sx: uint32(int16(value & 0xFFFF)),

		address: value & 0x03FFFFFF,
	}
}

func (insn *Instruction) String() string {
	return "unimplemented"
}
