package cpu

// Naming Mechanism & Nomenclature:
// -----------------------------
// op    		= bitwise operation (AND, XOR, OR, etc...)
// A     		= Implicit target for bitwise operations
// Reg8  		= Register source
// Imm8  		= Immediate value source
// HLPtr 		= Memory address pointed to by HL source
// bitIndex = bit position 0 beeing right most bit 7 being the leftmost bit
//
// Naming Pattern Examples:
// andAReg8(srcGet func() byte) int
// andAImm8(val byte) int
// andAHLPtr() int
// -----------------------------

// andAReg8 handles AND A, r8
// Set A to the bitwise AND between the value in r8 and A
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 0, H 1, C 0
func (c *CPU) andAReg8(srcVal byte) int {
	a := c.registers.A
	value := srcVal
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
func (c *CPU) orAReg8(srcVal byte) int {
	a := c.registers.A
	value := srcVal
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
func (c *CPU) xorAReg8(srcVal byte) int {
	a := c.registers.A
	value := srcVal
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

// bitIndexImm8 handles BIT u3, r8
// Test bit u3 in register r8, set the zero flag if bit not set
// cycles 2 | bytes 2 | flags Z (set if selected bit is 0), n 0, h 1
func (c *CPU) bitIndexImm8(bitIndex byte, srcVal byte) int {
	value := srcVal

	isZero := (value & (1 << bitIndex)) == 0

	c.registers.SetFlagZ(isZero)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(true)

	return 2
}

// bitIndexHLPtr handles BIT u3, [HL]
// Test bit u3 in the byte pointed by HL, set the zero flag if bit not set
// cycles 3 | bytes 2 | flags Z (set if selected bit is 0), n 0, h 1
func (c *CPU) bitIndexHlPtr(bitIndex byte) int {
	addr := c.registers.HL()
	value := c.bus.Read(addr)

	isZero := (value & (1 << bitIndex)) == 0

	c.registers.SetFlagZ(isZero)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(true)

	return 3
}

// resIndexReg8 handles RES u3, r8
// Set bit ue in register r8 to 0. Bit 0 is the right most one, bit 7 the left most one
// cycles 2 | bytes 2 | flags none affected
func (c *CPU) resIndexReg8(bitIndex byte, reg *byte) int {
	value := *reg
	res := value &^ (1 << bitIndex)
	*reg = res
	return 2
}

// resIndexHLPtr handles RES u3, [HL]
// Set the bit u3 in teh byte pointed by HL to 0. Bit 0 is the rightmost one, bit 7 is the left most one
// cycles 4 | bytes 2 | flags none affected
func (c *CPU) resIndexHLPtr(bitIndex byte) int {
	addr := c.registers.HL()
	value := c.bus.Read(addr)
	res := value &^ (1 << bitIndex)
	c.bus.Write(addr, res)
	return 4
}

// setIndexReg8 handles SET u3, r8
// Set the bit u3 in register r8 to 1. Bit 0 is the right most one, bit 7 the leftmost one
// cycles 2 | bytes 2 | flags none affected
func (c *CPU) setIndexReg8(bitIndex byte, reg *byte) int {
	value := *reg
	res := value | (1 << bitIndex)
	*reg = res
	return 2
}

// setIndexReg8 handles SET u3, r8
// set the bit u3 in the byte pointed by HL to 1. BIt 0 is teh rightmost one, b it 7 the leftmost one
// cycles 4 | bytes 2 | flags none affected
func (c *CPU) setIndexHLPtr(bitIndex byte) int {
	addr := c.registers.HL()
	value := c.bus.Read(addr)
	res := value | (1 << bitIndex)
	c.bus.Write(addr, res)
	return 4
}
