package cpu

// General Utilities
func (c *CPU) doAdvanceProgramCounter() {
	// Increase the program counter by 2 (which is the size of an operationCode)
	c.programCounter += 2
}

// 0x0*** Instructions
// 0x1*** Instructions
// 0x2*** Instructions
func (c *CPU) do2NNN(operationCode uint16) {
	// Store the current location of the program counter to the stack
	// And increment the stack pointer
	// to insert a new entry to the stack with the last progarm counter posiiton
	c.stack[c.stackPointer] = c.programCounter
	c.stackPointer++
	// change the program counter to the current subroutine location
	c.programCounter = operationCode & 0x0FFF
}

// 0x3*** Instructions
// 0x4*** Instructions
// 0x5*** Instructions
// 0x6*** Instructions
// 0x7*** Instructions
// 0x8*** Instructions
func (c *CPU) do8XY4(operationCode uint16) {
	// Get the value of X and Y from the operationCode
	registerXLocation, registerYLocation := (operationCode&0x0F00)>>8, (operationCode&0x00F0)>>4
	// Get the value of X + Y
	result := c.register[registerXLocation] + c.register[registerYLocation]
	// If there is an overflow, set the overflow flag on register VF (0xF) to 1
	if c.register[registerXLocation]+c.register[registerYLocation] > 0xFF {
		c.register[0xF] = 1
	} else {
		c.register[0xF] = 0
	}
	// Set the result
	c.register[registerXLocation] = result
	c.doAdvanceProgramCounter()
}

// 0x9*** Instructions
// 0xA*** Instructions
func (c *CPU) doANNN(operationCode uint16) {
	// Get last 3 * 4 bits of the operation code
	c.indexRegister = operationCode & 0x0FFF
	c.doAdvanceProgramCounter()
}

// 0xB*** Instructions
// 0xC*** Instructions
// 0xD*** Instructions
// 0xE*** Instructions
// 0xF*** Instructions
func (c *CPU) doFX33(operationCode uint16) {
	// Get the value at register X
	registerXValue := c.register[(operationCode&0x0F00)>>8]
	// Set the hundred's value of x to memory[I]
	c.memory[c.indexRegister] = registerXValue / 100
	// Set the ten's value of x to memory[I+1]
	c.memory[c.indexRegister+1] = (registerXValue / 10) % 10
	// Set the one's value of x to memory[I+2]
	c.memory[c.indexRegister+2] = (registerXValue % 100) % 10
	c.doAdvanceProgramCounter()
}
