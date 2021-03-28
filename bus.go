package main

import (
	"log"
)

// Memory map
// http://problemkaputt.de/psx-spx.htm#memorymap
const (
	MainRAM               = 0x00000000
	FirstExpansionRegion  = 0x1F000000
	Scratchpad            = 0x1F800000
	IOPorts               = 0x1F801000
	SecondExpansionRegion = 0x1F801000
	ThirdExpansionRegion  = 0x1F802000
	BIOSAddress           = 0x1FC00000
	CacheControl          = 0xFFFE0000

	MainRAMSize               = 2048 * 1024
	FirstExpansionRegionSize  = 8192 * 1024
	ScratchpadSize            = 1 * 1024
	IOPortsSize               = 8 * 1024
	SecondExpansionRegionSize = 8 * 1024
	ThirdExpansionRegionSize  = 2048 * 1024
	BIOSSize                  = 512 * 1024
	CacheControlSize          = 512
)

type Bus struct {
	mainRAM               []byte
	firstExpansionRegion  []byte
	scratchpad            []byte
	ioPorts               []byte
	secondExpansionRegion []byte
	thirdExpansionRegion  []byte
	bios                  []byte
	cacheControl          []byte
}

func NewBus(bios []byte) Bus {
	if len(bios) != BIOSSize {
		log.Fatal("Error: BIOSAddress size must be exactly 512 KiB")
	}

	return Bus{
		mainRAM:               make([]byte, MainRAMSize),
		firstExpansionRegion:  make([]byte, FirstExpansionRegionSize),
		scratchpad:            make([]byte, ScratchpadSize),
		ioPorts:               make([]byte, IOPortsSize),
		secondExpansionRegion: make([]byte, SecondExpansionRegionSize),
		thirdExpansionRegion:  make([]byte, ThirdExpansionRegionSize),
		bios:                  bios,
		cacheControl:          make([]byte, CacheControlSize),
	}
}

func inRange(value, start, size uint32) bool {
	return value >= start && value < start+size
}

func (bus *Bus) Map(address uint32) (uint32, []byte) {
	if inRange(address, CacheControl, CacheControlSize) {
		return address - CacheControl, bus.cacheControl
	}

	// Mask segment
	address = address & 0x1FFFFFFF

	switch {
	case inRange(address, MainRAM, MainRAMSize):
		return address - MainRAM, bus.mainRAM
	case inRange(address, FirstExpansionRegion, FirstExpansionRegionSize):
		return address - FirstExpansionRegion, bus.firstExpansionRegion
	case inRange(address, Scratchpad, ScratchpadSize):
		return address - Scratchpad, bus.scratchpad
	case inRange(address, IOPorts, IOPortsSize):
		return address - IOPorts, bus.ioPorts
	case inRange(address, SecondExpansionRegion, SecondExpansionRegionSize):
		return address - SecondExpansionRegion, bus.secondExpansionRegion
	case inRange(address, ThirdExpansionRegion, ThirdExpansionRegionSize):
		return address - ThirdExpansionRegion, bus.thirdExpansionRegion
	case inRange(address, BIOSAddress, BIOSSize):
		return address - BIOSAddress, bus.bios
	default:
		log.Fatalf("unknown memory region at address %x", address)
	}

	return 0, []byte{}
}

func (bus *Bus) LoadByte(address uint32) uint8 {
	address, data := bus.Map(address)
	return data[address]
}

func (bus *Bus) LoadHalfword(address uint32) uint16 {
	address, data := bus.Map(address)
	a := uint16(data[address+1])
	b := uint16(data[address])
	return (a << 8) | b
}

func (bus *Bus) LoadWord(address uint32) uint32 {
	address, data := bus.Map(address)
	a := uint32(data[address+3])
	b := uint32(data[address+2])
	c := uint32(data[address+1])
	d := uint32(data[address])
	return (a << 24) | (b << 16) | (c << 8) | d
}

func (bus *Bus) StoreByte(address uint32, value uint8) {
	address, data := bus.Map(address)
	data[address] = value
}

func (bus *Bus) StoreHalfword(address uint32, value uint16) {
	address, data := bus.Map(address)
	data[address+1] = uint8(value >> 8)
	data[address] = uint8(value)
}

func (bus *Bus) StoreWord(address uint32, value uint32) {
	address, data := bus.Map(address)

	data[address+3] = uint8(value >> 24)
	data[address+2] = uint8(value >> 16)
	data[address+1] = uint8(value >> 8)
	data[address] = uint8(value)
}
