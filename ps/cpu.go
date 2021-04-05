package ps

import (
	"log"
)

type CPU struct {
	// GPR is a General Purpose Registers.
	// The content of GPR[0] is always zero.
	// Attempts to alter the content of GPR[0] have no effect.
	// GPRNext is ...
	GPR, GPRNext []uint32

	// LoadDelaySlot emulates MIPS load delay
	LoadDelaySlot  uint32
	LoadDelayValue uint32

	// Pc is a program counter
	Pc, PcNext uint32

	// LO contains quotient
	// HI contains the remainder
	LO, HI uint32
}

func NewCPU() CPU {
	return CPU{
		GPR:            make([]uint32, 32),
		GPRNext:        make([]uint32, 32),
		LoadDelaySlot:  0,
		LoadDelayValue: 0,
		Pc:             0xBFC00000, // Bios start
		PcNext:         0xBFC00004,
	}
}

func (cpu *CPU) SetGPR(index, value uint32) {
	cpu.GPRNext[index] = value
}

func (cpu *CPU) GetGPR(index uint32) uint32 {
	return cpu.GPR[index]
}

func (cpu *CPU) SLL(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.GetGPR(instruction.Rt)<<instruction.ShiftAmount)
}

func (cpu *CPU) SRL(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.GetGPR(instruction.Rt)>>instruction.ShiftAmount)
}

func (cpu *CPU) SRA(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, uint32(int32(cpu.GetGPR(instruction.Rt))>>instruction.ShiftAmount))
}

func (cpu *CPU) SLLV(instruction Instruction) {
	s := cpu.GetGPR(instruction.Rs) & 0x1F
	cpu.SetGPR(instruction.Rd, cpu.GetGPR(instruction.Rt)<<s)
}

func (cpu *CPU) SRLV(instruction Instruction) {
	s := cpu.GetGPR(instruction.Rs) & 0x1F
	cpu.SetGPR(instruction.Rd, cpu.GetGPR(instruction.Rt)>>s)
}

func (cpu *CPU) SRAV(instruction Instruction) {
	s := cpu.GetGPR(instruction.Rs) & 0x1F
	cpu.SetGPR(instruction.Rd, uint32(int32(cpu.GetGPR(instruction.Rt))>>s))
}

func (cpu *CPU) JR(instruction Instruction) {
	cpu.PcNext = cpu.GetGPR(instruction.Rs)
}

func (cpu *CPU) JALR(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.Pc+8)
	cpu.PcNext = cpu.GetGPR(instruction.Rs)
}

func (cpu *CPU) SYSCALL(instruction Instruction, bus *Bus) {
	//TODO SystemCallException
}

func (cpu *CPU) BREAK(instruction Instruction, bus *Bus) {
	//TODO BreakpointException
}

func (cpu *CPU) MFHI(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.HI)
}

func (cpu *CPU) MTHI(instruction Instruction) {
	cpu.HI = cpu.GetGPR(instruction.Rs)
}

func (cpu *CPU) MFLO(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.LO)
}

func (cpu *CPU) MTLO(instruction Instruction) {
	cpu.LO = cpu.GetGPR(instruction.Rs)
}

func (cpu *CPU) MULT(instruction Instruction) {
	temp := int64(cpu.GetGPR(instruction.Rs)) * int64(cpu.GetGPR(instruction.Rt))
	cpu.LO = uint32(temp << 32)
	cpu.HI = uint32(temp >> 32)
}

func (cpu *CPU) MULTU(instruction Instruction) {
	temp := uint64(cpu.GetGPR(instruction.Rs)>>1) * uint64(cpu.GetGPR(instruction.Rt)>>1)
	cpu.LO = uint32(temp << 32)
	cpu.HI = uint32(temp >> 32)
}

func (cpu *CPU) DIV(instruction Instruction) {
	cpu.LO = cpu.GetGPR(instruction.Rs) / cpu.GetGPR(instruction.Rt)
	cpu.HI = cpu.GetGPR(instruction.Rs) % cpu.GetGPR(instruction.Rt)
}

func (cpu *CPU) DIVU(instruction Instruction) {
	cpu.LO = (cpu.GetGPR(instruction.Rs) >> 1) / (cpu.GetGPR(instruction.Rt) >> 1)
	cpu.HI = (cpu.GetGPR(instruction.Rs) >> 1) % (cpu.GetGPR(instruction.Rt) >> 1)
}

