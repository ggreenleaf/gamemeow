package cpu

// Naming Mechanism & Nomenclature:
// -----------------------------
// op8   = 8-bit arithmetic/logic operation (ADD, ADC, SUB, SBC, AND, XOR, OR, CP)
// op16  = 16-bit arithmetic operation (ADD, INC, DEC)
// A     = Implicit target for most 8-bit arithmetic
// Reg8  / Reg16 = Register source/destination
// Imm8  / Imm16 = Immediate value source
// HL    = Memory address pointed to by HL source/destination
//
// Naming Pattern Examples:
// add8AReg8(srcGet func() byte) int
// add8AImm8(val byte) int
// add8AHL() int
// inc8Reg8(dst func(byte), srcGet func() byte) int
// inc8HL() int
// add16HLReg16(srcGet func() uint16) int
// -----------------------------

// adc8AReg8 handles ADC A, r8
// Add the value in r8 plus the carry flag to A.
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 0, H (Set if bit 3 overflow), C (Set if bit 7 overflow)
func (c *CPU) adcAReg8(srcGet func() byte) int {
	a := c.registers.A
	value := srcGet()
	cy := byte(c.registers.FlagCyBit())
	res16 := uint16(a) + uint16(value) + uint16(cy)
	hCarry := (a&0x0F)+(value&0x0F)+cy > 0x0F
	cyCarry := res16 > 0xFF
	c.registers.A = byte(res16)
	c.registers.SetFlagZ(c.registers.A == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.registers.SetFlagCy(cyCarry)
	return 1
}

// adc8AHL handles ADC A, [HL]
// Add the byte pointed to by HL plus the carry flag To A
// cycles 2 | bytes 1 | flags Z (Set if result 0), N 0, H (Set if bit 3 overflow), C (Set if bit 7 overflow)
func (c *CPU) adcAHLPtr() int {
	value := c.bus.Read(c.registers.HL())
	a := c.registers.A
	cy := byte(c.registers.FlagCyBit())
	res16 := uint16(a) + uint16(value) + uint16(cy)
	hCarry := (a&0x0F)+(value&0x0F)+cy > 0x0F
	cyCarry := res16 > 0xFF
	c.registers.A = byte(res16)
	c.registers.SetFlagZ(c.registers.A == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.registers.SetFlagCy(cyCarry)
	return 2
}

// adc8AImm8 handles ADC A, n8
// Add the value n8 plus the carry flag to A
// cycles 2 | bytes 2 | flags Z (Set if result 0), N 0, H (Set if bit 3 overflow), C (Set if bit 7 overflow)
func (c *CPU) adcAImm8(srcVal byte) int {
	a := c.registers.A
	cy := byte(c.registers.FlagCyBit())
	res16 := uint16(a) + uint16(srcVal) + uint16(cy)
	hCarry := (a&0x0F)+(srcVal&0x0F)+cy > 0x0F
	cyCarry := res16 > 0xFF
	c.registers.A = byte(res16)
	c.registers.SetFlagZ(c.registers.A == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.registers.SetFlagCy(cyCarry)
	return 2
}

// add8AReg8 handles ADD A, r8
// Add the value in r8 to A
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 0, H (Set if bit 3 overflow), C (Set if bit 7 overflow)
func (c *CPU) addAReg8(srcGet func() byte) int {
	a := c.registers.A
	value := srcGet()
	res16 := uint16(a) + uint16(value)
	hCarry := (a&0x0F)+(value&0x0F) > 0x0F
	cyCarry := res16 > 0xFF
	c.registers.A = byte(res16)
	c.registers.SetFlagZ(c.registers.A == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.registers.SetFlagCy(cyCarry)

	return 1
}

// add8AHL handles ADD A, [HL]
// Add the byte pointed to by HL to A.
// cycles 2 | bytes 1 | flags Z (Set if result 0), N 0, H (Set if bit 3 overflow), C (Set if bit 7 overflow)
func (c *CPU) addAHLPtr() int {
	a := c.registers.A
	value := c.bus.Read(c.registers.HL())
	res16 := uint16(a) + uint16(value)
	hCarry := (a&0x0F)+(value&0x0F) > 0x0F
	cyCarry := res16 > 0xFF
	c.registers.A = byte(res16)
	c.registers.SetFlagZ(c.registers.A == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.registers.SetFlagCy(cyCarry)
	return 2
}

// add8AImm8 handles ADD A, n8
// Add the value n8 to A
// cycles 2 | bytes 2 | flags Z (Set if result 0), N 0, H (Set if bit 3 overflow), C (Set if bit 7 overflow)
func (c *CPU) addAImm8(srcVal byte) int {
	a := c.registers.A
	res16 := uint16(a) + uint16(srcVal)
	hCarry := (a&0x0F)+(srcVal&0x0F) > 0x0F
	cyCarry := res16 > 0xFF
	c.registers.A = byte(res16)
	c.registers.SetFlagZ(c.registers.A == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.registers.SetFlagCy(cyCarry)
	return 2
}

// add16HLReg16 handles ADD HL, r16
// Add the value in r16 to HL
// cycles 2 | bytes 1 | flags N 0, H (set if bit 11 overflow), C (set if bit 15 overflow)
func (c *CPU) addHLReg16(srcGet func() uint16) int {
	value := srcGet()
	hl := c.registers.HL()
	res32 := uint32(value) + uint32(hl)
	hCarry := (hl&0x0FFF)+(value&0x0FFF) > 0x0FFF
	cyCarry := res32 > 0xFFFF
	c.registers.SetHL(uint16(res32))
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.registers.SetFlagCy(cyCarry)
	return 2
}

// add16HLSP handles ADD HL, SP
// add the value in SP to HL
// cycles 2 | bytes 1 | flags N 0, H (set if bit 11 overflow), C (set if bit 15 overflow)
func (c *CPU) addHLSP() int {
	hl := c.registers.HL()
	sp := c.registers.SP
	res32 := uint32(sp) + uint32(hl)
	hCarry := (hl&0x0FFF)+(sp&0x0FFF) > 0x0FFF
	cyCarry := res32 > 0xFFFF
	c.registers.SetHL(uint16(res32))
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.registers.SetFlagCy(cyCarry)

	return 2
}

// cp8AHL handles CP A, [HL]
// Compare A with the byte pointed to by HL (subtract [HL] from A but don't store result)
// cycles 2 | bytes 1 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4), C (Set if borrow (A < src))
func (c *CPU) cpAHLPtr() int {
	a := c.registers.A
	value := c.bus.Read(c.registers.HL())
	hBorrow := (a & 0x0F) < (value & 0x0F)
	cyBorrow := a < value
	res := a - value
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.registers.SetFlagCy(cyBorrow)

	return 2
}

// cp8AImm8 handles CP A, n8
// Compare A with the value n8 (subtract n8 from A but don't store result)
// cycles 2 | bytes 2 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4), C (Set if borrow (A < src))
func (c *CPU) cpAImm8(srcVal byte) int {
	a := c.registers.A
	hBorrow := (a & 0x0F) < (srcVal & 0x0F)
	cyBorrow := a < srcVal
	res := a - srcVal
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.registers.SetFlagCy(cyBorrow)

	return 2
}

// cp8AReg8 handles CP A, r8
// Compare A with the value in r8 (subtract r8 from A but don't store result)
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4), C (Set if borrow (A < src))
func (c *CPU) cpAReg8(srcGet func() byte) int {
	a := c.registers.A
	value := srcGet()
	hBorrow := (a & 0x0F) < (value & 0x0F)
	cyBorrow := a < value
	res := a - value
	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.registers.SetFlagCy(cyBorrow)
	return 1
}

// dec16Reg16 handles DEC r16
// Decrement the value in r16
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) decReg16(dst func(uint16), srcGet func() uint16) int {
	value := srcGet()
	value--
	dst(value)
	return 2
}

// dec8HL handles DEC [HL]
// Decrement the byte pointed to by HL
// cycles 3 | bytes 1 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4)
func (c *CPU) decHLPtr() int {
	addr := c.registers.HL()
	value := c.bus.Read(addr)
	res := value - 1
	hBorrow := (value & 0x0F) == 0

	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.bus.Write(addr, res)
	return 3
}

// dec8Reg8 handles DEC r8
// Decrement the value in r8
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4)
func (c *CPU) decReg8(dst func(byte), srcGet func() byte) int {
	value := srcGet()

	res := value - 1
	hBorrow := (value & 0x0F) == 0

	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	dst(res)
	return 1
}

// inc16Reg16 handles INC r16
// Increment the value in r16
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) incReg16(dst func(uint16), srcGet func() uint16) int {
	value := srcGet()
	value++
	dst(value)
	return 2
}

// inc8HL handles INC [HL]
// Increment the byte pointed to by HL
// cycles 3 | bytes 1 | flags Z (Set if result 0), N 0, H (Set if bit 3 overflow)
func (c *CPU) incHLPtr() int {
	addr := c.registers.HL()
	value := c.bus.Read(addr)

	hCarry := (value & 0x0F) == 0x0F
	res := value + 1

	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.bus.Write(addr, res)

	return 3
}

// inc8Reg8 handles INC r8
// Increment the value in r8
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 0, H (Set if bit 3 overflow)
func (c *CPU) incReg8(dst func(byte), srcGet func() byte) int {

	value := srcGet()

	hCarry := (value & 0x0F) == 0x0F
	res := value + 1

	c.registers.SetFlagZ(res == 0)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	dst(res)

	return 1
}

// sbc8AHL handles SBC A, [HL]
// Subtract the byte pointed to by HL plus the carry flag from A
// cycles 2 | bytes 1 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4), C (Set if borrow (A < src + cy))
func (c *CPU) sbcAHLPtr() int {
	value := int(c.bus.Read(c.registers.HL()))
	cy := int(c.registers.FlagCyBit())
	a := int(c.registers.A)

	res := a - value - cy

	hBorrow := (a&0x0F)-(value&0x0F)-cy < 0
	cyBorrow := res < 0

	finalRes := byte(res)
	c.registers.A = finalRes
	c.registers.SetFlagZ(finalRes == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.registers.SetFlagCy(cyBorrow)

	return 2
}

// sbc8AImm8 handles SBC A, n8
// Subtract the value n8 plus the carry flag from A
// cycles 2 | bytes 2 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4), C (Set if borrow (A < src + cy))
func (c *CPU) sbcAImm8(srcVal byte) int {
	value := int(srcVal)
	cy := int(c.registers.FlagCyBit())
	a := int(c.registers.A)

	res := a - value - cy

	hBorrow := (a&0x0F)-(value&0x0F)-cy < 0
	cyBorrow := res < 0

	finalRes := byte(res)
	c.registers.A = finalRes
	c.registers.SetFlagZ(finalRes == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.registers.SetFlagCy(cyBorrow)

	return 2
}

// sbc8AReg8 handles SBC A, r8
// Subtract the value in r8 plus the carry flag from A
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4), C (Set if borrow (A < src + cy))
func (c *CPU) sbcAReg8(srcGet func() byte) int {
	value := int(srcGet())
	cy := int(c.registers.FlagCyBit())
	a := int(c.registers.A)

	res := a - value - cy

	hBorrow := (a&0x0F)-(value&0x0F)-cy < 0
	cyBorrow := res < 0

	finalRes := byte(res)
	c.registers.A = finalRes
	c.registers.SetFlagZ(finalRes == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.registers.SetFlagCy(cyBorrow)

	return 1
}

// sub8AHL handles SUB A, [HL]
// Subtract the byte pointed to by HL from A
// cycles 2 | bytes 1 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4), C (Set if borrow (A < src))
func (c *CPU) subAHLPtr() int {

	value := int(c.bus.Read(c.registers.HL()))
	a := int(c.registers.A)

	res := a - value

	hBorrow := (a&0x0F)-(value&0x0F) < 0
	cyBorrow := res < 0

	finalRes := byte(res)
	c.registers.A = finalRes
	c.registers.SetFlagZ(finalRes == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.registers.SetFlagCy(cyBorrow)
	return 2
}

// sub8AImm8 handles SUB A, n8
// Subtract the value n8 from A
// cycles 2 | bytes 2 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4), C (Set if borrow (A < src))
func (c *CPU) subAImm8(srcVal byte) int {

	value := int(srcVal)
	a := int(c.registers.A)

	res := a - value

	hBorrow := (a&0x0F)-(value&0x0F) < 0
	cyBorrow := res < 0

	finalRes := byte(res)
	c.registers.A = finalRes
	c.registers.SetFlagZ(finalRes == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.registers.SetFlagCy(cyBorrow)

	return 2
}

// sub8AReg8 handles SUB A, r8
// Subtract the value in r8 from A
// cycles 1 | bytes 1 | flags Z (Set if result 0), N 1, H (Set if borrow from bit 4), C (Set if borrow (A < src))
func (c *CPU) subAReg8(srcGet func() byte) int {
	value := int(srcGet())
	a := int(c.registers.A)

	res := a - value

	hBorrow := (a&0x0F)-(value&0x0F) < 0
	cyBorrow := res < 0

	finalRes := byte(res)
	c.registers.A = finalRes
	c.registers.SetFlagZ(finalRes == 0)
	c.registers.SetFlagN(true)
	c.registers.SetFlagH(hBorrow)
	c.registers.SetFlagCy(cyBorrow)

	return 1
}
