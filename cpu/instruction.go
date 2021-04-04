package cpu

import "math/rand"

// General Utilities
func (c *CPU) doAdvanceProgramCounter() {
	// Increase the program counter by 2 (which is the size of an operationCode)
	c.ProgramCounter += 2
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
	c.StackPointer--
	// Re-assign the program counter to the program counter on the previous stack
	c.ProgramCounter = c.Stack[c.StackPointer]

	c.doAdvanceProgramCounter()
}

// 0x1*** Instructions
func (c *CPU) do1NNN(operationCode uint16) {
	c.ProgramCounter = (operationCode & 0x0FFF)
}

// 0x2*** Instructions
func (c *CPU) do2NNN(operationCode uint16) {
	// Store the current location of the program counter to the stack
	// And increment the stack pointer
	// to insert a new entry to the stack with the last progarm counter posiiton
	c.Stack[c.StackPointer] = c.ProgramCounter
	c.StackPointer++
	// change the program counter to the current subroutine location
	c.ProgramCounter = operationCode & 0x0FFF
}

// 0x3*** Instructions
func (c *CPU) do3XNN(operationCode uint16) {
	// If value of register X == NN skip next instruction
	if c.Register[(operationCode&0x0F00)>>8] == (uint8(operationCode & 0x00FF)) {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0x4*** Instructions
func (c *CPU) do4XNN(operationCode uint16) {
	// If value of register X != NN skip next instruction
	if c.Register[(operationCode&0x0F00)>>8] != (uint8(operationCode & 0x00FF)) {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0x5*** Instructions
func (c *CPU) do5XY0(operationCode uint16) {
	// If value of register X == value of register Y skip next instruction
	if c.Register[(operationCode&0x0F00)>>8] == c.Register[(operationCode&0x00F0)>>4] {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0x6*** Instructions
func (c *CPU) do6XNN(operationCode uint16) {
	c.Register[(operationCode&0x0F00)>>8] = uint8(operationCode & 0x00FF)
	c.doAdvanceProgramCounter()
}

// 0x7*** Instructions
func (c *CPU) do7XNN(operationCode uint16) {
	c.Register[(operationCode&0x0F00)>>8] += uint8(operationCode & 0x00FF)
	c.doAdvanceProgramCounter()
}

// 0x8*** Instructions
func (c *CPU) do8XY0(operationCode uint16) {
	c.Register[(operationCode&0x0F00)>>8] = c.Register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY1(operationCode uint16) {
	c.Register[(operationCode&0x0F00)>>8] = c.Register[(operationCode&0x0F00)>>8] | c.Register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY2(operationCode uint16) {
	c.Register[(operationCode&0x0F00)>>8] = c.Register[(operationCode&0x0F00)>>8] & c.Register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY3(operationCode uint16) {
	c.Register[(operationCode&0x0F00)>>8] = c.Register[(operationCode&0x0F00)>>8] ^ c.Register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY4(operationCode uint16) {
	// Get the value of X and Y from the operationCode
	registerXLocation, registerYLocation := (operationCode&0x0F00)>>8, (operationCode&0x00F0)>>4
	// Get the value of X + Y
	result := c.Register[registerXLocation] + c.Register[registerYLocation]
	// If there is an overflow, set the overflow flag on register VF (0xF) to 1
	if c.Register[registerXLocation]+c.Register[registerYLocation] > 0xFF {
		c.Register[0xF] = 1
	} else {
		c.Register[0xF] = 0
	}
	// Set the result
	c.Register[registerXLocation] = result
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY5(operationCode uint16) {
	if c.Register[(operationCode&0x0F00)>>8] < c.Register[(operationCode&0x00F0)>>4] {
		c.Register[0xF] = 0
	} else {
		c.Register[0xF] = 1
	}
	c.Register[(operationCode&0x0F00)>>8] -= c.Register[(operationCode&0x00F0)>>4]
	c.doAdvanceProgramCounter()
}
func (c *CPU) do8XY6(operationCode uint16) {
	c.Register[0xF] = c.Register[(operationCode&0x0F00)>>8] & 0x1
	c.Register[(operationCode&0x0F00)>>8] >>= 1
	c.doAdvanceProgramCounter()
}

func (c *CPU) do8XY7(operationCode uint16) {
	if c.Register[(operationCode&0x0F00)>>8] > c.Register[(operationCode&0x00F0)>>4] {
		c.Register[0xF] = 0
	} else {
		c.Register[0xF] = 1
	}
	c.Register[(operationCode&0x0F00)>>8] = c.Register[(operationCode&0x00F0)>>4] - c.Register[(operationCode&0x0F00)>>8]
	c.doAdvanceProgramCounter()
}
func (c *CPU) do8XYE(operationCode uint16) {
	c.Register[0xF] = c.Register[(operationCode&0x0F00)>>8] >> 7
	c.Register[(operationCode&0x0F00)<<8] <<= 1
	c.doAdvanceProgramCounter()
}

// 0x9*** Instructions
func (c *CPU) do9XY0(operationCode uint16) {
	if c.Register[(operationCode&0x0F00)>>8] != c.Register[(operationCode&0x00F0)>>4] {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0xA*** Instructions
func (c *CPU) doANNN(operationCode uint16) {
	// Get last 3 * 4 bits of the operation code
	c.IndexRegister = operationCode & 0x0FFF
	c.doAdvanceProgramCounter()
}

// 0xB*** Instructions
func (c *CPU) doBNNN(operationCode uint16) {
	c.ProgramCounter = (operationCode & 0x0FFF) + uint16(c.Register[0])
}

// 0xC*** Instructions
func (c *CPU) doCXNN(operationCode uint16) {
	c.Register[(operationCode&0x0F00)>>8] = uint8(rand.Intn(0xFF)) & uint8(operationCode&0x00FF)
	c.doAdvanceProgramCounter()
}

// 0xD*** Instructions
func (c *CPU) doDXYN(operationCode uint16) {
	x, y, h := uint16(c.Register[(operationCode&0x0F00)>>8]), uint16(c.Register[(operationCode&0x00F0)>>4]), operationCode&0x000F
	pixelData := uint16(0)

	// Track if any pixels are flipped from set to unset.

	c.Register[0xF] = 0
	for yline := uint16(0); yline < h; yline++ {
		pixelData = uint16(c.Memory[c.IndexRegister+uint16(yline)])

		for xline := uint16(0); xline < 8; xline++ {

			index := (x + xline + ((y + yline) * 64))
			// Fix for offscreen rendering
			if index >= uint16(len(c.Screen)) {
				continue
			}
			if (pixelData & (0x80 >> xline)) != 0 {
				if c.Screen[index] == 1 {
					c.Register[0xF] = 1
				}
				c.Screen[index] ^= 1
			}
		}
	}

	c.ShouldDraw = true
	c.doAdvanceProgramCounter()
}

// 0xE*** Instructions
func (c *CPU) doEX9E(operationCode uint16) {
	if c.KeypadStates[c.Register[(operationCode&0x0F00)>>8]] != 0 {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}
func (c *CPU) doEXA1(operationCode uint16) {
	if c.KeypadStates[c.Register[(operationCode&0x0F00)>>8]] == 0 {
		c.doAdvanceProgramCounter()
	}
	c.doAdvanceProgramCounter()
}

// 0xF*** Instructions
func (c *CPU) doFX07(operationCode uint16) {
	c.Register[(operationCode&0x0F00)>>8] = c.DelayTimer
	c.doAdvanceProgramCounter()
}

func (c *CPU) doFX0A(operationCode uint16) {
	pressed := false
	for i, v := range c.KeypadStates {
		if v != 0 {
			c.Register[(operationCode&0x0F00)>>8] = uint8(i)
			pressed = true
		}
	}

	if !pressed {
		return
	}
	c.doAdvanceProgramCounter()
}

func (c *CPU) doFX15(operationCode uint16) {
	c.DelayTimer = c.Register[(operationCode&0x0F00)>>8]
	c.doAdvanceProgramCounter()
}
func (c *CPU) doFX18(operationCode uint16) {
	c.SoundTimer = c.Register[(operationCode&0x0F00)>>8]
	c.doAdvanceProgramCounter()
}
func (c *CPU) doFX1E(operationCode uint16) {
	if c.IndexRegister+uint16(c.Register[(operationCode&0x0F00)>>8]) > 0xFFF {
		c.Register[0xF] = 1
	} else {
		c.Register[0xF] = 0
	}
	c.IndexRegister += uint16(c.Register[(operationCode&0x0F00)>>8])
	c.doAdvanceProgramCounter()
}

func (c *CPU) doFX29(operationCode uint16) {
	c.IndexRegister = uint16(c.Register[(operationCode&0x0F00)>>8]) * 0x5
	c.doAdvanceProgramCounter()
}
func (c *CPU) doFX33(operationCode uint16) {
	// Get the value at register X
	registerXValue := c.Register[(operationCode&0x0F00)>>8]
	// Set the hundred's value of x to memory[I]
	c.Memory[c.IndexRegister] = registerXValue / 100
	// Set the ten's value of x to memory[I+1]
	c.Memory[c.IndexRegister+1] = (registerXValue / 10) % 10
	// Set the one's value of x to memory[I+2]
	c.Memory[c.IndexRegister+2] = (registerXValue % 100) % 10
	c.doAdvanceProgramCounter()
}

func (c *CPU) doFX55(operationCode uint16) {
	for i := 0; i <= int((operationCode&0x0F00)>>8); i++ {
		c.Memory[int(c.IndexRegister)+i] = c.Register[i]
	}

	// On the original system
	// When the operation is done
	// indexRegistered += X + 1

	c.IndexRegister += (operationCode&0x0F00)>>8 + 1
	c.doAdvanceProgramCounter()
}

func (c *CPU) doFX65(operationCode uint16) {
	for i := 0; i <= int((operationCode&0x0F00)>>8); i++ {
		c.Register[i] = c.Memory[int(c.IndexRegister)+i]
	}

	// On the original interpreter, when the operation is done, I = I + X + 1.
	c.IndexRegister += (operationCode&0x0F00)>>8 + 1
	c.doAdvanceProgramCounter()
}
