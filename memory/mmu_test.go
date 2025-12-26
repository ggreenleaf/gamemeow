package memory

import "testing"

func TestWRAM_ReadWrite(t *testing.T) {
	mmu := &MMU{}

	tests := []struct {
		name string
		addr uint16
		val  byte
	}{
		{"Single Random Address", 0xC005, 0x42},
		{"Start of WRAM", 0xC000, 0xA1},
		{"Middle of WRAM", 0xCFFF, 0xB2},
		{"End of WRAM", 0xDFFF, 0xC3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu.Write(tt.addr, tt.val)
			if got := mmu.Read(tt.addr); got != tt.val {
				t.Errorf("Read(%X) = %X; want %X", tt.addr, got, tt.val)
			}
		})
	}
}

func TestEchoRAM_Mirroring(t *testing.T) {
	mmu := &MMU{}

	// Case 1: Write to WRAM -> Read from Echo
	wramAddr := uint16(0xC000)
	echoAddr := uint16(0xE000)
	val := byte(0x34)

	mmu.Write(wramAddr, val)
	if got := mmu.Read(echoAddr); got != val {
		t.Errorf("Write to WRAM (0xC000) was not mirrored to Echo (0xE000). Got 0x%X, want 0x%X", got, val)
	}

	// Case 2: Write to Echo -> Read from WRAM
	val2 := byte(0x56)
	mmu.Write(echoAddr, val2)
	if got := mmu.Read(wramAddr); got != val2 {
		t.Errorf("Write to Echo (0xE000) was not reflected in WRAM (0xC000). Got 0x%X, want 0x%X", got, val2)
	}
}

func TestVRAM_ReadWrite(t *testing.T) {
	mmu := &MMU{}
	tests := []struct {
		name string
		addr uint16
		val  byte
	}{
		{"VRAM Start", 0x8000, 0x11},
		{"VRAM End", 0x9FFF, 0x22},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu.Write(tt.addr, tt.val)
			if got := mmu.Read(tt.addr); got != tt.val {
				t.Errorf("Read(%X) = %X; want %X", tt.addr, got, tt.val)
			}
		})
	}
}

func TestOAM_ReadWrite(t *testing.T) {
	mmu := &MMU{}
	tests := []struct {
		name string
		addr uint16
		val  byte
	}{
		{"OAM Start", 0xFE00, 0x33},
		{"OAM End", 0xFE9F, 0x44},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu.Write(tt.addr, tt.val)
			if got := mmu.Read(tt.addr); got != tt.val {
				t.Errorf("Read(%X) = %X; want %X", tt.addr, got, tt.val)
			}
		})
	}
}

func TestHRAM_ReadWrite(t *testing.T) {
	mmu := &MMU{}
	tests := []struct {
		name string
		addr uint16
		val  byte
	}{
		{"HRAM Start", 0xFF80, 0x66},
		{"HRAM End", 0xFFFE, 0x77},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu.Write(tt.addr, tt.val)
			if got := mmu.Read(tt.addr); got != tt.val {
				t.Errorf("Read(%X) = %X; want %X", tt.addr, got, tt.val)
			}
		})
	}
}

func TestIO_ReadWrite(t *testing.T) {
	mmu := &MMU{}
	tests := []struct {
		name string
		addr uint16
		val  byte
	}{
		{"IO Start", 0xFF00, 0x55},
		{"IO End", 0xFF7F, 0x99},
		{"IE Register", 0xFFFF, 0x88},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu.Write(tt.addr, tt.val)
			if got := mmu.Read(tt.addr); got != tt.val {
				t.Errorf("Read(%X) = %X; want %X", tt.addr, got, tt.val)
			}
		})
	}
}
