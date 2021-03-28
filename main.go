package main

import (
	"log"
	"os"
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

func (cpu *CPU) LUI(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rt] = instruction.Imm16 << 16
}

func (cpu *CPU) ORI(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rt] = cpu.GPR[instruction.Rs] | instruction.Imm16
}

func (cpu *CPU) SW(instruction Instruction, bus *Bus) {
	address := instruction.Imm16 | cpu.GPR[instruction.Rs]
	bus.StoreWord(address, cpu.GPR[instruction.Rt])
}

func (cpu *CPU) SLL(instruction Instruction, bus *Bus) {
	cpu.GPR[instruction.Rd] = cpu.GPR[instruction.Rt] << instruction.ShiftAmount
}

func (cpu *CPU) ADDIU(instruction Instruction, bus *Bus) {
	//TODO 32-bit-overflow
	cpu.GPR[instruction.Rt] = cpu.GPR[instruction.Rs] + instruction.Imm16sx
}

func (cpu *CPU) J(instruction Instruction, bus *Bus) {
	cpu.PcNext = cpu.Pc&0xF0000000 | (instruction.Address << 2)
}

func (cpu *CPU) Execute(instruction Instruction, bus *Bus) {
	switch instruction.Opcode {
	case 0x00:
		switch instruction.Function {
		case 0x00:
			cpu.SLL(instruction, bus)
		default:
			log.Fatalf("unknown special instruction: %02x", instruction.Function)
		}
	case 0x02:
		cpu.J(instruction, bus)
	case 0x09:
		cpu.ADDIU(instruction, bus)
	case 0x0d:
		cpu.ORI(instruction, bus)
	case 0x0F:
		cpu.LUI(instruction, bus)
	case 0x2b:
		cpu.SW(instruction, bus)
	default:
		log.Fatalf("unknown instruction: %02x", instruction.Opcode)
	}
}

func (cpu *CPU) Cycle(bus *Bus) {
	instruction := NewInstruction(bus.LoadWord(cpu.Pc))
	log.Printf("%08x %s", cpu.Pc, instruction.String())
	cpu.Pc = cpu.PcNext
	cpu.PcNext += 4
	cpu.Execute(instruction, bus)
}

func main() {
	bios, err := os.ReadFile("SCPH1001.bin")
	if err != nil {
		panic(err)
	}

	bus := NewBus(bios)
	cpu := NewCPU()

	for {
		cpu.Cycle(&bus)
	}
}
