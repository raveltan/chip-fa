package cpu

import "math/rand"

// General Utilities
func (c *CPU) doAdvanceProgramCounter() {
	// Increase the program counter by 2 (which is the size of an operationCode)
	c.programCounter += 2
}

// 0x0*** Instructions
func (c *CPU) do00E0() {
	for i := range c.Screen {
		c.Screen[i] = 0
	}
	c.ShouldDraw = true
	c.doAdvanceProgramCounter()
}
func (c *CPU) do000E() {
	// Decrease the stack pointer to the previous one
	c.stackPointer--
	// Re-assign the program counter to the program counter on the previous stack
	c.programCounter = c.stack[c.stackPointer]

	c.doAdvanceProgramCounter()
}

// 0x1*** Instructions
func (c *CPU) do1NNN(operationCode uint16) {
	c.programCounter = (operationCode & 0x0FFF)
}

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
func (c *CPU) do3XNN(operationCode uint16) {
	// If value of register X == NN skip next instruction
	if c.register[(operationCode&0x0F00)>>8] == (uint8(operationCode & 0x00FF)) {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0x4*** Instructions
func (c *CPU) do4XNN(operationCode uint16) {
	// If value of register X != NN skip next instruction
	if c.register[(operationCode&0x0F00)>>8] != (uint8(operationCode & 0x00FF)) {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0x5*** Instructions
func (c *CPU) do5XY0(operationCode uint16) {
	// If value of register X == value of register Y skip next instruction
	if c.register[(operationCode&0x0F00)>>8] == c.register[(operationCode&0x00F0)>>4] {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0x6*** Instructions
func (c *CPU) do6XNN(operationCode uint16) {
	c.register[(operationCode&0x0F00)>>8] = uint8(operationCode & 0x00FF)
	c.doAdvanceProgramCounter()
}

// 0x7*** Instructions
func (c *CPU) do7XNN(operationCode uint16) {
	c.register[(operationCode&0x0F00)>>8] += uint8(operationCode & 0x00FF)
	c.doAdvanceProgramCounter()
}

// 0x8*** Instructions
func (c *CPU) do8XY0(operationCode uint16) {
	c.register[(operationCode&0x0F00)>>8] = c.register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY1(operationCode uint16) {
	c.register[(operationCode&0x0F00)>>8] = c.register[(operationCode&0x0F00)>>8] | c.register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY2(operationCode uint16) {
	c.register[(operationCode&0x0F00)>>8] = c.register[(operationCode&0x0F00)>>8] & c.register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY3(operationCode uint16) {
	c.register[(operationCode&0x0F00)>>8] = c.register[(operationCode&0x0F00)>>8] ^ c.register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}

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

func (c *CPU) do8XY5(operationCode uint16) {
	if c.register[(operationCode&0x0F00)>>8] < c.register[(operationCode&0x00F0)>>4] {
		c.register[0xF] = 0
	} else {
		c.register[0xF] = 1
	}
	c.register[(operationCode&0x0F00)>>8] -= c.register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}
func (c *CPU) do8XY6(operationCode uint16) {
	c.register[0xF] = c.register[(operationCode&0x0F00)>>8] & 0x1
	c.register[(operationCode&0x0F00)>>8] >>= 1
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY7(operationCode uint16) {
	if c.register[(operationCode&0x0F00)>>8] > c.register[(operationCode&0x00F0)>>4] {
		c.register[0xF] = 0
	} else {
		c.register[0xF] = 1
	}
	c.register[(operationCode&0x0F00)>>8] = c.register[(operationCode&0x00F0)>>4] - c.register[(operationCode&0x0F00)>>8]
	c.doAdvanceProgramCounter()
}
func (c *CPU) do8XYE(operationCode uint16) {
	c.register[0xF] = c.register[(operationCode&0x0F00)>>8] >> 7
	c.register[(operationCode&0x0F00)<<8] <<= 1
	c.doAdvanceProgramCounter()
}

// 0x9*** Instructions
func (c *CPU) do9XY0(operationCode uint16) {
	if c.register[(operationCode&0x0F00)>>8] != c.register[(operationCode&0x00F0)>>4] {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0xA*** Instructions
func (c *CPU) doANNN(operationCode uint16) {
	// Get last 3 * 4 bits of the operation code
	c.indexRegister = operationCode & 0x0FFF
	c.doAdvanceProgramCounter()
}

// 0xB*** Instructions
func (c *CPU) doBNNN(operationCode uint16) {
	c.programCounter = (operationCode & 0x0FFF) + uint16(c.register[0])
}

// 0xC*** Instructions
func (c *CPU) doCXNN(operationCode uint16) {
	c.register[(operationCode&0x0F00)>>8] = uint8(rand.Intn(0xFF)) & uint8(operationCode&0x00FF)
	c.doAdvanceProgramCounter()
}

// 0xD*** Instructions
func (c *CPU) doDXYN(operationCode uint16) {
	x, y, h := c.register[(operationCode&0x0F00)>>8], c.register[(operationCode&0x00F0)>>4], operationCode&0x000F
	pixelData := uint8(0)

	// Track if any pixels are flipped from set to unset.
	c.register[0xF] = 0
	for yline := 0; yline < int(h); yline++ {
		pixelData = c.memory[c.indexRegister+uint16(yline)]

		for xline := 0; xline < 8; xline++ {
			if (pixelData & (0x80 >> xline)) != 0 {
				// FIXME: runtime error: index out of range [2101] with length 2048
				if c.Screen[(int(x)+xline+((int(y)+yline)*64))] == 1 {
					c.register[0xF] = 1
				}
				c.Screen[int(x)+xline+((int(y)+yline)*64)] ^= 1
			}
		}
	}

	c.ShouldDraw = true
	c.doAdvanceProgramCounter()
}

// 0xE*** Instructions
func (c *CPU) doEX9E(operationCode uint16) {
	if c.KeypadStates[c.register[(operationCode&0x0F00)>>8]] != 0 {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}
func (c *CPU) doEXA1(operationCode uint16) {
	if c.KeypadStates[c.register[(operationCode&0x0F00)>>8]] == 0 {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0xF*** Instructions
func (c *CPU) doFX07(operationCode uint16) {
	c.register[(operationCode&0x0F00)>>8] = c.delayTimer
	c.doAdvanceProgramCounter()
}

func (c *CPU) doFX0A(operationCode uint16) {
	pressed := false
	for i, v := range c.KeypadStates {
		if v != 0 {
			c.register[(operationCode&0x0F00)>>8] = uint8(i)
			pressed = true
		}
	}

	if !pressed {
		return
	}
	c.doAdvanceProgramCounter()
}

func (c *CPU) doFX15(operationCode uint16) {
	c.delayTimer = c.register[(operationCode&0x0F00)>>8]
	c.doAdvanceProgramCounter()
}
func (c *CPU) doFX18(operationCode uint16) {
	c.SoundTimer = c.register[(operationCode&0x0F00)>>8]
	c.doAdvanceProgramCounter()
}
func (c *CPU) doFX1E(operationCode uint16) {
	if c.indexRegister+uint16(c.register[(operationCode&0x0F00)>>8]) > 0xFFF {
		c.register[0xF] = 1
	} else {
		c.register[0xF] = 0
	}
	c.indexRegister += uint16(c.register[(operationCode&0x0F00)>>8])
	c.doAdvanceProgramCounter()
}

func (c *CPU) doFX29(operationCode uint16) {
	c.indexRegister = uint16(c.register[(operationCode&0x0F00)>>8]) * 0x5
	c.doAdvanceProgramCounter()
}
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

func (c *CPU) doFX55(operationCode uint16) {
	for i := 0; i <= int((operationCode&0x0F00)>>8); i++ {
		c.memory[int(c.indexRegister)+i] = c.register[i]
	}

	// On the original system
	// When the operation is done
	// indexRegistered += X + 1

	c.indexRegister += (operationCode&0x0F00)>>8 + 1
	c.doAdvanceProgramCounter()
}

func (c *CPU) doFX65(operationCode uint16) {
	for i := 0; i <= int((operationCode&0x0F00)>>8); i++ {
		c.register[i] = c.memory[int(c.indexRegister)+i]
	}

	// On the original interpreter, when the operation is done, I = I + X + 1.
	c.indexRegister += (operationCode&0x0F00)>>8 + 1
	c.doAdvanceProgramCounter()
}