func (cpu *CPU) ADD(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, uint32(int32(cpu.GetGPR(instruction.Rs))+int32(cpu.GetGPR(instruction.Rt))))
}

func (cpu *CPU) ADDU(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.GetGPR(instruction.Rs)+cpu.GetGPR(instruction.Rt))
}

func (cpu *CPU) SUB(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, uint32(int32(cpu.GetGPR(instruction.Rs))-int32(cpu.GetGPR(instruction.Rt))))
}

func (cpu *CPU) SUBU(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.GetGPR(instruction.Rs)-cpu.GetGPR(instruction.Rt))
}

func (cpu *CPU) AND(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.GetGPR(instruction.Rs)&cpu.GetGPR(instruction.Rt))
}

func (cpu *CPU) OR(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.GetGPR(instruction.Rs)|cpu.GetGPR(instruction.Rt))
}

func (cpu *CPU) XOR(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, cpu.GetGPR(instruction.Rs)^cpu.GetGPR(instruction.Rt))
}

func (cpu *CPU) NOR(instruction Instruction) {
	cpu.SetGPR(instruction.Rd, 0xFFFFFFFF^(cpu.GetGPR(instruction.Rs)|cpu.GetGPR(instruction.Rt)))
}

func (cpu *CPU) SLT(instruction Instruction) {
	if cpu.GetGPR(instruction.Rs) < cpu.GetGPR(instruction.Rt) {
		cpu.SetGPR(instruction.Rd, 1)
	} else {
		cpu.SetGPR(instruction.Rd, 0)
	}
}

func (cpu *CPU) SLTU(instruction Instruction) {
	if (cpu.GetGPR(instruction.Rs) >> 1) < (cpu.GetGPR(instruction.Rt) >> 1) {
		cpu.SetGPR(instruction.Rd, 1)
	} else {
		cpu.SetGPR(instruction.Rd, 0)
	}
}

func (cpu *CPU) BLTZ(instruction Instruction) {
	address := cpu.PcNext + (instruction.Imm16sx << 2)
	if (cpu.GetGPR(instruction.Rs) >> 31) == 1 {
		cpu.PcNext = address
	}
}

func (cpu *CPU) BGEZ(instruction Instruction) {
	address := cpu.PcNext + (instruction.Imm16sx << 2)
	if (cpu.GetGPR(instruction.Rs) >> 31) == 0 {
		cpu.PcNext = address
	}
}

func (cpu *CPU) BLTZAL(instruction Instruction) {
	address := cpu.PcNext + (instruction.Imm16sx << 2)
	if (cpu.GetGPR(instruction.Rs) >> 31) == 1 {
		cpu.PcNext = address
	}
	cpu.SetGPR(31, cpu.Pc+8)
}

func (cpu *CPU) BGEZAL(instruction Instruction) {
	address := cpu.PcNext + (instruction.Imm16sx << 2)
	if (cpu.GetGPR(instruction.Rs) >> 31) == 0 {
		cpu.PcNext = address
	}
	cpu.SetGPR(31, cpu.Pc+8)
}

func (cpu *CPU) J(instruction Instruction) {
	cpu.PcNext = cpu.Pc&0xF0000000 | (instruction.Address << 2)
}

func (cpu *CPU) JAL(instruction Instruction) {
	cpu.SetGPR(31, cpu.Pc+8)
	cpu.PcNext = cpu.Pc&0xF0000000 | (instruction.Address << 2)
}

func (cpu *CPU) BEQ(instruction Instruction) {
	address := cpu.PcNext + (instruction.Imm16sx << 2)
	if cpu.GetGPR(instruction.Rs) == cpu.GetGPR(instruction.Rt) {
		cpu.PcNext = address
	}
}

func (cpu *CPU) BNE(instruction Instruction) {
	address := cpu.PcNext + (instruction.Imm16sx << 2)
	if cpu.GetGPR(instruction.Rs) != cpu.GetGPR(instruction.Rt) {
		cpu.PcNext = address
	}
}

func (cpu *CPU) BLEZ(instruction Instruction) {
	address := cpu.PcNext + (instruction.Imm16sx << 2)
	if (cpu.GetGPR(instruction.Rs)>>31) == 1 || cpu.GetGPR(instruction.Rs) == 0 {
		cpu.PcNext = address
	}
}

func (cpu *CPU) BGTZ(instruction Instruction) {
	address := cpu.PcNext + (instruction.Imm16sx << 2)
	if (cpu.GetGPR(instruction.Rs)>>31) == 0 && cpu.GetGPR(instruction.Rs) != 0 {
		cpu.PcNext = address
	}
}

