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
	return 0
}

func (m *MMU) Write(addr uint16, data byte) {

}
