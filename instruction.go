package main

// See https://en.wikipedia.org/wiki/MIPS_architecture#Instruction_formats
type Instruction struct {
	// Opcode is primary opcode value
	Opcode uint32

	// Rs, Rt, Rd are register indices
	Rs, Rt, Rd uint32
	// ShiftAmount is offset value
	ShiftAmount uint32
	// Function is secondary opcode value
	Function uint32

	// Imm16 is immediate value extended to 32 bits
	Imm16 uint32
	// Imm16sx is sign extended immediate value
	Imm16sx uint32

	// Address is address value
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

func (inst *Instruction) String() string {
	return "unimplemented"
}
