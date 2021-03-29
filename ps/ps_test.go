package ps

import (
	"testing"
)

func assertEqual(t *testing.T, actual interface{}, expected interface{}) {
	t.Helper()
	if actual != expected {
		t.Fatalf("assertion failed: %v != %v", actual, expected)
	}
}

func TestInstructionDecoding(t *testing.T) {
	instruction := NewInstruction(0x1234ABCD)
	assertEqual(t, instruction, Instruction{
		Opcode: 0x4,

		Rs: 0x11,
		Rt: 0x14,
		Rd: 0x15,

		ShiftAmount: 0xF,
		Function:    0xD,

		Imm16:   0xABCD,
		Imm16sx: 0xFFFFABCD,

		Address: 0x0234ABCD,
	})
}

func TestLoadStore(t *testing.T) {
	bus := NewBus(make([]byte, 512*1024))

	var addr uint32 = 0xBFC01234

	bus.StoreWord(addr, 0xABCD1234)
	bus.StoreHalfword(addr+4, 0x5A51)
	bus.StoreByte(addr+6, 0xEA)

	assertEqual(t, bus.LoadWord(addr), uint32(0xABCD1234))
	assertEqual(t, bus.LoadHalfword(addr+4), uint16(0x5A51))
	assertEqual(t, bus.LoadByte(addr+6), uint8(0xEA))

	assertEqual(t, bus.LoadWord(addr+1), uint32(0x51ABCD12))
	assertEqual(t, bus.LoadHalfword(addr+5), uint16(0xEA5A))
	assertEqual(t, bus.LoadByte(addr), uint8(0x34))
}
