package main

import (
	"os"
	"github.com/weqqr/ps/ps"
)


func main() {
	bios, err := os.ReadFile("SCPH1001.bin")
	if err != nil {
		panic(err)
	}

	bus := ps.NewBus(bios)
	cpu := ps.NewCPU()

	for {
		cpu.Cycle(&bus)
	}
}
