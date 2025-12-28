package cpu

import "testing"

// createTestCPU helper creates a CPU with a mock memory bus for testing.
// It returns the CPU and the mock memory so we can inspect/write to it.
func createTestCPU() (*CPU, *mockMemory) {
	mem := &mockMemory{
		data: make(map[uint16]byte),
	}
	regs := &Registers{}
	cpu := &CPU{
		registers: regs,
		bus:       mem,
	}
	return cpu, mem
}

func TestLoadInstructions(t *testing.T) {
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
			groupName: "8-bit Loads",
			cases: []testCase{
				{
					name: "LD A, B",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.B = 0x42
					},
					run: func(c *CPU) int {
						return c.loadReg8Reg8(&c.registers.A, c.registers.B)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x42 {
							t.Errorf("A = %X, want 0x42", c.registers.A)
						}
						if cycles != 1 {
							t.Errorf("Cycles = %d, want 1", cycles)
						}
					},
				},
				{
					name:  "LD B, n8",
					setup: func(c *CPU, m *mockMemory) {}, // No setup needed
					run: func(c *CPU) int {
						return c.loadReg8Imm8(&c.registers.B, 0x99)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.B != 0x99 {
							t.Errorf("B = %X, want 0x99", c.registers.B)
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LD C, C (No-op check)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.C = 0xAA
					},
					run: func(c *CPU) int {
						return c.loadReg8Reg8(&c.registers.C, c.registers.C)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.C != 0xAA {
							t.Errorf("C = %X, want 0xAA", c.registers.C)
						}
						if cycles != 1 {
							t.Errorf("Cycles = %d, want 1", cycles)
						}
					},
				},
				{
					name:  "LD E, n8",
					setup: func(c *CPU, m *mockMemory) {},
					run: func(c *CPU) int {
						return c.loadReg8Imm8(&c.registers.E, 0x12)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.E != 0x12 {
							t.Errorf("E = %X, want 0x12", c.registers.E)
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
			},
		},
		{
			groupName: "16-bit Loads",
			cases: []testCase{
				{
					name:  "LD BC, n16",
					setup: func(c *CPU, m *mockMemory) {},
					run: func(c *CPU) int {
						return c.loadReg16Imm16(c.registers.SetBC, 0x1234)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.BC() != 0x1234 {
							t.Errorf("BC = %X, want 0x1234", c.registers.BC())
						}
						if cycles != 3 {
							t.Errorf("Cycles = %d, want 3", cycles)
						}
					},
				},
				{
					name:  "LD DE, n16",
					setup: func(c *CPU, m *mockMemory) {},
					run: func(c *CPU) int {
						return c.loadReg16Imm16(c.registers.SetDE, 0x5678)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.DE() != 0x5678 {
							t.Errorf("DE = %X, want 0x5678", c.registers.DE())
						}
						if cycles != 3 {
							t.Errorf("Cycles = %d, want 3", cycles)
						}
					},
				},
				{
					name:  "LD HL, n16",
					setup: func(c *CPU, m *mockMemory) {},
					run: func(c *CPU) int {
						return c.loadReg16Imm16(c.registers.SetHL, 0x9ABC)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.HL() != 0x9ABC {
							t.Errorf("HL = %X, want 0x9ABC", c.registers.HL())
						}
						if cycles != 3 {
							t.Errorf("Cycles = %d, want 3", cycles)
						}
					},
				},
				{
					name:  "LD SP, n16",
					setup: func(c *CPU, m *mockMemory) {},
					run: func(c *CPU) int {
						return c.loadSPImm16(0xDEAD)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.SP != 0xDEAD {
							t.Errorf("SP = %X, want 0xDEAD", c.registers.SP)
						}
						if cycles != 3 {
							t.Errorf("Cycles = %d, want 3", cycles)
						}
					},
				},
				{
					name: "LD SP, HL",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0xF00D)
					},
					run: func(c *CPU) int {
						return c.loadSPHL()
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.SP != 0xF00D {
							t.Errorf("SP = %X, want 0xF00D", c.registers.SP)
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
			},
		},
		{
			groupName: "Indirect Loads",
			cases: []testCase{
				{
					name: "LD [HL], A",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0xC000)
						c.registers.A = 0x77
					},
					run: func(c *CPU) int {
						return c.storeHLPtrReg8(c.registers.A)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xC000] != 0x77 {
							t.Errorf("Mem[C000] = %X, want 0x77", m.data[0xC000])
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LD B, [HL]",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0xC200)
						m.data[0xC200] = 0x99
					},
					run: func(c *CPU) int {
						return c.loadReg8HLPtr(&c.registers.B)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.B != 0x99 {
							t.Errorf("B = %X, want 0x99", c.registers.B)
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LD [HL], n8",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0xC300)
					},
					run: func(c *CPU) int {
						return c.storeHLPtrImm8(0x33)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xC300] != 0x33 {
							t.Errorf("Mem[C300] = %X, want 0x33", m.data[0xC300])
						}
						if cycles != 3 {
							t.Errorf("Cycles = %d, want 3", cycles)
						}
					},
				},
				{
					name: "LD A, [n16]",
					setup: func(c *CPU, m *mockMemory) {
						m.data[0xC123] = 0x55
					},
					run: func(c *CPU) int {
						return c.loadAImm16(0xC123)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x55 {
							t.Errorf("A = %X, want 0x55", c.registers.A)
						}
						if cycles != 4 {
							t.Errorf("Cycles = %d, want 4", cycles)
						}
					},
				},
				{
					name: "LD [n16], A",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0x88
					},
					run: func(c *CPU) int {
						return c.storeImm16A(0xC001)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xC001] != 0x88 {
							t.Errorf("Mem[C001] = %X, want 0x88", m.data[0xC001])
						}
						if cycles != 4 {
							t.Errorf("Cycles = %d, want 4", cycles)
						}
					},
				},
				{
					name: "LD A, [BC]",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetBC(0xCABC)
						m.data[0xCABC] = 0x66
					},
					run: func(c *CPU) int {
						return c.loadAReg16(c.registers.BC)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x66 {
							t.Errorf("A = %X, want 0x66", c.registers.A)
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LD A, [DE]",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetDE(0xDEDE)
						m.data[0xDEDE] = 0x77
					},
					run: func(c *CPU) int {
						return c.loadAReg16(c.registers.DE)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x77 {
							t.Errorf("A = %X, want 0x77", c.registers.A)
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LD [BC], A",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetBC(0xB00B)
						c.registers.A = 0x88
					},
					run: func(c *CPU) int {
						return c.storeReg16A(c.registers.BC)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xB00B] != 0x88 {
							t.Errorf("Mem[B00B] = %X, want 0x88", m.data[0xB00B])
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LD [DE], A",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetDE(0xD00D)
						c.registers.A = 0x99
					},
					run: func(c *CPU) int {
						return c.storeReg16A(c.registers.DE)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xD00D] != 0x99 {
							t.Errorf("Mem[D00D] = %X, want 0x99", m.data[0xD00D])
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
			},
		},
		{
			groupName: "High Ram (LDH) Loads",
			cases: []testCase{
				{
					name: "LDH [a8], A",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.A = 0xAA
					},
					run: func(c *CPU) int {
						return c.storeHighImm8A(0x10) // LDH [$FF10], A
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xFF10] != 0xAA {
							t.Errorf("Mem[FF10] = %X, want 0xAA", m.data[0xFF10])
						}
						if cycles != 3 {
							t.Errorf("Cycles = %d, want 3", cycles)
						}
					},
				},
				{
					name: "LDH A, [C]",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.C = 0x44
						m.data[0xFF44] = 0xBB // Accessing LCD register as an example
					},
					run: func(c *CPU) int {
						return c.loadHighAC()
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0xBB {
							t.Errorf("A = %X, want 0xBB", c.registers.A)
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LDH [C], A",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.C = 0x05
						c.registers.A = 0x77
					},
					run: func(c *CPU) int {
						return c.storeHighCA()
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xFF05] != 0x77 {
							t.Errorf("Mem[FF05] = %X, want 0x77", m.data[0xFF05])
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LDH A, [n8]",
					setup: func(c *CPU, m *mockMemory) {
						m.data[0xFF80] = 0x44
					},
					run: func(c *CPU) int {
						return c.loadHighAImm8(0x80)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x44 {
							t.Errorf("A = %X, want 0x44", c.registers.A)
						}
						if cycles != 3 {
							t.Errorf("Cycles = %d, want 3", cycles)
						}
					},
				},
				{
					name: "LDH A, [n8] (Min Offset 0x00)",
					setup: func(c *CPU, m *mockMemory) {
						m.data[0xFF00] = 0x01
					},
					run: func(c *CPU) int {
						return c.loadHighAImm8(0x00)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x01 {
							t.Errorf("A = %X, want 0x01", c.registers.A)
						}
						if cycles != 3 {
							t.Errorf("Cycles = %d, want 3", cycles)
						}
					},
				},
				{
					name: "LDH A, [n8] (Max Offset 0xFF)",
					setup: func(c *CPU, m *mockMemory) {
						m.data[0xFFFF] = 0x0F
					},
					run: func(c *CPU) int {
						return c.loadHighAImm8(0xFF)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x0F {
							t.Errorf("A = %X, want 0x0F", c.registers.A)
						}
						if cycles != 3 {
							t.Errorf("Cycles = %d, want 3", cycles)
						}
					},
				},
			},
		},
		{
			groupName: "HL Increment/Decrement",
			cases: []testCase{
				{
					name: "LD [HL+], A",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0xD000)
						c.registers.A = 0x55
					},
					run: func(c *CPU) int {
						return c.storeHLPtrIncA()
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xD000] != 0x55 {
							t.Errorf("Mem[D000] = %X, want 0x55", m.data[0xD000])
						}
						if c.registers.HL() != 0xD001 {
							t.Errorf("HL = %X, want 0xD001", c.registers.HL())
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LD A, [HL+]",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0xC0F0)
						m.data[0xC0F0] = 0x11
					},
					run: func(c *CPU) int {
						return c.loadAHLPtrInc()
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x11 {
							t.Errorf("A = %X, want 0x11", c.registers.A)
						}
						if c.registers.HL() != 0xC0F1 {
							t.Errorf("HL = %X, want 0xC0F1", c.registers.HL())
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LD A, [HL-]",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0xC0F0)
						m.data[0xC0F0] = 0x22
					},
					run: func(c *CPU) int {
						return c.loadAHLPtrDec()
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.A != 0x22 {
							t.Errorf("A = %X, want 0x22", c.registers.A)
						}
						if c.registers.HL() != 0xC0EF {
							t.Errorf("HL = %X, want 0xC0EF", c.registers.HL())
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
				{
					name: "LD [HL-], A",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SetHL(0xD500)
						c.registers.A = 0x33
					},
					run: func(c *CPU) int {
						return c.storeHLPtrDecA()
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xD500] != 0x33 {
							t.Errorf("Mem[D500] = %X, want 0x33", m.data[0xD500])
						}
						if c.registers.HL() != 0xD4FF {
							t.Errorf("HL = %X, want 0xD4FF", c.registers.HL())
						}
						if cycles != 2 {
							t.Errorf("Cycles = %d, want 2", cycles)
						}
					},
				},
			},
		},
		{
			groupName: "Stack Pointer Operations",
			cases: []testCase{
				{
					name: "LD HL, SP+e8 (Positive offset)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SP = 0x1000
					},
					run: func(c *CPU) int {
						return c.loadHLSPSigned8(0x05)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.HL() != 0x1005 {
							t.Errorf("HL = %X, want 0x1005", c.registers.HL())
						}
						// Carry from bit 3: (0x00 & 0x0F) + (0x05 & 0x0F) = 0x05 (No carry)
						// Carry from bit 7: (0x00 & 0xFF) + 0x05 = 0x05 (No carry)
						if c.registers.FlagH() {
							t.Errorf("Flag H should be false")
						}
						if c.registers.FlagCy() {
							t.Errorf("Flag Cy should be false")
						}
					},
				},
				{
					name: "LD HL, SP+e8 (Negative offset & Flags)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SP = 0x1002
					},
					run: func(c *CPU) int {
						return c.loadHLSPSigned8(-3) // 0xFD
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.HL() != 0x0FFF {
							t.Errorf("HL = %X, want 0x0FFF", c.registers.HL())
						}
						// SP low byte = 0x02, Offset bits = 0xFD
						// H: (0x2 & 0xF) + (0xD & 0xF) = 0x2 + 0xD = 0xF (No carry)
						// Cy: 0x02 + 0xFD = 0xFF (No carry)
						if c.registers.FlagH() {
							t.Error("H should be false")
						}
						if c.registers.FlagCy() {
							t.Error("Cy should be false")
						}
					},
				},
				{
					name: "LD HL, SP+e8 (H-Carry check)",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SP = 0x100F
					},
					run: func(c *CPU) int {
						return c.loadHLSPSigned8(0x01)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if c.registers.HL() != 0x1010 {
							t.Errorf("HL = %X, want 0x1010", c.registers.HL())
						}
						// H: 0xF + 0x1 = 0x10 (Carry!)
						if !c.registers.FlagH() {
							t.Error("H should be true")
						}
						if c.registers.FlagCy() {
							t.Error("Cy should be false")
						}
					},
				},
				{
					name: "LD [n16], SP",
					setup: func(c *CPU, m *mockMemory) {
						c.registers.SP = 0xABCD
					},
					run: func(c *CPU) int {
						return c.storeImm16SP(0xC000)
					},
					expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
						if m.data[0xC000] != 0xCD {
							t.Errorf("Mem[C000] = %X, want 0xCD", m.data[0xC000])
						}
						if m.data[0xC001] != 0xAB {
							t.Errorf("Mem[C001] = %X, want 0xAB", m.data[0xC001])
						}
						if cycles != 5 {
							t.Errorf("Cycles = %d, want 5", cycles)
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
