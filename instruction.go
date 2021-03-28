package main

import (
	"fmt"
)

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
	switch inst.Opcode {
	case 0x00:
		switch inst.Function {
		case 0x00:
			return fmt.Sprintf("SLL rd = %X, rt = %X, ShiftAmount = %d", inst.Rd, inst.Rt, inst.ShiftAmount)
		case 0x02:
			return "SRL"
		case 0x03:
			return "SRA"
		case 0x04:
			return "SLLV"
		case 0x06:
			return "SRLV"
		case 0x07:
			return "SRAV"
		case 0x08:
			return "JR"
		case 0x09:
			return "JALR"
		case 0x0C:
			return "SYSCALL"
		case 0x0D:
			return "BREAK"
		case 0x10:
			return "MFHI"
		case 0x11:
			return "MTHI"
		case 0x12:
			return "MFLO"
		case 0x13:
			return "MTLO"
		case 0x18:
			return "MULT"
		case 0x19:
			return "MULTU"
		case 0x1A:
			return "DIV"
		case 0x1B:
			return "DIVU"
		case 0x20:
			return "ADD"
		case 0x21:
			return "ADDU"
		case 0x22:
			return "SUB"
		case 0x23:
			return "SUBU"
		case 0x24:
			return "AND"
		case 0x25:
			return "OR"
		case 0x26:
			return "XOR"
		case 0x27:
			return "NOR"
		case 0x2A:
			return "SLT"
		case 0x2B:
			return "SLTU"
		default:
			return "N/A"
		}
	case 0x01:
		return "BcondZ"
	case 0x02:
		return fmt.Sprintf("J target = %X", inst.Address)
	case 0x03:
		return "JAL"
	case 0x04:
		return "BEQ"
	case 0x05:
		return "BNE"
	case 0x06:
		return "BLEZ"
	case 0x07:
		return "BGTZ"
	case 0x08:
		return "ADDI"
	case 0x09:
		return fmt.Sprintf("ADDIU rt = %X, rs = %X, Imm = %04X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0A:
		return "SLTI"
	case 0x0B:
		return "SLTIU"
	case 0x0C:
		return "ANDI"
	case 0x0D:
		return fmt.Sprintf("ORI rt = %X, rs = %X, Imm = %04X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0E:
		return "XORI"
	case 0x0F:
		return fmt.Sprintf("LUI rt = %X, Imm = %04X",inst.Rt, inst.Imm16)
	case 0x10:
		return "COP0"
	case 0x11:
		return "COP1"
	case 0x12:
		return "COP2"
	case 0x13:
		return "COP3"
	case 0x20:
		return "LB"
	case 0x21:
		return "LH"
	case 0x22:
		return "LWL"
	case 0x23:
		return "LW"
	case 0x24:
		return "LBU"
	case 0x25:
		return "LHU"
	case 0x26:
		return "LWR"
	case 0x28:
		return "SB"
	case 0x29:
		return "SH"
	case 0x2A:
		return "SWL"
	case 0x2B:
		return fmt.Sprintf("SW rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x2E:
		return "SWR"
	case 0x30:
		return "LWC0"
	case 0x31:
		return "LWC1"
	case 0x32:
		return "LWC2"
	case 0x33:
		return "LWC3"
	case 0x38:
		return "SWC0"
	case 0x39:
		return "SWC1"
	case 0x3A:
		return "SWC2"
	case 0x3B:
		return "SWC3"
	default:
		return "N/A"

	}

}