func (cpu *CPU) ADDI(instruction Instruction) {
	cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rs)+instruction.Imm16sx)
}

func (cpu *CPU) ADDIU(instruction Instruction) {
	//TODO 32-bit-overflow
	cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rs)+instruction.Imm16sx)
}

func (cpu *CPU) SLTI(instruction Instruction) {
	if cpu.GetGPR(instruction.Rs) < instruction.Imm16sx {
		cpu.SetGPR(instruction.Rt, 1)
	} else {
		cpu.SetGPR(instruction.Rt, 0)
	}
}

func (cpu *CPU) SLTIU(instruction Instruction) {
	if (cpu.GetGPR(instruction.Rs) >> 1) < instruction.Imm16sx {
		cpu.SetGPR(instruction.Rt, 1)
	} else {
		cpu.SetGPR(instruction.Rt, 0)
	}
}

func (cpu *CPU) ANDI(instruction Instruction) {
	cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rs)&instruction.Imm16)
}

func (cpu *CPU) ORI(instruction Instruction) {
	cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rs)|instruction.Imm16)
}

func (cpu *CPU) XORI(instruction Instruction) {
	cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rs)^instruction.Imm16)
}

func (cpu *CPU) LUI(instruction Instruction) {
	cpu.SetGPR(instruction.Rt, instruction.Imm16<<16)
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
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	value := bus.LoadByte(address)
	cpu.LoadDelaySlot = instruction.Rt
	cpu.LoadDelayValue = uint32(int8(value))
}

func (cpu *CPU) LH(instruction Instruction, bus *Bus) {
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	value := bus.LoadHalfword(address)
	cpu.LoadDelaySlot = instruction.Rt
	cpu.LoadDelayValue = uint32(int16(value))
}

func (cpu *CPU) LW(instruction Instruction, bus *Bus) {
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	value := bus.LoadWord(address)
	cpu.LoadDelaySlot = instruction.Rt
	cpu.LoadDelayValue = value
}

func (cpu *CPU) LWL(instruction Instruction, bus *Bus) {
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	var temp uint32 = 0

	switch address % 4 {
	case 0:
		cpu.SetGPR(instruction.Rt, bus.LoadWord(address/4))
		break
	case 1:
		temp = bus.LoadWord(address/4) & 0xFFFFFF << 8
		cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rt)&0xFF|temp)
		break
	case 2:
		temp = bus.LoadWord(address/4) & 0xFFFF << 16
		cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rt)&0xFFFF|temp)
		break
	case 3:
		temp = bus.LoadWord(address/4) & 0xFF << 24
		cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rt)&0xFFFFFF|temp)
		break
	}
}

func (cpu *CPU) LBU(instruction Instruction, bus *Bus) {
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	value := bus.LoadByte(address)
	cpu.LoadDelaySlot = instruction.Rt
	cpu.LoadDelayValue = uint32(value)
}

func (cpu *CPU) LHU(instruction Instruction, bus *Bus) {
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	value := bus.LoadHalfword(address)
	cpu.LoadDelaySlot = instruction.Rt
	cpu.LoadDelayValue = uint32(value)
}

func (cpu *CPU) LWR(instruction Instruction, bus *Bus) {
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	var temp uint32 = 0

	switch address % 4 {
	case 0:
		cpu.SetGPR(instruction.Rt, bus.LoadWord(address/4))
		break
	case 1:
		temp = (bus.LoadWord(address/4) & 0xFF000000 >> 24) & 0xFF
		cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rt)&0xFFFFFF00|temp)
		break
	case 2:
		temp = (bus.LoadWord(address/4) & 0xFFFF0000 >> 16) & 0xFFFF
		cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rt)&0xFFFF0000|temp)
		break
	case 3:
		temp = (bus.LoadWord(address/4) & 0xFFFFFF00 >> 8) & 0xFFFFFF
		cpu.SetGPR(instruction.Rt, cpu.GetGPR(instruction.Rt)&0xFF000000|temp)
		break
	}
}

func (cpu *CPU) SB(instruction Instruction, bus *Bus) {
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	bus.StoreByte(address, uint8(cpu.GetGPR(instruction.Rt)&0xFF))
}

func (cpu *CPU) SH(instruction Instruction, bus *Bus) {
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	bus.StoreHalfword(address, uint16(cpu.GetGPR(instruction.Rt)&0xFFFF))
}

