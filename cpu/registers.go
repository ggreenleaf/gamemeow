package cpu

// The  Flags Register (lower 8 bits of the AF register)
const (
	FlagZero      = 1 << 7 // Bit 7 z
	FlagSubtract  = 1 << 6 // Bit 6 n
	FlagHalfCarry = 1 << 5 // Bit 5 h
	FlagCarry     = 1 << 4 // Bit 4 c
)

type Registers struct {
	a, f byte
	b, c byte
	d, e byte
	h, l byte
	sp   uint16
	pc   uint16
}

func (r *Registers) AF() uint16 {
	return (uint16(r.a) << 8) | uint16(r.f)
}

func (r *Registers) SetAF(value uint16) {
	r.a = byte(value >> 8)   // shifting right to get the top bits
	r.f = byte(value & 0xF0) // the flags only use the top 4 bits of F so we need to zero out the bottom 4 bits
}

func (r *Registers) BC() uint16 {
	return (uint16(r.b) << 8) | uint16(r.c)
}

func (r *Registers) SetBC(value uint16) {
	r.b = byte(value >> 8) // shifting right to get the top bits
	r.c = byte(value)      // cast truncates top keep bottom bits
}

func (r *Registers) DE() uint16 {
	return (uint16(r.d) << 8) | uint16(r.e)
}

func (r *Registers) SetDE(value uint16) {
	r.d = byte(value >> 8) // shifting right to get the top bits
	r.e = byte(value)      // cast truncates top keep bottom bits
}

func (r *Registers) HL() uint16 {
	return (uint16(r.h) << 8) | uint16(r.l)
}

func (r *Registers) SetHL(value uint16) {
	r.h = byte(value >> 8) // shifting right to get the top bits
	r.l = byte(value)      // cast truncates top keep bottom bits
}

func (r *Registers) Z() bool {
	return r.f&FlagZero != 0
}

func (r *Registers) SetZ(value bool) {
	r.setFlag(FlagZero, value)
}

func (r *Registers) N() bool {
	return r.f&FlagSubtract != 0
}

func (r *Registers) SetN(value bool) {
	r.setFlag(FlagSubtract, value)
}

func (r *Registers) H() bool {
	return r.f&FlagHalfCarry != 0
}

func (r *Registers) SetH(value bool) {
	r.setFlag(FlagHalfCarry, value)
}

func (r *Registers) Cy() bool {
	return r.f&FlagCarry != 0
}

func (r *Registers) SetCy(value bool) {
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
