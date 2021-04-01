package cpu

// 0x0*** Instructions
// 0x1*** Instructions
// 0x2*** Instructions
// 0x3*** Instructions
// 0x4*** Instructions
// 0x5*** Instructions
// 0x6*** Instructions
// 0x7*** Instructions
// 0x8*** Instructions
// 0x9*** Instructions
// 0xA*** Instructions
func (c *CPU) doANNN(operationCode uint16) {
	c.indexRegister = operationCode & 0x0FFF
	c.programCounter += 2
}

// 0xB*** Instructions
// 0xC*** Instructions
// 0xD*** Instructions
// 0xE*** Instructions
// 0xF*** Instructions
