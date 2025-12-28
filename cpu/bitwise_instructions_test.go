package cpu

import (
	"testing"
)

func TestBitwiseInstructions(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(c *CPU, m *mockMemory)
		run      func(c *CPU) int
		expected func(t *testing.T, c *CPU, m *mockMemory, cycles int)
	}{
		// --- AND ---
		{
			name: "AND A, B (Result Zero)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x0F
				c.registers.B = 0xF0
			},
			run: func(c *CPU) int {
				return c.andAReg8(func() byte { return c.registers.B })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x00 {
					t.Errorf("A = %02X, want 00", c.registers.A)
				}
				if !c.registers.FlagZ() {
					t.Error("Z flag should be set")
				}
				if !c.registers.FlagH() {
					t.Error("H flag should be set for AND")
				}
				if c.registers.FlagN() || c.registers.FlagCy() {
					t.Errorf("Flags incorrect: N:%v C:%v", c.registers.FlagN(), c.registers.FlagCy())
				}
				if cycles != 1 {
					t.Errorf("Cycles = %d, want 1", cycles)
				}
			},
		},
		{
			name: "AND A, [HL]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0xFF
				c.registers.SetHL(0xC000)
				m.data[0xC000] = 0xAA
			},
			run: func(c *CPU) int {
				return c.andAHLPtr()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0xAA {
					t.Errorf("A = %02X, want AA", c.registers.A)
				}
				if !c.registers.FlagH() {
					t.Error("H flag should be set")
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},
		{
			name: "AND A, n8",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x55
			},
			run: func(c *CPU) int {
				return c.andAImm8(0x55)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x55 {
					t.Errorf("A = %02X, want 55", c.registers.A)
				}
				if c.registers.FlagZ() {
					t.Error("Z flag should not be set")
				}
				if !c.registers.FlagH() {
					t.Error("H flag should be set")
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},

		// --- XOR ---
		{
			name: "XOR A, A (Zero)",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0xFF
			},
			run: func(c *CPU) int {
				return c.xorAReg8(func() byte { return c.registers.A })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x00 {
					t.Errorf("A = %02X, want 00", c.registers.A)
				}
				if !c.registers.FlagZ() {
					t.Error("Z flag should be set")
				}
				if c.registers.FlagN() || c.registers.FlagH() || c.registers.FlagCy() {
					t.Errorf("Flags incorrect: N:%v H:%v C:%v", c.registers.FlagN(), c.registers.FlagH(), c.registers.FlagCy())
				}
				if cycles != 1 {
					t.Errorf("Cycles = %d, want 1", cycles)
				}
			},
		},
		{
			name: "XOR A, [HL]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x0F
				c.registers.SetHL(0xD000)
				m.data[0xD000] = 0xF0
			},
			run: func(c *CPU) int {
				return c.xorAHLPtr()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0xFF {
					t.Errorf("A = %02X, want FF", c.registers.A)
				}
				if c.registers.FlagZ() {
					t.Error("Z flag should not be set")
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},
		{
			name: "XOR A, n8",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0xAA
			},
			run: func(c *CPU) int {
				return c.xorAImm8(0x55)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0xFF {
					t.Errorf("A = %02X, want FF", c.registers.A)
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},

		// --- OR ---
		{
			name: "OR A, B",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0xF0
				c.registers.B = 0x0F
			},
			run: func(c *CPU) int {
				return c.orAReg8(func() byte { return c.registers.B })
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0xFF {
					t.Errorf("A = %02X, want FF", c.registers.A)
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
			name: "OR A, [HL]",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0x00
				c.registers.SetHL(0xC000)
				m.data[0xC000] = 0x00
			},
			run: func(c *CPU) int {
				return c.orAHLPtr()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x00 {
					t.Errorf("A = %02X, want 00", c.registers.A)
				}
				if !c.registers.FlagZ() {
					t.Error("Z flag should be set")
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},
		{
			name: "OR A, n8",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0xAA
			},
			run: func(c *CPU) int {
				return c.orAImm8(0x55)
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0xFF {
					t.Errorf("A = %02X, want FF", c.registers.A)
				}
				if cycles != 2 {
					t.Errorf("Cycles = %d, want 2", cycles)
				}
			},
		},
		{
			name: "CPL A",
			setup: func(c *CPU, m *mockMemory) {
				c.registers.A = 0xFF
			},
			run: func(c *CPU) int {
				return c.cpl()
			},
			expected: func(t *testing.T, c *CPU, m *mockMemory, cycles int) {
				if c.registers.A != 0x00 {
					t.Errorf("A = %02X, want 00", c.registers.A)
				}
				if !c.registers.FlagN() || !c.registers.FlagH() {
					t.Errorf("Flags incorrect: N:%v H:%v", c.registers.FlagN(), c.registers.FlagH())
				}
				if cycles != 1 {
					t.Errorf("Cycles = %d, want 1", cycles)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, mem := createTestCPU()
			tt.setup(cpu, mem)
			cycle := tt.run(cpu)
			tt.expected(t, cpu, mem, cycle)
		})
	}
}
