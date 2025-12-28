package cpu

import (
	"testing"
)

func TestArithmeticInstructions(t *testing.T) {
	type testCase struct {
		name     string
		setup    func(c *CPU, m *mockMemory)
		run      func(c *CPU) int
		expected func(t *testing.T, c *CPU, m *mockMemory, cycles int)
	}

	tests := []struct {
		groupName string
		cases     []testCase
	}{
		{
			groupName: "ADD (8-bit)",
			cases: []testCase{
				{
					name: "ADD A, B (No Carry)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0x10
						c.registers.B = 0x20
					},
					run: func(c *CPU) int {
						return c.addAReg8(c.registers.B)
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
						return c.addAReg8(c.registers.C)
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
						return c.addAReg8(c.registers.D)
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
						return c.addAHLPtr()
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
						return c.addAImm8(0x10)
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
				{
					name: "ADD A, A (Double)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0x40
					},
					run: func(c *CPU) int {
						return c.addAReg8(c.registers.A)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x80 {
							t.Errorf("A = %02X, want 80", c.registers.A)
						}
						if c.registers.FlagH() {
							t.Error("H flag should NOT be set for 0x40+0x40")
						}
					},
				},
				{
					name: "ADD A, A (Double with Carry)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0x80
					},
					run: func(c *CPU) int {
						return c.addAReg8(c.registers.A)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x00 {
							t.Errorf("A = %02X, want 00", c.registers.A)
						}
						if !c.registers.FlagCy() {
							t.Error("Carry flag should be set")
						}
						if !c.registers.FlagZ() {
							t.Error("Zero flag should be set")
						}
					},
				},
			},
		},
		{
			groupName: "ADC (8-bit)",
			cases: []testCase{
				{
					name: "ADC A, E (With Carry Set)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0x01
						c.registers.E = 0x01
						c.registers.SetFlagCy(true)
					},
					run: func(c *CPU) int {
						return c.adcAReg8(c.registers.E)
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
						return c.adcAReg8(c.registers.H)
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
						return c.adcAHLPtr()
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
						return c.adcAImm8(0x0F)
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
			},
		},
		{
			groupName: "SBC (8-bit)",
			cases: []testCase{
				{
					name: "SBC A, B (With Carry)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0x10
						c.registers.B = 0x01
						c.registers.SetFlagCy(true)
					},
					run: func(c *CPU) int {
						return c.sbcAReg8(c.registers.B)
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
						return c.sbcAImm8(0x00)
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
						return c.sbcAHLPtr()
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x0A { // 0x10 - 0x05 - 1 = 0x0A
							t.Errorf("A = %02X, want 0A", c.registers.A)
						}
					},
				},
			},
		},
		{
			groupName: "SUB (8-bit)",
			cases: []testCase{
				{
					name: "SUB A, B",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0x30
						c.registers.B = 0x10
					},
					run: func(c *CPU) int {
						return c.subAReg8(c.registers.B)
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
						return c.subAReg8(c.registers.C)
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
						return c.subAHLPtr()
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
						return c.subAImm8(0x01)
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
				{
					name: "SUB A, A (Zero)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0xAA
					},
					run: func(c *CPU) int {
						return c.subAReg8(c.registers.A)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x00 {
							t.Errorf("A = %02X, want 00", c.registers.A)
						}
						if !c.registers.FlagZ() || !c.registers.FlagN() {
							t.Errorf("Flags incorrect: Z:%v N:%v", c.registers.FlagZ(), c.registers.FlagN())
						}
						if c.registers.FlagH() || c.registers.FlagCy() {
							t.Error("H and C flags should be clear")
						}
					},
				},
			},
		},
		{
			groupName: "CP (8-bit)",
			cases: []testCase{
				{
					name: "CP A, L (Equal)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0x42
						c.registers.L = 0x42
					},
					run: func(c *CPU) int {
						return c.cpAReg8(c.registers.L)
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
						return c.cpAHLPtr()
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
						return c.cpAImm8(0x20)
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
				{
					name: "CP A, A (Zero)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0x99
					},
					run: func(c *CPU) int {
						return c.cpAReg8(c.registers.A)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if !c.registers.FlagZ() {
							t.Error("Z flag should be set")
						}
						if !c.registers.FlagN() {
							t.Error("N flag should be set")
						}
					},
				},
			},
		},
		{
			groupName: "INC (8-bit)",
			cases: []testCase{
				{
					name: "INC B",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.B = 0x0F
					},
					run: func(c *CPU) int {
						return c.incReg8(&c.registers.B)
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
					name: "INC [HL] (Half-Carry)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0xD001)
						m.data[0xD001] = 0x0F
					},
					run: func(c *CPU) int {
						return c.incHLPtr()
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
			},
		},
		{
			groupName: "DEC (8-bit)",
			cases: []testCase{
				{
					name: "DEC C",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.C = 0x10
					},
					run: func(c *CPU) int {
						return c.decReg8(&c.registers.C)
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
						return c.decHLPtr()
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
			},
		},
		{
			groupName: "16-bit Operations",
			cases: []testCase{
				{
					name: "ADD HL, BC",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0x1000)
						c.registers.SetBC(0x0100)
					},
					run: func(c *CPU) int {
						return c.addHLReg16(c.registers.BC)
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
						return c.addHLReg16(c.registers.DE)
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
					name: "ADD HL, HL (Double HL)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0x4000)
					},
					run: func(c *CPU) int {
						return c.addHLReg16(c.registers.HL)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.HL() != 0x8000 {
							t.Errorf("HL = %04X, want 8000", c.registers.HL())
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
						return c.addHLSP()
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
						return c.decReg16(func(v uint16) { c.registers.SP = v }, func() uint16 { return c.registers.SP })
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.SP != 0xFFFF {
							t.Errorf("SP = %04X, want FFFF", c.registers.SP)
						}
					},
				},
				{
					name: "DEC BC",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetBC(0x0200)
					},
					run: func(c *CPU) int {
						return c.decReg16(c.registers.SetBC, c.registers.BC)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.BC() != 0x01FF {
							t.Errorf("BC = %04X, want 01FF", c.registers.BC())
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
						return c.incReg16(c.registers.SetBC, c.registers.BC)
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
				{
					name: "INC SP",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SP = 0x00FF
					},
					run: func(c *CPU) int {
						return c.incReg16(func(v uint16) { c.registers.SP = v }, func() uint16 { return c.registers.SP })
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.SP != 0x0100 {
							t.Errorf("SP = %04X, want 0100", c.registers.SP)
						}
					},
				},
			},
		},
	}

	for _, group := range tests {
		t.Run(group.groupName, func(t *testing.T) {
			for _, tt := range group.cases {
				t.Run(tt.name, func(t *testing.T) {
					cpu, mem := createTestCPU()
					tt.setup(cpu, mem)
					cycles := tt.run(cpu)
					tt.expected(t, cpu, mem, cycles)
				})
			}
		})
	}
}
