package cpu

import "testing"

type mockMemory struct {
	data map[uint16]byte
}

func (m *mockMemory) Read(addr uint16) byte {
	return m.data[addr]
}

func (m *mockMemory) Write(addr uint16, value byte) {
	if m.data == nil {
		m.data = make(map[uint16]byte)
	}
	m.data[addr] = value
}

func TestCPU_Fetch(t *testing.T) {
	mockBus := &mockMemory{
		data: map[uint16]byte{
			0x0000: 0x42, // Random instruction byte
		},
	}

	cpu := &CPU{
		registers: &Registers{pc: 0x0000},
		bus:       mockBus,
	}

	instruction := cpu.fetch()
	if instruction != 0x42 {
		t.Errorf("fetch() = 0x%X; want 0x42", instruction)
	}

	if cpu.registers.pc != 0x0001 {
		t.Errorf("After fetch, PC = 0x%X; want 0x0001", cpu.registers.pc)
	}
}
