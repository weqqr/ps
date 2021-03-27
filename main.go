package main

import (
	"log"
	"os"
)

type CPU struct {
	Halted bool
	GPR    []uint32
	Pc     uint32
}

func NewCPU() CPU {
	return CPU{
		Halted: false,
		GPR:    make([]uint32, 32),
		Pc:     0xBFC00000, // Bios start
	}
}

func (cpu *CPU) Execute(instruction Instruction) {
	switch instruction.Opcode {
	default:
		log.Fatalf("unknown instruction: %02x", instruction.Opcode)
	}
}

func (cpu *CPU) Cycle(bus *Bus) {
	instruction := NewInstruction(bus.ReadWord(cpu.Pc))
	cpu.Pc += 4
	cpu.Execute(instruction)
}

func main() {
	bios, err := os.ReadFile("SCPH1001.bin")
	if err != nil {
		panic(err)
	}

	bus := NewBus(bios)
	cpu := NewCPU()

	for !cpu.Halted {
		cpu.Cycle(&bus)
	}

}
