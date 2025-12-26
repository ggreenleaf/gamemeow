package memory

// Memory Map Constants
// Source: https://gbdev.io/pandocs/Memory_Map.html
const (
	CartridgeROMEnd = 0x7FFF

	VRAMStart = 0x8000
	VRAMEnd   = 0x9FFF

	CartridgeRAMStart = 0xA000
	CartridgeRAMEnd   = 0xBFFF

	WRAMStart = 0xC000
	WRAMEnd   = 0xDFFF

	EchoRAMStart = 0xE000
	EchoRAMEnd   = 0xFDFF

	OAMStart = 0xFE00
	OAMEnd   = 0xFE9F

	UnusableStart = 0xFEA0
	UnusableEnd   = 0xFEFF

	IOStart = 0xFF00
	IOEnd   = 0xFF7F

	HRAMStart = 0xFF80
	HRAMEnd   = 0xFFFE

	IEAddr = 0xFFFF
)

type Cartridge interface {
	Read(addr uint16) byte
	Write(addr uint16, data byte)
}

// MMU (Memory Management Unit)
// In our emulator, this struct acts as both the "Address Decoder" (routing requests)
// and the "Storage Container" (holding the actual byte slices for WRAM, VRAM, etc).
type MMU struct {
	cartridge Cartridge
	vram      [8192]byte
	wram      [8192]byte
	oam       [160]byte
	hram      [127]byte

	// interrupt enable register
	ie byte

	// temp for now
	io [128]byte
}

func (m *MMU) Read(addr uint16) byte {
	switch {
	case addr <= CartridgeROMEnd:
		return m.cartridge.Read(addr)

	case addr >= VRAMStart && addr <= VRAMEnd:
		return m.vram[addr-VRAMStart]

	case addr >= CartridgeRAMStart && addr <= CartridgeRAMEnd:
		return m.cartridge.Read(addr)

	case addr >= WRAMStart && addr <= WRAMEnd:
		return m.wram[addr-WRAMStart]

	case addr >= EchoRAMStart && addr <= EchoRAMEnd:
		// Echo RAM is a mirror of WRAM (0xC000 - 0xDDFF)
		// We subtract EchoStart (0xE000) to get the offset, effectively mapping
		// 0xE000 -> 0x0000 (Start of WRAM)
		return m.wram[addr-EchoRAMStart]

	case addr >= OAMStart && addr <= OAMEnd:
		return m.oam[addr-OAMStart]

	case addr >= UnusableStart && addr <= UnusableEnd:
		return 0xFF

	case addr >= IOStart && addr <= IOEnd:
		return m.io[addr-IOStart]

	case addr >= HRAMStart && addr <= HRAMEnd:
		return m.hram[addr-HRAMStart]

	case addr == IEAddr:
		return m.ie
	}
	return 0xFF
}

func (m *MMU) Write(addr uint16, data byte) {
	switch {
	case addr <= CartridgeROMEnd:
		m.cartridge.Write(addr, data)

	case addr >= VRAMStart && addr <= VRAMEnd:
		m.vram[addr-VRAMStart] = data

	case addr >= CartridgeRAMStart && addr <= CartridgeRAMEnd:
		m.cartridge.Write(addr, data)

	case addr >= WRAMStart && addr <= WRAMEnd:
		m.wram[addr-WRAMStart] = data

	case addr >= EchoRAMStart && addr <= EchoRAMEnd:
		// Echo RAM is a mirror of WRAM (0xC000 - 0xDDFF)
		// We subtract EchoStart (0xE000) to get the offset, effectively mapping
		// 0xE000 -> 0x0000 (Start of WRAM)
		m.wram[addr-EchoRAMStart] = data

	case addr >= OAMStart && addr <= OAMEnd:
		m.oam[addr-OAMStart] = data

	case addr >= UnusableStart && addr <= UnusableEnd:
		return

	case addr >= IOStart && addr <= IOEnd:
		m.io[addr-IOStart] = data

	case addr >= HRAMStart && addr <= HRAMEnd:
		m.hram[addr-HRAMStart] = data

	case addr == IEAddr:
		m.ie = data
	}
}
