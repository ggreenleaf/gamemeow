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
	tests := []struct {
		name     string
		setup    func(c *CPU, m *mockMemory)
		run      func(c *CPU) int
		expected func(t *testing.T, c *CPU, m *mockMemory, cycles int)
	}{
		{
			name: "LD A, B",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.B = 0x42
			},
			run: func(c *CPU) int {
				return c.loadReg8Reg8(func(v byte) { c.registers.A = v }, func() byte { return c.registers.B })
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
				return c.loadReg8Imm8(func(v byte) { c.registers.B = v }, 0x99)
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
			name: "LD [HL], A",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetHL(0xC000)
				c.registers.A = 0x77
			},
			run: func(c *CPU) int {
				return c.storeHLReg8(c.registers.A)
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
			name: "LDH [a8], A",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0xAA
			},
			run: func(c *CPU) int {
				return c.storeHighImm16A(0x10) // LDH [$FF10], A
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
			name: "LD [HL+], A",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetHL(0xD000)
				c.registers.A = 0x55
			},
			run: func(c *CPU) int {
				return c.storeHLIncA()
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
			name: "LD A, [HL+]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetHL(0xC0F0)
				m.data[0xC0F0] = 0x11
			},
			run: func(c *CPU) int {
				return c.loadAHLInc()
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
				return c.loadAHLDec()
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
				return c.storeHLDecA()
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
			name: "LD B, [HL]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.SetHL(0xC200)
				m.data[0xC200] = 0x99
			},
			run: func(c *CPU) int {
				return c.loadReg8HL(func(v byte) { c.registers.B = v })
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
				return c.storeHLImm8(0x33)
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
