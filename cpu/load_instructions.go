package cpu

// Naming Mechanism & Nomenclature:
// -----------------------------
// load  = CPU <- Source (Reading into a register)
// store = Memory <- Source (Writing to a memory address)
//
// Parameters:
// dst     = Target setter function: func(byte) or func(uint16)
// srcVal  = A direct 8-bit or 16-bit value (r8, n8, n16)
// srcGet  = Function to get a source value: func() byte
// addrPtr = Function to get a 16-bit address (used in [r16] or [HL]): func() uint16
// srcAddr = A direct 16-bit address
// -----------------------------

// loadReg8Reg8 handles LD r8, r8
// Copy (aka Load) the value in register on the right into the register on the left.
// Storing a register into itself is a no-op; however, some Game Boy emulators
// interpret LD B,B as a breakpoint, or LD D,D as a debug message (such as BGB).
// cycles 1 | bytes 1 | flags None affected.
func (c *CPU) loadReg8Reg8(dst *byte, srcVal byte) int {
	*dst = srcVal
	return 1
}

// loadReg8Imm8 handles LD r8, n8
// Copy the value n8 into register r8.
// cycles 2 | bytes 2 | flags None affected.
func (c *CPU) loadReg8Imm8(dst *byte, srcVal byte) int {
	*dst = srcVal
	return 2
}

// loadReg16Imm16 handles LD r16, n16
// Copy the value n16 into register r16.
// cycles 3 | bytes 3 | flags None affected.
func (c *CPU) loadReg16Imm16(dst func(uint16), srcVal uint16) int {
	dst(srcVal)
	return 3
}

// storeHLReg8 handles LD [HL], r8
// Copy the value in register r8 into the byte pointed to by HL.
// cycles 2 | bytes 1 | flags None affected.
func (c *CPU) storeHLPtrReg8(srcVal byte) int {
	c.bus.Write(c.registers.HL(), srcVal)
	return 2
}

// storeHLImm8 handles LD [HL], n8
// Copy the value n8 into the byte pointed to by HL.
// cycles 3 | bytes 2 | flags None affected.
func (c *CPU) storeHLPtrImm8(srcVal byte) int {
	c.bus.Write(c.registers.HL(), srcVal)
	return 3
}

// loadReg8HL handles LD r8, [HL]
// Copy the value pointed to by HL into register r8.
// cycles 2 | bytes 1 | flags None affected.
func (c *CPU) loadReg8HLPtr(dst *byte) int {
	*dst = c.bus.Read(c.registers.HL())
	return 2
}

// storeReg16A handles LD [r16], A
// Copy the value in register A into the byte pointed to by r16
// cycles 2 | bytes 1 | flags None affected
func (c *CPU) storeReg16A(addrPtr func() uint16) int {
	c.bus.Write(addrPtr(), c.registers.A)
	return 2
}

// storeImm16A handles LD [n16], A
// Copy the value in register A into the byte at address n16
// cycles 4 | bytes 3 | flags None affected
func (c *CPU) storeImm16A(targetAddr uint16) int {
	c.bus.Write(targetAddr, c.registers.A)
	return 4
}

// storeHighImm16A handles LDH [n16], A
// Copy the value in register A into the byte at address n16, provided the address is between 0xFF00 and 0xFFFF
// cycles 3 | bytes 2 | flags none affected
func (c *CPU) storeHighImm8A(offset byte) int {
	value := c.registers.A
	addr := uint16(offset) + 0xFF00
	c.bus.Write(addr, value)
	return 3
}

// storeHighCA handles LDH [C], A
// Copy the value in register A into the byte at address 0xFF00+C
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) storeHighCA() int {
	value := c.registers.A
	addr := uint16(c.registers.C) + 0xFF00
	c.bus.Write(addr, value)
	return 2
}

// loadAReg16 handles LD A, [r16]
// Copy the byte pointed to by r16 into register A
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) loadAReg16(addrPtr func() uint16) int {
	c.registers.A = c.bus.Read(addrPtr())
	return 2
}

