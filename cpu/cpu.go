package cpu

// MemoryBus defines the interface for memory access.
// The CPU only needs to know how to Read and Write.
type MemoryBus interface {
	Read(addr uint16) byte
	Write(addr uint16, value byte)
}

type CPU struct {
	registers *Registers
	bus       MemoryBus
}

func (c *CPU) RunNextInstruction() {
	opcode := c.fetch()
	c.execute(opcode)
}

func (c *CPU) fetch() byte {
	value := c.bus.Read(c.registers.pc)
	c.registers.pc++
	return value
}

func (c *CPU) execute(opcode byte) {

}
