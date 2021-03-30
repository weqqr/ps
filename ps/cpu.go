package ps

import (
	"log"
)

type CPU struct {
	// GPR - General Purpose Registers.
	// The content of GPR[0] is always zero.
	// Attempts to alter the content of GPR[0] have no effect.
	GPR []uint32
	// Pc is a program counter
	Pc uint32
	// PcNext is a next program counter
	PcNext uint32
}

func NewCPU() CPU {
	return CPU{
		GPR:    make([]uint32, 32),
		Pc:     0xBFC00000, // Bios start
		PcNext: 0xBFC00004,
	}
}

func (cpu *CPU) SLL(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rt] << instruction.ShiftAmount
}

func (cpu *CPU) SRL(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rt] >> instruction.ShiftAmount
}

func (cpu *CPU) SRA(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = uint32(int32(cpu.GPR[instruction.Rt]) << instruction.ShiftAmount)
}

func (cpu *CPU) SLLV(instruction Instruction, bus *Bus) {
	s := cpu.GPR[instruction.Rs] & 0x1F
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rt] << s
}

func (cpu *CPU) SRLV(instruction Instruction, bus *Bus) {
	s := cpu.GPR[instruction.Rs] & 0x1F
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rt] >> s
}

func (cpu *CPU) SRAV(instruction Instruction, bus *Bus) {
	s := cpu.GPR[instruction.Rs] & 0x1F
	cpu.GPR[instruction.Rd] = uint32(int32(cpu.GPR[instruction.Rt]) >> s)
}

func (cpu *CPU) JR(instruction Instruction, bus *Bus) {
	cpu.PcNext = cpu.GPR[instruction.Rs]
}

func (cpu *CPU) JALR(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = cpu.Pc + 8
	cpu.PcNext = cpu.GPR[instruction.Rs]
}

func (cpu *CPU) SYSCALL(instruction Instruction, bus *Bus) {
	//TODO SystemCallException
}

func (cpu *CPU) BREAK(instruction Instruction, bus *Bus) {
	//TODO BreakpointException
}

func (cpu *CPU) MFHI(instruction Instruction, bus *Bus) {
	//TODO implementation of the quotient and remainder
}

func (cpu *CPU) MTHI(instruction Instruction, bus *Bus) {
	//TODO implementation of the quotient and remainder
}

func (cpu *CPU) MFLO(instruction Instruction, bus *Bus) {
	//TODO implementation of the quotient and remainder
}

func (cpu *CPU) MTLO(instruction Instruction, bus *Bus) {
	//TODO implementation of the quotient and remainder
}

func (cpu *CPU) MULT(instruction Instruction, bus *Bus) {
	//TODO implementation of the quotient and remainder
}

func (cpu *CPU) MULTU(instruction Instruction, bus *Bus) {
	//TODO implementation of the quotient and remainder
}

func (cpu *CPU) DIV(instruction Instruction, bus *Bus) {
	//TODO implementation of the quotient and remainder
}

func (cpu *CPU) DIVU(instruction Instruction, bus *Bus) {
	//TODO implementation of the quotient and remainder
}

func (cpu *CPU) ADD(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = uint32(int32(cpu.GPR[instruction.Rs]) + int32(cpu.GPR[instruction.Rt]))
}

func (cpu *CPU) ADDU(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rs] + cpu.GPR[instruction.Rt]
}

func (cpu *CPU) SUB(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = uint32(int32(cpu.GPR[instruction.Rs]) - int32(cpu.GPR[instruction.Rt]))
}

func (cpu *CPU) SUBU(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rs] - cpu.GPR[instruction.Rt]
}

func (cpu *CPU) AND(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rs] & cpu.GPR[instruction.Rt]
}

func (cpu *CPU) OR(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rs] | cpu.GPR[instruction.Rt]
}

func (cpu *CPU) XOR(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rs] ^ cpu.GPR[instruction.Rt]
}

func (cpu *CPU) NOR(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = ^(cpu.GPR[instruction.Rs] | cpu.GPR[instruction.Rt])
}

func (cpu *CPU) SLT(instruction Instruction, bus *Bus) {
	if cpu.GPR[instruction.Rs] < cpu.GPR[instruction.Rt] {
		cpu.GPR[instruction.Rd] = 1
	} else {
		cpu.GPR[instruction.Rd] = 0
	}
}

func (cpu *CPU) SLTU(instruction Instruction, bus *Bus) {
	if (cpu.GPR[instruction.Rs] >> 1) < (cpu.GPR[instruction.Rt] >> 1) {
		cpu.GPR[instruction.Rd] = 1
	} else {
		cpu.GPR[instruction.Rd] = 0
	}
}

func (cpu *CPU) BLTZ(instruction Instruction, bus *Bus) {
	instruction.Address = instruction.Imm16 << 2
	if (cpu.GPR[instruction.Rs] >> 31) == 1 {
		cpu.PcNext += instruction.Address
	}
}

