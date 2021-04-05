package ps

import (
	"fmt"
)

// See https://en.wikipedia.org/wiki/MIPS_architecture#Instruction_formats
type Instruction struct {
	// Raw contains raw, unparsed instruction
	Raw uint32

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
		Raw: value,

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

func (inst Instruction) String() string {
	switch inst.Opcode {
	case 0x00:
		switch inst.Function {
		case 0x00:
			return fmt.Sprintf("SLL     $%d, $%d, %d", inst.Rd, inst.Rt, inst.ShiftAmount)
		case 0x02:
			return fmt.Sprintf("SRL     $%d, $%d, %d", inst.Rd, inst.Rt, inst.ShiftAmount)
		case 0x03:
			return fmt.Sprintf("SRA     $%d, $%d, %d", inst.Rd, inst.Rt, inst.ShiftAmount)
		case 0x04:
			return fmt.Sprintf("SLLV    $%d, $%d, $%d", inst.Rd, inst.Rt, inst.Rs)
		case 0x06:
			return fmt.Sprintf("SRLV    $%d, $%d, $%d", inst.Rd, inst.Rt, inst.Rs)
		case 0x07:
			return fmt.Sprintf("SRAV    $%d, $%d, $%d", inst.Rd, inst.Rt, inst.Rs)
		case 0x08:
			return fmt.Sprintf("JR      $%d", inst.Rs)
		case 0x09:
			return fmt.Sprintf("JALR    $%d, $%d", inst.Rd, inst.Rs)
		case 0x0C:
			return "SYSCALL"
		case 0x0D:
			return "BREAK"
		case 0x10:
			return fmt.Sprintf("MFHI    $%d", inst.Rd)
		case 0x11:
			return fmt.Sprintf("MTHI    $%d", inst.Rs)
		case 0x12:
			return fmt.Sprintf("MFLO    $%d", inst.Rd)
		case 0x13:
			return fmt.Sprintf("MTLO    $%d", inst.Rs)
		case 0x18:
			return fmt.Sprintf("MULT    $%d, $%d", inst.Rt, inst.Rs)
		case 0x19:
			return fmt.Sprintf("MULTU   $%d, $%d", inst.Rt, inst.Rs)
		case 0x1A:
			return fmt.Sprintf("DIV     $%d, $%d", inst.Rt, inst.Rs)
		case 0x1B:
			return fmt.Sprintf("DIVU    $%d, $%d", inst.Rt, inst.Rs)
		case 0x20:
			return fmt.Sprintf("ADD     $%d, $%d, $%d", inst.Rd, inst.Rs, inst.Rt)
		case 0x21:
			return fmt.Sprintf("ADDU    $%d, $%d, $%d", inst.Rd, inst.Rs, inst.Rt)
		case 0x22:
			return fmt.Sprintf("SUB     $%d, $%d, $%d", inst.Rd, inst.Rt, inst.Rs)
		case 0x23:
			return fmt.Sprintf("SUBU    $%d, $%d, $%d", inst.Rd, inst.Rs, inst.Rt)
		case 0x24:
			return fmt.Sprintf("AND     $%d, $%d, $%d", inst.Rd, inst.Rs, inst.Rt)
		case 0x25:
			return fmt.Sprintf("OR      $%d, $%d, $%d", inst.Rd, inst.Rs, inst.Rt)
		case 0x26:
			return fmt.Sprintf("XOR     $%d, $%d, $%d", inst.Rd, inst.Rs, inst.Rt)
		case 0x27:
			return fmt.Sprintf("NOR     $%d, $%d, $%d", inst.Rd, inst.Rs, inst.Rt)
		case 0x2A:
			return fmt.Sprintf("SLT     $%d, $%d, $%d", inst.Rd, inst.Rs, inst.Rt)
		case 0x2B:
			return fmt.Sprintf("SLTU    $%d, $%d, $%d", inst.Rd, inst.Rs, inst.Rt)
		default:
			return fmt.Sprintf("Invalid SPECIAL opcode: %02Xh", inst.Function)
		}
	case 0x01:
		return fmt.Sprintf("Invalid BCOND opcode: %02Xh", inst.Rs)
	case 0x02:
		return fmt.Sprintf("J       %08Xh", inst.Address)
	case 0x03:
		return fmt.Sprintf("JAL     %08Xh", inst.Address)
	case 0x04:
		return fmt.Sprintf("BEQ     $%d, $%d, %Xh", inst.Rs, inst.Rt, inst.Imm16)
	case 0x05:
		return fmt.Sprintf("BNE     $%d, $%d, %Xh", inst.Rs, inst.Rt, inst.Imm16)
	case 0x06:
		return fmt.Sprintf("BLEZ    $%d, %Xh", inst.Rs, inst.Imm16)
	case 0x07:
		return fmt.Sprintf("BGTZ    $%d, %Xh", inst.Rs, inst.Imm16)
	case 0x08:
		return fmt.Sprintf("ADDI    $%d, $%d, %Xh", inst.Rt, inst.Rs, inst.Imm16sx)
	case 0x09:
		return fmt.Sprintf("ADDIU   $%d, $%d, %Xh", inst.Rt, inst.Rs, inst.Imm16sx)
	case 0x0A:
		return fmt.Sprintf("SLTI    $%d, $%d, %Xh", inst.Rt, inst.Rs, inst.Imm16sx)
	case 0x0B:
		return fmt.Sprintf("SLTIU   $%d, $%d, %Xh", inst.Rt, inst.Rs, inst.Imm16sx)
	case 0x0C:
		return fmt.Sprintf("ANDI    $%d, $%d, %Xh", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0D:
		return fmt.Sprintf("ORI     $%d, $%d, %Xh", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0E:
		return fmt.Sprintf("XORI    $%d, $%d, %Xh", inst.Rt, inst.Rs, inst.Imm16)
	case 0x0F:
		return fmt.Sprintf("LUI     $%d, %Xh", inst.Rt, inst.Imm16)
	case 0x10:
		return "COP0"
	case 0x12:
		return "COP2"
	case 0x20:
		return fmt.Sprintf("LB      $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x21:
		return fmt.Sprintf("LH      $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x22:
		return fmt.Sprintf("LWL     $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x23:
		return fmt.Sprintf("LW      $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x24:
		return fmt.Sprintf("LBU     $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x25:
		return fmt.Sprintf("LHU     $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x26:
		return fmt.Sprintf("LWR     $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x28:
		return fmt.Sprintf("SB      $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x29:
		return fmt.Sprintf("SH      $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x2A:
		return fmt.Sprintf("SWL     $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x2B:
		return fmt.Sprintf("SW      $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
	case 0x2E:
		return fmt.Sprintf("SWR     $%d, %Xh($%d)", inst.Rt, inst.Imm16sx, inst.Rs)
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
		return fmt.Sprintf("Unknown opcode: %02xh", inst.Opcode)
	}
}
