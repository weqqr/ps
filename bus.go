package main

import (
	"log"
)

// Memory map
const (
	BiosStart = 0x1FC00000
	BiosSize  = 512 * 1024
)

type Bus struct {
	bios []byte
}

func NewBus(bios []byte) Bus {
	if len(bios) != 512*1024 {
		log.Fatal("Error: BIOS size must be exactly 512 KiB")
	}

	return Bus{
		bios: bios,
	}
}

func inRange(value, start, size uint32) bool {
	return value >= start && value < start+size
}

func (bus *Bus) Map(address uint32) (uint32, []byte) {
	// Mask segment
	address = address & 0x1FFFFFFF

	switch {
	case inRange(address, BiosStart, BiosSize):
		return address - BiosStart, bus.bios
	default:
		log.Fatalf("unknown memory region at address %x", address)
	}

	return 0, []byte{}
}

func (bus *Bus) ReadByte(address uint32) uint8 {
	address, data := bus.Map(address)
	return data[address]
}

func (bus *Bus) ReadHalfword(address uint32) uint16 {
	address, data := bus.Map(address)
	a := uint16(data[address+1])
	b := uint16(data[address])
	return (a << 8) | b
}

func (bus *Bus) ReadWord(address uint32) uint32 {
	address, data := bus.Map(address)
	a := uint32(data[address+3])
	b := uint32(data[address+2])
	c := uint32(data[address+1])
	d := uint32(data[address])
	return (a << 24) | (b << 16) | (c << 8) | d
}
