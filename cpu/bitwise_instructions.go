package cpu

// Naming Mechanism & Nomenclature:
// -----------------------------
// op    = bitwise operation (AND, XOR, OR)
// A     = Implicit target for bitwise operations
// Reg8  = Register source
// Imm8  = Immediate value source
// HLPtr = Memory address pointed to by HL source
//
// Naming Pattern Examples:
// andAReg8(srcGet func() byte) int
// andAImm8(val byte) int
// andAHLPtr() int
// -----------------------------

// andAReg8 handles AND A, r8
// Set A to the bitwise AND between the value in r8 and A
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 0, H 1, C 0
func (c *CPU) andAReg8(srcGet func() byte) int {
	a := c.registers.A
	value := srcGet()
	res := a & value
	c.registers.A = res
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(true)
	c.registers.SetFlagCy(false)
	return 1
}

// andAHLPtr handles AND A, [HL]
// Set A to the bitwise AND between the byte pointed to by HL and A
// cycles 2 | bytes 1 | flags Z (Set if result 0), N 0, H 1, C 0
func (c *CPU) andAHLPtr() int {
	a := c.registers.A

	addr := c.registers.HL()
	value := c.bus.Read(addr)
	res := a & value

	c.registers.A = res
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(true)
	c.registers.SetFlagCy(false)
	return 2
}

// andAImm8 handles AND A, n8
// Set A to the bitwise AND between the value n8 and A
// cycles 2 | bytes 2 | flags Z (Set if result 0), N 0, H 1, C 0
func (c *CPU) andAImm8(srcVal byte) int {
	a := c.registers.A
	res := a & srcVal

	c.registers.A = res
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(true)
	c.registers.SetFlagCy(false)
	return 2
}

// cpl handles cpl
// Complement accumulator (A = ~A) also called bitwise NOT
// cycles 1 | bytes 1 | flags N 1, H 1
func (c *CPU) cpl() int {
	c.registers.A = ^c.registers.A
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(true)
	return 1
}

// orAReg8 handles OR A, r8
// Set A to the bitwise OR between the value in r8 and A
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 0, H 0, C 0
func (c *CPU) orAReg8(srcGet func() byte) int {
	a := c.registers.A
	value := srcGet()
	res := a | value
	c.registers.A = res
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(false)
	c.registers.SetFlagCy(false)
	return 1
}

// orAHLPtr handles OR A, [HL]
// Set A to the bitwise OR between the byte pointed to by HL and A
// cycles 2 | bytes 1 | flags Z (Set if result 0), N 0, H 0, C 0
func (c *CPU) orAHLPtr() int {
	a := c.registers.A

	addr := c.registers.HL()
	value := c.bus.Read(addr)
	res := a | value

	c.registers.A = res
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(false)
	c.registers.SetFlagCy(false)
	return 2
}

// orAImm8 handles OR A, n8
// Set A to the bitwise OR between the value n8 and A
// cycles 2 | bytes 2 | flags Z (Set if result 0), N 0, H 0, C 0
func (c *CPU) orAImm8(srcVal byte) int {
	a := c.registers.A
	res := a | srcVal

	c.registers.A = res
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(false)
	c.registers.SetFlagCy(false)
	return 2
}

// xorAReg8 handles XOR A, r8
// Set A to the bitwise XOR between the value in r8 and A
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 0, H 0, C 0
func (c *CPU) xorAReg8(srcGet func() byte) int {
	a := c.registers.A
	value := srcGet()
	res := a ^ value
	c.registers.A = res
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(false)
	c.registers.SetFlagCy(false)
	return 1
}

// xorAHLPtr handles XOR A, [HL]
// Set A to the bitwise XOR between the byte pointed to by HL and A
// cycles 2 | bytes 1 | flags Z (Set if result 0), N 0, H 0, C 0
func (c *CPU) xorAHLPtr() int {
	a := c.registers.A

	addr := c.registers.HL()
	value := c.bus.Read(addr)
	res := a ^ value

	c.registers.A = res
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(false)
	c.registers.SetFlagCy(false)
	return 2
}

// xorAImm8 handles XOR A, n8
// Set A to the bitwise XOR between the value n8 and A
// cycles 2 | bytes 2 | flags Z (Set if result 0), N 0, H 0, C 0
func (c *CPU) xorAImm8(srcVal byte) int {
	a := c.registers.A
	res := a ^ srcVal

	c.registers.A = res
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(false)
	c.registers.SetFlagCy(false)
	return 2
}
