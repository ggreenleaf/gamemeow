package cpu

// The  Flags Register (lower 8 bits of the AF register)
const (
	FlagZero      = 1 << 7 // Bit 7 z
	FlagSubtract  = 1 << 6 // Bit 6 n
	FlagHalfCarry = 1 << 5 // Bit 5 h
	FlagCarry     = 1 << 4 // Bit 4 c
)

type Registers struct {
	A, f byte
	B, C byte
	D, E byte
	H, L byte
	SP   uint16
	PC   uint16
}

func (r *Registers) AF() uint16 {
	return (uint16(r.A) << 8) | uint16(r.f)
}

func (r *Registers) SetAF(value uint16) {
	r.A = byte(value >> 8)   // shifting right to get the top bits
	r.f = byte(value & 0xF0) // the flags only use the top 4 bits of F so we need to zero out the bottom 4 bits
}

func (r *Registers) BC() uint16 {
	return (uint16(r.B) << 8) | uint16(r.C)
}

func (r *Registers) SetBC(value uint16) {
	r.B = byte(value >> 8) // shifting right to get the top bits
	r.C = byte(value)      // cast truncates top keep bottom bits
}

func (r *Registers) DE() uint16 {
	return (uint16(r.D) << 8) | uint16(r.E)
}

func (r *Registers) SetDE(value uint16) {
	r.D = byte(value >> 8) // shifting right to get the top bits
	r.E = byte(value)      // cast truncates top keep bottom bits
}

func (r *Registers) HL() uint16 {
	return (uint16(r.H) << 8) | uint16(r.L)
}

func (r *Registers) SetHL(value uint16) {
	r.H = byte(value >> 8) // shifting right to get the top bits
	r.L = byte(value)      // cast truncates top keep bottom bits
}

func (r *Registers) FlagZ() bool {
	return r.f&FlagZero != 0
}

func (r *Registers) SetFlagZ(value bool) {
	r.setFlag(FlagZero, value)
}

func (r *Registers) FlagN() bool {
	return r.f&FlagSubtract != 0
}

func (r *Registers) SetFlagN(value bool) {
	r.setFlag(FlagSubtract, value)
}

func (r *Registers) FlagH() bool {
	return r.f&FlagHalfCarry != 0
}

func (r *Registers) SetFlagH(value bool) {
	r.setFlag(FlagHalfCarry, value)
}

func (r *Registers) FlagCy() bool {
	return r.f&FlagCarry != 0
}

func (r *Registers) SetFlagCy(value bool) {
	r.setFlag(FlagCarry, value)
}

// setFlag is a helper to set or clear a specific flag bit
func (r *Registers) setFlag(flag byte, value bool) {
	if value {
		r.f |= flag
	} else {
		r.f &^= flag
	}
}

func (r *Registers) FlagCyBit() uint16 {
	if r.FlagCy() {
		return 1
	}
	return 0
}
