package cpu

import (
	"testing"
)

func TestArithmeticInstructions(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(c *CPU, m *mockMemory)
		run      func(c *CPU) int
		expected func(t *testing.T, c *CPU, m *mockMemory, cycles int)
	}{
		// --- 8-bit ADD ---
		{
			name: "ADD A, B (No Carry)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x10
				c.registers.B = 0x20
			},
			run: func(c *CPU) int {
				return c.add8AReg8(func() byte { return c.registers.B })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x30 {
					t.Errorf("A = %02X, want 30", c.registers.A)
				}
				if c.registers.FlagZ() || c.registers.FlagN() || c.registers.FlagH() || c.registers.FlagCy() {
					t.Errorf("Flags incorrect: Z:%v N:%v H:%v C:%v", c.registers.FlagZ(), c.registers.FlagN(), c.registers.FlagH(), c.registers.FlagCy())
				}
				if cycles != 1 {
					t.Errorf("Cycles = %d, want 1", cycles)
				}
			},
		},
		{
			name: "ADD A, C (Half-Carry)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x0F
				c.registers.C = 0x01
			},
			run: func(c *CPU) int {
				return c.add8AReg8(func() byte { return c.registers.C })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x10 {
					t.Errorf("A = %02X, want 10", c.registers.A)
				}
				if !c.registers.FlagH() {
					t.Error("H flag should be set")
				}
				if c.registers.FlagZ() || c.registers.FlagN() || c.registers.FlagCy() {
					t.Errorf("Flags incorrect: Z:%v N:%v C:%v", c.registers.FlagZ(), c.registers.FlagN(), c.registers.FlagCy())
				}
			},
		},
		{
			name: "ADD A, D (Carry & Zero)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x80
				c.registers.D = 0x80
			},
			run: func(c *CPU) int {
				return c.add8AReg8(func() byte { return c.registers.D })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x00 {
					t.Errorf("A = %02X, want 00", c.registers.A)
				}
				if !c.registers.FlagCy() || !c.registers.FlagZ() {
					t.Errorf("Flags incorrect: C:%v Z:%v", c.registers.FlagCy(), c.registers.FlagZ())
				}
				if c.registers.FlagN() || c.registers.FlagH() {
					t.Errorf("Flags incorrect: N:%v H:%v", c.registers.FlagN(), c.registers.FlagH())
				}
			},
		},
		{
			name: "ADD A, [HL]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x12
				c.registers.SetHL(0xC000)
				m.data[0xC000] = 0x23
			},
			run: func(c *CPU) int {
				return c.add8AHL()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x35 {
					t.Errorf("A = %02X, want 35", c.registers.A)
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},
		{
			name: "ADD A, n8",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x50
			},
			run: func(c *CPU) int {
				return c.add8AImm8(0x10)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x60 {
					t.Errorf("A = %02X, want 60", c.registers.A)
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},

		// --- 8-bit ADC ---
		{
			name: "ADC A, E (With Carry Set)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x01
				c.registers.E = 0x01
				c.registers.SetFlagCy(true)
			},
			run: func(c *CPU) int {
				return c.adc8AReg8(func() byte { return c.registers.E })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x03 { // 1 + 1 + 1
					t.Errorf("A = %02X, want 03", c.registers.A)
				}
				if c.registers.FlagCy() {
					t.Error("Carry flag should be cleared")
				}
			},
		},
		{
			name: "ADC A, H (Half-Carry with Carry)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x0E
				c.registers.H = 0x01
				c.registers.SetFlagCy(true)
			},
			run: func(c *CPU) int {
				return c.adc8AReg8(func() byte { return c.registers.H })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x10 { // 0xE + 0x1 + 0x1 = 0x10
					t.Errorf("A = %02X, want 10", c.registers.A)
				}
				if !c.registers.FlagH() {
					t.Error("H flag should be set")
				}
			},
		},
		{
			name: "ADC A, [HL]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x01
				c.registers.SetHL(0xC000)
				m.data[0xC000] = 0x80
				c.registers.SetFlagCy(true)
			},
			run: func(c *CPU) int {
				return c.adc8AHL()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x82 { // 0x01 + 0x80 + 1 = 0x82
					t.Errorf("A = %02X, want 82", c.registers.A)
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},
		{
			name: "ADC A, n8",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x10
				c.registers.SetFlagCy(true)
			},
			run: func(c *CPU) int {
				return c.adc8AImm8(0x0F)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x20 { // 0x10 + 0x0F + 1 = 0x20
					t.Errorf("A = %02X, want 20", c.registers.A)
				}
				if !c.registers.FlagH() {
					t.Error("H flag should be set")
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},

		// --- 8-bit SBC ---
		{
			name: "SBC A, B (With Carry)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x10
				c.registers.B = 0x01
				c.registers.SetFlagCy(true)
			},
			run: func(c *CPU) int {
				return c.sbc8AReg8(func() byte { return c.registers.B })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x0E { // 0x10 - 0x01 - 1 = 0x0E
					t.Errorf("A = %02X, want 0E", c.registers.A)
				}
			},
		},
		{
			name: "SBC A, n8 (Borrow)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x00
				c.registers.SetFlagCy(true)
			},
			run: func(c *CPU) int {
				return c.sbc8AImm8(0x00)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0xFF { // 0x00 - 0x00 - 1 = -1 (0xFF)
					t.Errorf("A = %02X, want FF", c.registers.A)
				}
				if !c.registers.FlagCy() || !c.registers.FlagH() {
					t.Errorf("Flags incorrect: C:%v H:%v", c.registers.FlagCy(), c.registers.FlagH())
				}
			},
		},
		{
			name: "SBC A, [HL]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x10
				c.registers.SetHL(0xC000)
				m.data[0xC000] = 0x05
				c.registers.SetFlagCy(true)
			},
			run: func(c *CPU) int {
				return c.sbc8AHL()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x0A { // 0x10 - 0x05 - 1 = 0x0A
					t.Errorf("A = %02X, want 0A", c.registers.A)
				}
			},
		},

		// --- 8-bit SUB ---
		{
			name: "SUB A, B",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x30
				c.registers.B = 0x10
			},
			run: func(c *CPU) int {
				return c.sub8AReg8(func() byte { return c.registers.B })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x20 {
					t.Errorf("A = %02X, want 20", c.registers.A)
				}
				if !c.registers.FlagN() {
					t.Error("N flag should be set")
				}
				if c.registers.FlagZ() || c.registers.FlagH() || c.registers.FlagCy() {
					t.Errorf("Flags incorrect: Z:%v H:%v C:%v", c.registers.FlagZ(), c.registers.FlagH(), c.registers.FlagCy())
				}
			},
		},
		{
			name: "SUB A, C (Borrow/Carry)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x01
				c.registers.C = 0x02
			},
			run: func(c *CPU) int {
				return c.sub8AReg8(func() byte { return c.registers.C })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0xFF {
					t.Errorf("A = %02X, want FF", c.registers.A)
				}
				if !c.registers.FlagCy() || !c.registers.FlagH() || !c.registers.FlagN() {
					t.Errorf("Flags incorrect: C:%v H:%v N:%v", c.registers.FlagCy(), c.registers.FlagH(), c.registers.FlagN())
				}
			},
		},
		{
			name: "SUB A, [HL]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x30
				c.registers.SetHL(0xD000)
				m.data[0xD000] = 0x10
			},
			run: func(c *CPU) int {
				return c.sub8AHL()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x20 {
					t.Errorf("A = %02X, want 20", c.registers.A)
				}
				if !c.registers.FlagN() {
					t.Error("N flag should be set")
				}
			},
		},
		{
			name: "SUB A, n8",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x10
			},
			run: func(c *CPU) int {
				return c.sub8AImm8(0x01)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x0F {
					t.Errorf("A = %02X, want 0F", c.registers.A)
				}
				if !c.registers.FlagH() {
					t.Error("H flag should be set")
				}
			},
		},

		// --- 8-bit CP ---
		{
			name: "CP A, L (Equal)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x42
				c.registers.L = 0x42
			},
			run: func(c *CPU) int {
				return c.cp8AReg8(func() byte { return c.registers.L })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x42 {
					t.Error("CP should not modify A")
				}
				if !c.registers.FlagZ() || !c.registers.FlagN() {
					t.Errorf("Flags incorrect: Z:%v N:%v", c.registers.FlagZ(), c.registers.FlagN())
				}
			},
		},
		{
			name: "CP A, [HL]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x10
				c.registers.SetHL(0xD000)
				m.data[0xD000] = 0x11
			},
			run: func(c *CPU) int {
				return c.cp8AHL()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x10 {
					t.Error("CP should not modify A")
				}
				if !c.registers.FlagCy() || !c.registers.FlagN() {
					t.Errorf("Flags incorrect: C:%v N:%v", c.registers.FlagCy(), c.registers.FlagN())
				}
			},
		},
		{
			name: "CP A, n8",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x20
			},
			run: func(c *CPU) int {
				return c.cp8AImm8(0x20)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x20 {
					t.Error("CP should not modify A")
				}
				if !c.registers.FlagZ() {
					t.Error("Z flag should be set")
				}
			},
		},

		// --- 8-bit INC/DEC ---
		{
			name: "INC B",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.B = 0x0F
			},
			run: func(c *CPU) int {
				return c.inc8Reg8(func(v byte) { c.registers.B = v }, func() byte { return c.registers.B })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.B != 0x10 {
					t.Errorf("B = %02X, want 10", c.registers.B)
				}
				if !c.registers.FlagH() || c.registers.FlagN() {
					t.Errorf("Flags incorrect: H:%v N:%v", c.registers.FlagH(), c.registers.FlagN())
				}
			},
		},
		{
			name: "DEC C",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.C = 0x10
			},
			run: func(c *CPU) int {
				return c.dec8Reg8(func(v byte) { c.registers.C = v }, func() byte { return c.registers.C })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.C != 0x0F {
					t.Errorf("C = %02X, want 0F", c.registers.C)
				}
				if !c.registers.FlagH() || !c.registers.FlagN() {
					t.Errorf("Flags incorrect: H:%v N:%v", c.registers.FlagH(), c.registers.FlagN())
				}
			},
		},
		{
			name: "DEC [HL] (Zero Flag)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetHL(0xD000)
				m.data[0xD000] = 0x01
			},
			run: func(c *CPU) int {
				return c.dec8HL()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if m.data[0xD000] != 0x00 {
					t.Errorf("Mem[D000] = %02X, want 00", m.data[0xD000])
				}
				if !c.registers.FlagZ() || !c.registers.FlagN() {
					t.Errorf("Flags incorrect: Z:%v N:%v", c.registers.FlagZ(), c.registers.FlagN())
				}
				if cycles != 3 {
					t.Errorf("Cycles = %d, want 3", cycles)
				}
			},
		},
		{
			name: "INC [HL] (Half-Carry)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetHL(0xD001)
				m.data[0xD001] = 0x0F
			},
			run: func(c *CPU) int {
				return c.inc8HL()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if m.data[0xD001] != 0x10 {
					t.Errorf("Mem[D001] = %02X, want 10", m.data[0xD001])
				}
				if !c.registers.FlagH() {
					t.Error("H flag should be set")
				}
			},
		},

		// --- 16-bit ADD/INC/DEC ---
		{
			name: "ADD HL, BC",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetHL(0x1000)
				c.registers.SetBC(0x0100)
			},
			run: func(c *CPU) int {
				return c.add16HLReg16(c.registers.BC)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.HL() != 0x1100 {
					t.Errorf("HL = %04X, want 1100", c.registers.HL())
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},
		{
			name: "ADD HL, DE (H-Carry bit 11)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetHL(0x0FFF)
				c.registers.SetDE(0x0001)
			},
			run: func(c *CPU) int {
				return c.add16HLReg16(c.registers.DE)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.HL() != 0x1000 {
					t.Errorf("HL = %04X, want 1000", c.registers.HL())
				}
				if !c.registers.FlagH() {
					t.Error("H flag should be set (bit 11 overflow)")
				}
			},
		},
		{
			name: "ADD HL, SP",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetHL(0x1000)
				c.registers.SP = 0x0500
			},
			run: func(c *CPU) int {
				return c.add16HLSP()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.HL() != 0x1500 {
					t.Errorf("HL = %04X, want 1500", c.registers.HL())
				}
			},
		},
		{
			name: "DEC SP",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SP = 0x0000
			},
			run: func(c *CPU) int {
				return c.dec16Reg16(func(v uint16) { c.registers.SP = v }, func() uint16 { return c.registers.SP })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.SP != 0xFFFF {
					t.Errorf("SP = %04X, want FFFF", c.registers.SP)
				}
			},
		},
		{
			name: "INC BC",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetBC(0xFFFF)
				c.registers.SetFlagZ(true) // Should NOT be modified
			},
			run: func(c *CPU) int {
				return c.inc16Reg16(c.registers.SetBC, c.registers.BC)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.BC() != 0x0000 {
					t.Errorf("BC = %04X, want 0000", c.registers.BC())
				}
				if !c.registers.FlagZ() {
					t.Error("Z flag should remain true")
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, mem := createTestCPU()
			tt.setup(cpu, mem)
			cycles := tt.run(cpu)
			tt.expected(t, cpu, mem, cycles)
		})
	}
}
