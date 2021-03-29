package ps

import "log"

type Vector3 struct {
	X, Y, Z int32
}

var ZeroVector3 Vector3 = Vector3{
	X: 0,
	Y: 0,
	Z: 0,
}

type ColorU8 struct {
	X, Y, Z uint8
}

// GTE emulates PlayStation's Geometry Transformation Engine, a coprocessor
// designed for linear algebra
type GTE struct {
	V0, V1, V2       Vector3
	RGBC             uint8
	OTZ              uint16
	IR0              int16
	IR               Vector3
	SX0, SY0         int16
	SX1, SY1         int16
	SX2, SY2         int16
	SXP, SYP         int16
	SZ0              uint16
	SZ1              uint16
	SZ2              uint16
	SZ3              uint16
	RGB0, RGB1, RGB2 ColorU8
	MAC0             int32
	MAC1, MAC2, MAC3 int32
	IRGB, ORGB       uint16
	LZCS, LZCR       int32

	RT1, RT2, RT3 Vector3
	TR            Vector3
	L1, L2, L3    Vector3
	BK            Vector3
	LR, LG, LB    Vector3
	FC            Vector3
	OFX, OFY      int32
	H             uint16
	DQA           int16
	DQB           int32
	ZSF3, ZSF4    int16
	FLAG          uint32
}

func NewGTE() GTE {
	return GTE{}
}

type GTEParameters struct {
	SF               bool
	MVMVAMatrix      uint32
	MVMVAVector      uint32
	MVMVATranslation uint32
	LM               bool
}

func NewGTEParameters(instruction uint32) GTEParameters {
	return GTEParameters{
		SF:               (instruction >> 19) == 1,
		MVMVAMatrix:      (instruction >> 17) & 0x3,
		MVMVAVector:      (instruction >> 15) & 0x3,
		MVMVATranslation: (instruction >> 13) & 0x3,
		LM:               (instruction >> 10) == 1,
	}
}

func (g *GTE) MVMVA(p GTEParameters) {
	m1 := []Vector3{g.RT1, g.L1, g.LR, ZeroVector3}[p.MVMVAMatrix]
	m2 := []Vector3{g.RT2, g.L2, g.LG, ZeroVector3}[p.MVMVAMatrix]
	m3 := []Vector3{g.RT3, g.L3, g.LB, ZeroVector3}[p.MVMVAMatrix]
	v := []Vector3{g.V0, g.V1, g.V2, g.IR}[p.MVMVAVector]
	t := []Vector3{g.TR, g.BK, g.FC, ZeroVector3}[p.MVMVATranslation]

	shift := 0
	if p.SF {
		shift = 12
	}
	g.MAC1 = (t.X*0x1000 + m1.X*v.X + m1.Y*v.X + m1.Z*v.X) >> shift
	g.MAC2 = (t.Y*0x1000 + m2.X*v.Y + m2.Y*v.Y + m2.Z*v.Y) >> shift
	g.MAC3 = (t.Z*0x1000 + m3.X*v.Z + m3.Y*v.Z + m3.Z*v.Z) >> shift

	g.IR.X = g.MAC1
	g.IR.Y = g.MAC2
	g.IR.Z = g.MAC3
}

func (g *GTE) Execute(instruction Instruction) {
	log.Printf("Executing GTE instruction: %s", GTEOpcodeName(instruction.Function))
	parameters := NewGTEParameters(instruction.Raw)
	switch instruction.Function {
	case 0x01:
		// g.RTPS(parameters)
	case 0x06:
		// g.NCLIP()
	case 0x0C:
		// g.OP()
	case 0x10:
		// g.DPCS()
	case 0x11:
		// g.INTPL()
	case 0x12:
		g.MVMVA(parameters)
	case 0x13:
		// g.NCDS()
	case 0x14:
		// g.CDP()
	case 0x16:
		// g.NCDT()
	case 0x1B:
		// g.NCCS()
	case 0x1C:
		// g.CC()
	case 0x1E:
		// g.NCS()
	case 0x20:
		// g.NCT()
	case 0x28:
		// g.SQR()
	case 0x29:
		// g.DCPL()
	case 0x2A:
		// g.DPCT()
	case 0x2D:
		// g.AVSZ3()
	case 0x2E:
		// g.AVSZ4()
	case 0x30:
		// g.RTPT()
	case 0x3D:
		// g.GPF()
	case 0x3E:
		// g.GPL()
	case 0x3F:
		// g.NCCT()
	default:
		log.Fatalf("unknown GTE instruction %02x", instruction.Function)
	}
}

func GTEOpcodeName(opcode uint32) string {
	names := map[uint32]string{
		0x01: "RTPS",
		0x06: "NCLIP",
		0x0C: "OP",
		0x10: "DPCS",
		0x11: "INTPL",
		0x12: "MVMVA",
		0x13: "NCDS",
		0x14: "CDP",
		0x16: "NCDT",
		0x1B: "NCCS",
		0x1C: "CC",
		0x1E: "NCS",
		0x20: "NCT",
		0x28: "SQR",
		0x29: "DCPL",
		0x2A: "DPCT",
		0x2D: "AVSZ3",
		0x2E: "AVSZ4",
		0x30: "RTPT",
		0x3D: "GPF",
		0x3E: "GPL",
		0x3F: "NCCT",
	}

	if name, ok := names[opcode]; ok {
		return name
	} else {
		return "N/A"
	}
}