// loadAImm16 handles LD A, [n16]
// Copy the byte at address n16 into register A
// cycles 4 | bytes 3 | flags none affected
func (c *CPU) loadAImm16(srcAddr uint16) int {
	value := c.bus.Read(srcAddr)
	c.registers.A = value
	return 4
}

// loadHighAImm8 handles LDH A, [n8]
// Copy the byte at address 0xFF00+n8 into register A
// cycles 3 | bytes 2 | flags none affected
func (c *CPU) loadHighAImm8(offset byte) int {
	c.registers.A = c.bus.Read(uint16(offset) + 0xFF00)
	return 3
}

// loadHighAC handles LDH A, [C]
// Copy the byte at address 0xFF00+C into register A
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) loadHighAC() int {
	value := c.bus.Read(uint16(c.registers.C) + 0xFF00)
	c.registers.A = value
	return 2
}

// storeHLIncA handles LD [HLI], A
// Copy the value in register A into the byte pointed by HL and increment HL afterwards
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) storeHLPtrIncA() int {
	value := c.registers.A
	addr := c.registers.HL()
	c.bus.Write(addr, value)
	c.registers.SetHL(c.registers.HL() + 1)
	return 2
}

// storeHLDecA handles LD [HLD], A
// Copy the value in register A into the byte pointed by HL and decrement HL afterwards
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) storeHLPtrDecA() int {
	value := c.registers.A
	addr := c.registers.HL()
	c.bus.Write(addr, value)
	c.registers.SetHL(c.registers.HL() - 1)
	return 2
}

// loadAHLDec handles LD A, [HLD]
// Copy the byte pointed to by HL into register A, and decrement HL afterwards
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) loadAHLPtrDec() int {
	value := c.bus.Read(c.registers.HL())
	c.registers.A = value
	c.registers.SetHL(c.registers.HL() - 1)
	return 2
}

// loadAHLInc handles LD A, [HLI]
// Copy the byte pointed to by HL into register A, and increment HL afterwards
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) loadAHLPtrInc() int {
	value := c.bus.Read(c.registers.HL())
	c.registers.A = value
	c.registers.SetHL(c.registers.HL() + 1)
	return 2
}

// loadSPImm16 handles LD SP, n16
// Copy the value n16 into register SP
// cycles 3 | bytes 3 | flags none affected
func (c *CPU) loadSPImm16(srcVal uint16) int {
	c.registers.SP = srcVal
	return 3
}

// storeImm16SP handles LD [n16], SP
// Copy SP & 0xFF at address n16 and SP >> 8 at address n16 + 1
// cyles 5 | bytes 3 | flags none affected
func (c *CPU) storeImm16SP(srcAddr uint16) int {
	c.bus.Write(srcAddr, byte(c.registers.SP)) // this will auto truncate the lower 8 bits
	c.bus.Write(srcAddr+1, byte(c.registers.SP>>8))
	return 5
}

// loadHLSPSigned8 handles LD HL, SP+e8
// Add the signed value e8 to SP and copy the result in HL
// cycles 3 | bytes 2 | flags (Z 0) (N 0) (H Set if overflow from bit 3) (C Set if overflow from bit 7)
func (c *CPU) loadHLSPSigned8(signed8 int8) int {
	spLow := uint16(c.registers.SP & 0xFF)
	unsignedOffset := uint16(uint8(signed8))
	c.registers.SetHL(c.registers.SP + uint16(int16(signed8)))

	hCarry := (spLow&0x0F)+(unsignedOffset&0x0F) > 0x0F
	cyCarry := spLow+unsignedOffset > 0xFF

	c.registers.SetFlagZ(false)
	c.registers.SetFlagN(false)
	c.registers.SetFlagH(hCarry)
	c.registers.SetFlagCy(cyCarry)

	return 3
}

// loadSPHL handles LD SP, HL
// Copy register HL into register SP
// cycles 2 | bytes 1 | flags none affected
func (c *CPU) loadSPHL() int {
	c.registers.SP = c.registers.HL()
	return 2
}
