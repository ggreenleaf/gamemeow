package cpu

import "testing"

func TestRegister_BC(t *testing.T) {
	tests := []struct {
		name      string
		setValue  uint16
		expectedB byte
		expectedC byte
	}{
		{"Zero", 0x0000, 0x00, 0x00},
		{"Full", 0xFFFF, 0xFF, 0xFF},
		{"Mixed", 0xAABB, 0xAA, 0xBB},
		{"LowByte", 0x00FF, 0x00, 0xFF},
		{"HighByte", 0xFF00, 0xFF, 0x00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &Registers{}
			reg.SetBC(tt.setValue)

			if reg.b != tt.expectedB {
				t.Errorf("SetBC(0x%X): expected B=0x%X, got 0x%X", tt.setValue, tt.expectedB, reg.b)
			}
			if reg.c != tt.expectedC {
				t.Errorf("SetBC(0x%X): expected C=0x%X, got 0x%X", tt.setValue, tt.expectedC, reg.c)
			}
			if got := reg.BC(); got != tt.setValue {
				t.Errorf("BC(): expected 0x%X, got 0x%X", tt.setValue, got)
			}
		})
	}
}

func TestRegister_DE(t *testing.T) {
	tests := []struct {
		name      string
		setValue  uint16
		expectedD byte
		expectedE byte
	}{
		{"Zero", 0x0000, 0x00, 0x00},
		{"Full", 0xFFFF, 0xFF, 0xFF},
		{"Mixed", 0xCCDD, 0xCC, 0xDD},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &Registers{}
			reg.SetDE(tt.setValue)

			if reg.d != tt.expectedD {
				t.Errorf("SetDE(0x%X): expected D=0x%X, got 0x%X", tt.setValue, tt.expectedD, reg.d)
			}
			if reg.e != tt.expectedE {
				t.Errorf("SetDE(0x%X): expected E=0x%X, got 0x%X", tt.setValue, tt.expectedE, reg.e)
			}
			if got := reg.DE(); got != tt.setValue {
				t.Errorf("DE(): expected 0x%X, got 0x%X", tt.setValue, got)
			}
		})
	}
}

func TestRegister_HL(t *testing.T) {
	tests := []struct {
		name      string
		setValue  uint16
		expectedH byte
		expectedL byte
	}{
		{"Zero", 0x0000, 0x00, 0x00},
		{"Full", 0xFFFF, 0xFF, 0xFF},
		{"Mixed", 0x1234, 0x12, 0x34},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &Registers{}
			reg.SetHL(tt.setValue)

			if reg.h != tt.expectedH {
				t.Errorf("SetHL(0x%X): expected H=0x%X, got 0x%X", tt.setValue, tt.expectedH, reg.h)
			}
			if reg.l != tt.expectedL {
				t.Errorf("SetHL(0x%X): expected L=0x%X, got 0x%X", tt.setValue, tt.expectedL, reg.l)
			}
			if got := reg.HL(); got != tt.setValue {
				t.Errorf("HL(): expected 0x%X, got 0x%X", tt.setValue, got)
			}
		})
	}
}

func TestRegister_AF(t *testing.T) {
	tests := []struct {
		name         string
		setValue     uint16
		expectedA    byte
		expectedF    byte
		expectedRead uint16
	}{
		{
			name:         "Aligns to F0 Mask",
			setValue:     0xFFFF,
			expectedA:    0xFF,
			expectedF:    0xF0, // Lower 4 bits cleared
			expectedRead: 0xFFF0,
		},
		{
			name:         "Zero",
			setValue:     0x0000,
			expectedA:    0x00,
			expectedF:    0x00,
			expectedRead: 0x0000,
		},
		{
			name:         "Preserves Upper Nibble",
			setValue:     0x12F3, // F3 = 1111 0011
			expectedA:    0x12,
			expectedF:    0xF0, // 3 is masked out
			expectedRead: 0x12F0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &Registers{}
			reg.SetAF(tt.setValue)

			if reg.a != tt.expectedA {
				t.Errorf("SetAF(0x%X): expected A=0x%X, got 0x%X", tt.setValue, tt.expectedA, reg.a)
			}
			if reg.f != tt.expectedF {
				t.Errorf("SetAF(0x%X): expected F=0x%X, got 0x%X", tt.setValue, tt.expectedF, reg.f)
			}
			if got := reg.AF(); got != tt.expectedRead {
				t.Errorf("AF(): expected 0x%X, got 0x%X", tt.expectedRead, got)
			}
		})
	}
}

func TestRegister_Flags(t *testing.T) {
	// Table of flag setter/getter pairs
	// We can't easily iterate functions, but we can iterate scenarios
	type flagOp func(*Registers) bool
	type flagSet func(*Registers, bool)

	tests := []struct {
		name   string
		setter flagSet
		getter flagOp
	}{
		{"Z", (*Registers).SetZ, (*Registers).Z},
		{"N", (*Registers).SetN, (*Registers).N},
		{"H", (*Registers).SetH, (*Registers).H},
		{"Cy", (*Registers).SetCy, (*Registers).Cy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &Registers{}

			// Initially false
			if tt.getter(reg) {
				t.Errorf("Flag %s initialized to true, expected false", tt.name)
			}

			// Set true
			tt.setter(reg, true)
			if !tt.getter(reg) {
				t.Errorf("Flag %s not set to true", tt.name)
			}

			// Set false
			tt.setter(reg, false)
			if tt.getter(reg) {
				t.Errorf("Flag %s not set to false", tt.name)
			}
		})
	}
}
