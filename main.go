package main

import (
	"log"
	"os"
)

type CPU struct {
	halted bool
}

func (cpu *CPU) Cycle(bus *Bus) {
	// instruction = cpu.fetchInstruction(bus)
	// cpu.execute(instruction)
}

func main() {
	bios, err := os.ReadFile("SCPH1001.bin")
	if err != nil {
		panic(err)
	}

	bus := NewBus(bios)

	log.Printf("%#v", NewInstruction(bus.ReadWord(0xBFC00000)))
}
