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

// GTEParameters contains flags and indices required by certain operations
type GTEParameters struct {
	// Shift indicates amount of bits fractional values should be shifted right
	Shift uint32

	// LM indicates if results should be saturated (i. e. clamped to 0..0x7FFF)
	LM bool

	// MVMVAMatrix, MVMVAVector and MVMVATranslation are 2-bit values that
	// select Matrix, Vector and Translation vector for MVMVA (matrix-vector
	// multiplication with translation vector addition)
	MVMVAMatrix      uint32
	MVMVAVector      uint32
	MVMVATranslation uint32
}

func NewGTEParameters(instruction uint32) GTEParameters {
	var shift uint32 = 0
	if ((instruction >> 19) & 0x1) == 1 {
		shift = 12
	}

	return GTEParameters{
		Shift:            shift,
		LM:               ((instruction >> 10) & 0x1) == 1,
		MVMVAMatrix:      (instruction >> 17) & 0x3,
		MVMVAVector:      (instruction >> 15) & 0x3,
		MVMVATranslation: (instruction >> 13) & 0x3,
	}
}

func saturate(value int16) int16 {
	if value < 0 {
		return 0
	}
	return value
}

func (g *GTE) DCPL() {
	/*
	[MAC1, MAC2, MAC3] = [R * IR1, G * IR2, B * IR3] SHL 4
	[MAC1, MAC2, MAC3] = MAC + (FC - MAC) * IRO
	[MAC1, MAC2, MAC3] = [MAC1, MAC2, MAC3] SAR (sf * 12)
	Color FIFO = [MAC1 / 16, MAC2 / 16, MAC3 / 16, CODE], [IR1, IR2, IR3] = [MAC1, MAC2, MAC3]
	*/
	[g.MAC1, g.MAC2, g.MAC3] = [g.R * g.IR1, g.G * g.IR2, g.B * g.IR3] << 4
	[g.MAC1, g.MAC2, g.MAC3] = g.MAC + (g.FC - g.MAC) * g.IR0
	[g.MAC1, g.MAC2, g.MAC3] = [MAC1, MAC2, MAC3] >> (g.sf * 12)
	Color FIFO = [g.MAC1 / 16, g.MAC2 / 16, MAC3 / 16, g.CODE]
	[g.IR1, g.IR2, g.IR3] = [g.MAC1, g.MAC2, g.MAC3]
}

func (g *GTE) DPCS() {
	[g.MAC1, g.MAC2, g.MAC3] = [g.R, g.G, g.B] >> 16
	[g.MAC1, g.MAC2, g.MAC3] = g.MAC + (g.FC - g.MAC) * g.IR0
	[g.MAC1, g.MAC2, g.MAC3] = [MAC1, MAC2, MAC3] >> (g.sf * 12)
	Color FIFO = [g.MAC1 / 16, g.MAC2 / 16, MAC3 / 16, g.CODE]
	[g.IR1, g.IR2, g.IR3] = [g.MAC1, g.MAC2, g.MAC3]
}

func (g *GTE) DPCT() {
	[g.MAC1, g.MAC2, g.MAC3] = [g.R, g.G, g.B] >> 16
	[g.MAC1, g.MAC2, g.MAC3] = g.MAC + (g.FC - g.MAC) * g.IR0
	[g.MAC1, g.MAC2, g.MAC3] = [MAC1, MAC2, MAC3] >> (g.sf * 12)
	Color FIFO = [g.MAC1 / 16, g.MAC2 / 16, MAC3 / 16, g.CODE]
	[g.IR1, g.IR2, g.IR3] = [g.MAC1, g.MAC2, g.MAC3]
}




func (g *GTE) NCLIP() {
	g.MAC0 = int32(g.SX0*g.SY1 + g.SX1*g.SY2 + g.SX2*g.SY0 - g.SX0*g.SY2 - g.SX1*g.SY0 - g.SX2*g.SY1)
}

func (g *GTE) OP(p GTEParameters) {
	D1 := g.RT1.X
	D2 := g.RT2.Y
	D3 := g.RT3.Z

	g.MAC1 = (g.IR.Z*D2 - g.IR.Y*D3) >> p.Shift
	g.MAC2 = (g.IR.X*D3 - g.IR.Z*D1) >> p.Shift
	g.MAC3 = (g.IR.Y*D1 - g.IR.X*D2) >> p.Shift

	g.IR.X, g.IR.Y, g.IR.Z = g.MAC1, g.MAC2, g.MAC3
}

func (g *GTE) MVMVA(p GTEParameters) {
	m1 := []Vector3{g.RT1, g.L1, g.LR, ZeroVector3}[p.MVMVAMatrix]
	m2 := []Vector3{g.RT2, g.L2, g.LG, ZeroVector3}[p.MVMVAMatrix]
	m3 := []Vector3{g.RT3, g.L3, g.LB, ZeroVector3}[p.MVMVAMatrix]
	v := []Vector3{g.V0, g.V1, g.V2, g.IR}[p.MVMVAVector]
	t := []Vector3{g.TR, g.BK, g.FC, ZeroVector3}[p.MVMVATranslation]

	g.MAC1 = (t.X*0x1000 + m1.X*v.X + m1.Y*v.X + m1.Z*v.X) >> p.Shift
	g.MAC2 = (t.Y*0x1000 + m2.X*v.Y + m2.Y*v.Y + m2.Z*v.Y) >> p.Shift
	g.MAC3 = (t.Z*0x1000 + m3.X*v.Z + m3.Y*v.Z + m3.Z*v.Z) >> p.Shift

	g.IR.X, g.IR.Y, g.IR.Z = g.MAC1, g.MAC2, g.MAC3
}

func (g *GTE) SQR(p GTEParameters) {
	g.MAC1 = (g.TR.X * g.TR.X) >> p.Shift
	g.MAC2 = (g.TR.Y * g.TR.Y) >> p.Shift
	g.MAC3 = (g.TR.Z * g.TR.Z) >> p.Shift

	g.IR.X, g.IR.Y, g.IR.Z = g.MAC1, g.MAC2, g.MAC3
}

func (g *GTE) AVSZ3() {
	g.MAC0 = int32(g.ZSF3 * int16(g.SZ1+g.SZ2+g.SZ3))
	g.OTZ = uint16(saturate(int16(g.MAC0 / 0x1000)))
}

func (g *GTE) AVSZ4() {
	g.MAC0 = int32(g.ZSF4 * int16(g.SZ0+g.SZ1+g.SZ2+g.SZ3))
	g.OTZ = uint16(saturate(int16(g.MAC0 / 0x1000)))
}

func (g *GTE) Execute(instruction Instruction) {
	log.Printf("Executing GTE instruction: %s", GTEOpcodeName(instruction.Function))
	parameters := NewGTEParameters(instruction.Raw)
	switch instruction.Function {
	case 0x01:
		// g.RTPS(parameters)
	case 0x06:
		g.NCLIP()
	case 0x0C:
		g.OP(parameters)
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
		g.SQR(parameters)
	case 0x29:
		// g.DCPL()
	case 0x2A:
		// g.DPCT()
	case 0x2D:
		g.AVSZ3()
	case 0x2E:
		g.AVSZ4()
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