func (cpu *CPU) SWL(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) SW(instruction Instruction, bus *Bus) {
	address := instruction.Imm16sx + cpu.GetGPR(instruction.Rs)
	bus.StoreWord(address, cpu.GetGPR(instruction.Rt))
}

func (cpu *CPU) SWR(instruction Instruction, bus *Bus) {
	//TODO Later
}

func (cpu *CPU) Execute(instruction Instruction, bus *Bus) {
	switch instruction.Opcode {
	case 0x00:
		switch instruction.Function {
		case 0x00:
			cpu.SLL(instruction)
		case 0x02:
			cpu.SRL(instruction)
		case 0x03:
			cpu.SRA(instruction)
		case 0x04:
			cpu.SLLV(instruction)
		case 0x06:
			cpu.SRLV(instruction)
		case 0x07:
			cpu.SRAV(instruction)
		case 0x08:
			cpu.JR(instruction)
		case 0x09:
			cpu.JALR(instruction)
		//case 0x0C:
		//	cpu.SYSCALL(instruction, bus)
		//case 0x0D:
		//	cpu.BREAK(instruction, bus)
		case 0x10:
			cpu.MFHI(instruction)
		case 0x11:
			cpu.MTHI(instruction)
		case 0x12:
			cpu.MFLO(instruction)
		case 0x13:
			cpu.MTLO(instruction)
		case 0x18:
			cpu.MULT(instruction)
		case 0x19:
			cpu.MULTU(instruction)
		case 0x1A:
			cpu.DIV(instruction)
		case 0x1B:
			cpu.DIVU(instruction)
		case 0x20:
			cpu.ADD(instruction)
		case 0x21:
			cpu.ADDU(instruction)
		case 0x22:
			cpu.SUB(instruction)
		case 0x23:
			cpu.SUBU(instruction)
		case 0x24:
			cpu.AND(instruction)
		case 0x25:
			cpu.OR(instruction)
		case 0x26:
			cpu.XOR(instruction)
		case 0x27:
			cpu.NOR(instruction)
		case 0x2A:
			cpu.SLT(instruction)
		case 0x2B:
			cpu.SLTU(instruction)
		default:
			log.Fatalf("unknown special instruction: %02x", instruction.Function)
		}
	case 0x01:
		switch instruction.Rt {
		case 0x00:
			cpu.BLTZ(instruction)
		case 0x01:
			cpu.BGEZ(instruction)
		case 0x0A:
			cpu.BLTZAL(instruction)
		case 0x0B:
			cpu.BGEZAL(instruction)
		default:
			log.Fatalf("unknown bcondz instruction: %02x", instruction.Function)
		}
	case 0x02:
		cpu.J(instruction)
	case 0x03:
		cpu.JAL(instruction)
	case 0x04:
		cpu.BEQ(instruction)
	case 0x05:
		cpu.BNE(instruction)
	case 0x06:
		cpu.BLEZ(instruction)
	case 0x07:
		cpu.BGTZ(instruction)
	case 0x08:
		cpu.ADDI(instruction)
	case 0x09:
		cpu.ADDIU(instruction)
	case 0x0A:
		cpu.SLTI(instruction)
	case 0x0B:
		cpu.SLTIU(instruction)
	case 0x0C:
		cpu.ANDI(instruction)
	case 0x0D:
		cpu.ORI(instruction)
	case 0x0E:
		cpu.XORI(instruction)
	case 0x0F:
		cpu.LUI(instruction)
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
			log.Fatalf("unknown coprocessor opcode instruction: %02x", instruction.Rs)
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
	//case 0x2A:
	//	cpu.SWL(instruction, bus)
	case 0x2B:
		cpu.SW(instruction, bus)
	//case 0x2E:
	//	cpu.SWR(instruction, bus)
	default:
		log.Fatalf("unknown primary instruction: %02x", instruction.Opcode)
	}
}

func (cpu *CPU) Cycle(bus *Bus) {
	instruction := NewInstruction(bus.LoadWord(cpu.Pc))
	log.Printf("%08x %s", cpu.Pc, instruction)
	cpu.PcNext += 4
	cpu.Pc = cpu.PcNext
	cpu.SetGPR(cpu.LoadDelaySlot, cpu.LoadDelayValue)
	cpu.LoadDelaySlot = 0
	cpu.LoadDelayValue = 0

	cpu.Execute(instruction, bus)
	copy(cpu.GPR, cpu.GPRNext)
}
