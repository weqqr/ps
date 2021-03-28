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

		ShiftAmount: (value >> 6) & 0x1F,
		Function:    value & 0x3F,

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
			return fmt.Sprintf("SRL rd = %X, rt = %X, ShiftAmount = %d", inst.Rd, inst.Rt, inst.ShiftAmount)
		case 0x03:
			return fmt.Sprintf("SRA rd = %X, rt = %X, ShiftAmount = %d", inst.Rd, inst.Rt, inst.ShiftAmount)
		case 0x04:
			return fmt.Sprintf("SLLV rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x06:
			return fmt.Sprintf("SRLV rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x07:
			return fmt.Sprintf("SRAV rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x08:
			return fmt.Sprintf("JR rs = %d", inst.Rs)
		case 0x09:
			return fmt.Sprintf("JALR rd = %X, rs = %d", inst.Rd, inst.Rs)
		case 0x0C:
			return "SYSCALL"
		case 0x0D:
			return "BREAK"
		case 0x10:
			return fmt.Sprintf("MFHI rd = %X", inst.Rd)
		case 0x11:
			return fmt.Sprintf("MTHI rs = %d", inst.Rs)
		case 0x12:
			return fmt.Sprintf("MFLO rd = %X", inst.Rd)
		case 0x13:
			return fmt.Sprintf("MTLO rs = %d", inst.Rs)
		case 0x18:
			return fmt.Sprintf("MULT rt = %X, rs = %d", inst.Rt, inst.Rs)
		case 0x19:
			return fmt.Sprintf("MULTU rt = %X, rs = %d", inst.Rt, inst.Rs)
		case 0x1A:
			return fmt.Sprintf("DIV rt = %X, rs = %d", inst.Rt, inst.Rs)
		case 0x1B:
			return fmt.Sprintf("DIVU rt = %X, rs = %d", inst.Rt, inst.Rs)
		case 0x20:
			return fmt.Sprintf("ADD rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x21:
			return fmt.Sprintf("ADDU rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x22:
			return fmt.Sprintf("SUB rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x23:
			return fmt.Sprintf("SUBU rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x24:
			return fmt.Sprintf("AND rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x25:
			return fmt.Sprintf("OR rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x26:
			return fmt.Sprintf("XOR rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x27:
			return fmt.Sprintf("NOR rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x2A:
			return fmt.Sprintf("SLT rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		case 0x2B:
			return fmt.Sprintf("SLTU rd = %X, rt = %X, rs = %d", inst.Rd, inst.Rt, inst.Rs)
		default:
			return "N/A"
		}
	case 0x01:
		return "BcondZ"
	case 0x02:
		return fmt.Sprintf("J target = %X", inst.Address)
	case 0x03:
		return fmt.Sprintf("JAL target = %X", inst.Address)
	case 0x04:
		return fmt.Sprintf("BEQ rt = %d, rs = %X, offset = %X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x05:
		return fmt.Sprintf("BNE rt = %d, rs = %X, offset = %X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x06:
		return fmt.Sprintf("BLEZ rs = %d, offset = %X", inst.Rs, inst.Imm16)
	case 0x07:
		return fmt.Sprintf("BGTZ rs = %d, offset = %X", inst.Rs, inst.Imm16)
	case 0x08:
		return fmt.Sprintf("ADDI rt = %X, rs = %X, Imm = %04X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x09:
		return fmt.Sprintf("ADDIU rt = %X, rs = %X, Imm = %04X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0A:
		return fmt.Sprintf("SLTI rt = %X, rs = %X, Imm = %04X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0B:
		return fmt.Sprintf("SLTIU rt = %X, rs = %X, Imm = %04X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0C:
		return fmt.Sprintf("ANDI rt = %X, rs = %X, Imm = %04X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0D:
		return fmt.Sprintf("ORI rt = %X, rs = %X, Imm = %04X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0E:
		return fmt.Sprintf("XORI rt = %X, rs = %X, Imm = %04X", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0F:
		return fmt.Sprintf("LUI rt = %X, Imm = %04X", inst.Rt, inst.Imm16)
	case 0x10:
		return "COP0"
	case 0x11:
		return "COP1"
	case 0x12:
		return "COP2"
	case 0x13:
		return "COP3"
	case 0x20:
		return fmt.Sprintf("LB rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x21:
		return fmt.Sprintf("LH rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x22:
		return fmt.Sprintf("LWL rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x23:
		return fmt.Sprintf("LW rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x24:
		return fmt.Sprintf("LBU rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x25:
		return fmt.Sprintf("LHU rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x26:
		return fmt.Sprintf("LWR rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x28:
		return fmt.Sprintf("SB rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x29:
		return fmt.Sprintf("SH rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x2A:
		return fmt.Sprintf("SWL rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x2B:
		return fmt.Sprintf("SW rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
	case 0x2E:
		return fmt.Sprintf("SWR rt = %X, offset = %X, base = %X", inst.Rt, inst.Imm16, inst.Rs)
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