func (cpu *CPU) BGEZ(instruction Instruction, bus *Bus) {
	instruction.Address = instruction.Imm16 << 2
	if (cpu.GPR[instruction.Rs] >> 31) == 0 {
		cpu.PcNext += instruction.Address
	}
}

func (cpu *CPU) BLTZAL(instruction Instruction, bus *Bus) {
	instruction.Address = instruction.Imm16 << 2
	if (cpu.GPR[instruction.Rs] >> 31) == 1 {
		cpu.PcNext += instruction.Address
	}
	cpu.GPR[31] = cpu.Pc + 8
}

func (cpu *CPU) BGEZAL(instruction Instruction, bus *Bus) {
	instruction.Address = instruction.Imm16 << 2
	if (cpu.GPR[instruction.Rs] >> 31) == 0 {
		cpu.PcNext += instruction.Address
	}
	cpu.GPR[31] = cpu.Pc + 8
}

func (cpu *CPU) J(instruction Instruction, bus *Bus) {
	cpu.PcNext = cpu.Pc&0xF0000000 | (instruction.Address << 2)
}

func (cpu *CPU) JAL(instruction Instruction, bus *Bus) {
	cpu.GPR[31] = cpu.Pc + 8
	cpu.PcNext = cpu.Pc&0xF0000000 | (instruction.Address << 2)
}

func (cpu *CPU) BEQ(instruction Instruction, bus *Bus) {
	instruction.Address = instruction.Imm16sx << 2
	if cpu.GPR[instruction.Rs] == cpu.GPR[instruction.Rt] {
		cpu.PcNext += instruction.Address
	}
}

func (cpu *CPU) BNE(instruction Instruction, bus *Bus) {
	instruction.Address = instruction.Imm16sx << 2
	if cpu.GPR[instruction.Rs] != cpu.GPR[instruction.Rt] {
		cpu.PcNext += instruction.Address
	}
}

func (cpu *CPU) BLEZ(instruction Instruction, bus *Bus) {
	instruction.Address = instruction.Imm16sx << 2
	if (cpu.GPR[instruction.Rs]>>31) == 1 || cpu.GPR[instruction.Rs] == 0 {
		cpu.PcNext += instruction.Address
	}
}

func (cpu *CPU) BGTZ(instruction Instruction, bus *Bus) {
	instruction.Address = instruction.Imm16sx << 2
	if (cpu.GPR[instruction.Rs]>>31) == 0 && cpu.GPR[instruction.Rs] != 0 {
		cpu.PcNext += instruction.Address
	}
}

func (cpu *CPU) ADDI(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rt] = cpu.GPR[instruction.Rs] + instruction.Imm16
}

func (cpu *CPU) ADDIU(instruction Instruction, bus *Bus) {
	//TODO 32-bit-overflow
	cpu.GPR[instruction.Rt] = cpu.GPR[instruction.Rs] + instruction.Imm16sx
}

func (cpu *CPU) SLTI(instruction Instruction, bus *Bus) {
	if cpu.GPR[instruction.Rs] < instruction.Imm16 {
		cpu.GPR[instruction.Rd] = 1
	} else {
		cpu.GPR[instruction.Rd] = 0
	}
}

func (cpu *CPU) SLTIU(instruction Instruction, bus *Bus) {
	if (cpu.GPR[instruction.Rs] >> 1) < instruction.Imm16 {
		cpu.GPR[instruction.Rd] = 1
	} else {
		cpu.GPR[instruction.Rd] = 0
	}
}

func (cpu *CPU) ANDI(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rt] = instruction.Imm16 & cpu.GPR[instruction.Rs]
}

func (cpu *CPU) ORI(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rt] = cpu.GPR[instruction.Rs] | instruction.Imm16
}

func (cpu *CPU) XORI(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rt] = cpu.GPR[instruction.Rs] ^ instruction.Imm16
}

func (cpu *CPU) LUI(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rt] = instruction.Imm16 << 16
}

func (cpu *CPU) MFC(instruction Instruction, bus *Bus) {
	//TODO Coprocessor
}

func (cpu *CPU) CFC(instruction Instruction, bus *Bus) {
	//TODO Coprocessor
}

func (cpu *CPU) MTC(instruction Instruction, bus *Bus) {
	//TODO Coprocessor
}

func (cpu *CPU) CTC(instruction Instruction, bus *Bus) {
	//TODO Coprocessor
}

func (cpu *CPU) LB(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) LH(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) LW(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) LWL(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) LBU(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) LHU(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) LWR(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) SB(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) SH(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) SWL(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) SW(instruction Instruction, bus *Bus) {
	address := instruction.Imm16 + cpu.GPR[instruction.Rs]
	bus.StoreWord(address, cpu.GPR[instruction.Rt])
}

func (cpu *CPU) SWR(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) Execute(instruction Instruction, bus *Bus) {
	switch instruction.Opcode {
	case 0x00:
		switch instruction.Function {
		case 0x00:
			cpu.SLL(instruction, bus)
		case 0x02:
			cpu.SRL(instruction, bus)
		case 0x03:
			cpu.SRA(instruction, bus)
		case 0x04:
			cpu.SLLV(instruction, bus)
		case 0x06:
			cpu.SRLV(instruction, bus)
		case 0x07:
			cpu.SRAV(instruction, bus)
		case 0x08:
			cpu.JR(instruction, bus)
		case 0x09:
			cpu.JALR(instruction, bus)
		case 0x0C:
			cpu.SYSCALL(instruction, bus)
		case 0x0D:
			cpu.BREAK(instruction, bus)
		case 0x10:
			cpu.MFHI(instruction, bus)
		case 0x11:
			cpu.MTHI(instruction, bus)
		case 0x12:
			cpu.MFLO(instruction, bus)
		case 0x13:
			cpu.MTLO(instruction, bus)
		case 0x18:
			cpu.MULT(instruction, bus)
		case 0x19:
			cpu.MULTU(instruction, bus)
		case 0x1A:
			cpu.DIV(instruction, bus)
		case 0x1B:
			cpu.DIVU(instruction, bus)
		case 0x20:
			cpu.ADD(instruction, bus)
		case 0x21:
			cpu.ADDU(instruction, bus)
		case 0x22:
			cpu.SUB(instruction, bus)
		case 0x23:
			cpu.SUBU(instruction, bus)
		case 0x24:
			cpu.AND(instruction, bus)
		case 0x25:
			cpu.OR(instruction, bus)
		case 0x26:
			cpu.XOR(instruction, bus)
		case 0x27:
			cpu.NOR(instruction, bus)
		case 0x2A:
			cpu.SLT(instruction, bus)
		case 0x2B:
			cpu.SLTU(instruction, bus)
		default:
			log.Fatalf("unknown special instruction: %02x", instruction.Function)
		}
	case 0x01:
		switch instruction.Rt {
		case 0x00:
			cpu.BLTZ(instruction, bus)
		case 0x01:
			cpu.BGEZ(instruction, bus)
		case 0x0A:
			cpu.BLTZAL(instruction, bus)
		case 0x0B:
			cpu.BGEZAL(instruction, bus)
		default:
			log.Fatalf("unknown bcondz instruction: %02x", instruction.Function)
		}
	case 0x02:
		cpu.J(instruction, bus)
	case 0x03:
		cpu.JAL(instruction, bus)
	case 0x04:
		cpu.BEQ(instruction, bus)
	case 0x05:
		cpu.BNE(instruction, bus)
	case 0x06:
		cpu.BLEZ(instruction, bus)
	case 0x07:
		cpu.BGTZ(instruction, bus)
	case 0x08:
		cpu.ADDI(instruction, bus)
	case 0x09:
		cpu.ADDIU(instruction, bus)
	case 0x0A:
		cpu.SLTI(instruction, bus)
	case 0x0B:
		cpu.SLTIU(instruction, bus)
	case 0x0C:
		cpu.ANDI(instruction, bus)
	case 0x0D:
		cpu.ORI(instruction, bus)
	case 0x0E:
		cpu.XORI(instruction, bus)
	case 0x0F:
		cpu.LUI(instruction, bus)
	case 0x10:
		switch instruction.Rs {
		case 0x0:
			cpu.MFC(instruction, bus)
		case 0x2:
			cpu.CFC(instruction, bus)
		case 0x4:
			cpu.MTC(instruction, bus)
		case 0x6:
			cpu.CTC(instruction, bus)
		default:
			log.Fatalf("unknown coprocessor opcode instruction: %02x", instruction.Opcode)
		}
	case 0x20:
		cpu.LB(instruction, bus)
	case 0x21:
		cpu.LH(instruction, bus)
	case 0x23:
		cpu.LW(instruction, bus)
	case 0x22:
		cpu.LWL(instruction, bus)
	case 0x24:
		cpu.LBU(instruction, bus)
	case 0x25:
		cpu.LHU(instruction, bus)
	case 0x26:
		cpu.LWR(instruction, bus)
	case 0x28:
		cpu.SB(instruction, bus)
	case 0x29:
		cpu.SH(instruction, bus)
	case 0x2A:
		cpu.SWL(instruction, bus)
	case 0x2B:
		cpu.SW(instruction, bus)
	case 0x2E:
		cpu.SWR(instruction, bus)
	default:
		log.Fatalf("unknown primary instruction: %02x", instruction.Opcode)
	}
}

func (cpu *CPU) Cycle(bus *Bus) {
	instruction := NewInstruction(bus.LoadWord(cpu.Pc))
	log.Printf("%08x %s", cpu.Pc, instruction.String())
	cpu.Pc = cpu.PcNext
	cpu.PcNext += 4
	cpu.Execute(instruction, bus)
}
