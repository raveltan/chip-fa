package cpu

import (
	"fmt"
	"io/ioutil"
)

type CPU struct {
	// Most common implementation of Chip8 uses 4k of Memory
	// https://en.wikipedia.org/wiki/CHIP-8#Memory
	Memory [4096]uint8

	// Chip8's Register is a 1 byte general purpose Register V0,V1...VE.
	Register [16]uint8

	// Index Register
	IndexRegister uint16

	// Program Counter
	ProgramCounter uint16

	// Chip8's screen is 64 x 32 black and white screen.
	Screen [64 * 32]uint8

	// Timers (60Hz)
	DelayTimer uint8
	// Will buzz the system when it reaches 0
	// Updated on the rendering side
	SoundTimer uint8

	// Chip8's Stack
	Stack [16]uint16
	// Stack pointer
	StackPointer uint16

	// Keypad states tracker
	KeypadStates [16]uint8

	// Draw Flag
	// Should be exported to be used on the Emulator struct
	// FIXME: currently not used due the internal working of the graphic libarry
	// that always draws 60Hz thus this flag is ignored.
	ShouldDraw bool

	StopForDebuggingCallback func()
}

func (c *CPU) Boot() {
	// Chip8's application entry point is 0x200
	// Program counter should start at application entry point
	// Adreesses before 0x200 is commonly used by the interpreter
	c.ProgramCounter = 0x200

	// Load fontset to memory
	for i, v := range fontset {
		c.Memory[i] = v
	}
}

func (c *CPU) LoadROM(file string) error {
	rom, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if len(rom) > maxRomSize {
		panic("error: ROM size is too big for this system. Max size: 3232 (0xEA0 - 0x200)")
	}

	for i := 0; i < len(rom); i++ {
		c.Memory[0x200+i] = rom[i]
	}

	return nil
}

func (c *CPU) DoCycle() {
	// Fetch operationCode
	// gets opCode on the memory address specified by the programCounter
	// memory is one byte while program counter is 2 bytes, thus we need to fetch [addr] and [addr+1]

	// Combine 2 1 bytes memory address to a single 2 bytes operation code with bitwise and or operation
	// -----------
	// Example
	// c.memory[c.programCounter] = 0xFF = 0b11111111
	// c.memory[c.programCounter + 1] = 0x10 = 0b00010000
	// ->
	// 0xFF << 8 = 0b1111111100000000 = 0xFF00
	// 0xFF00 | 0x10 = 0xFF10
	// Resulting 0xFF10 as the operationCode
	// -----------

	currentOperationCode := uint16(c.Memory[c.ProgramCounter])<<8 | uint16(c.Memory[c.ProgramCounter+1])

	// Decode operationCode
	// operationCode table: https://en.wikipedia.org/wiki/CHIP-8#Opcode_table

	// Check for the first 4 bit of the operationCode by using bitwise
	// and branch out for each.
	// Example
	// 0x8123
	// Get operationCode by doing operationCode & 0xF000
	// 0x8123 & 0xF000
	// => 0x8000
	// Which is the first 4 bit of the operationCode
	switch currentOperationCode & 0xF000 {
	case 0x0000:
		switch currentOperationCode & 0x000F {
		// 0NNN: Call machine code routine (not used in most ROMS). (NOT IMPLEMENTED)
		case 0x0000:
			// 00E0: Clear Screen.
			c.do00E0()
		case 0x0001:
			// 0001: PROGRAM BREAKPOINT
			// TODO: Add to documentation
			c.do0001()
		case 0x000E:
			// 00EE: Retrun from subroutine.
			c.do000E()
		default:
			panic(fmt.Sprintf("error: Unknown operationCode (0x%x)", currentOperationCode))
		}
	case 0x1000:
		// 1NNN: Jumps to address NNN.
		c.do1NNN(currentOperationCode)
	case 0x2000:
		// 2NNN: Call subroutine at NNN.
		c.do2NNN(currentOperationCode)
	case 0x3000:
		// 3XNN: Skips the next instruction if VX equals NN.
		// (Usually the next instruction is a jump to skip a code block)
		c.do3XNN(currentOperationCode)
	case 0x4000:
		// 4XNN: Skips the next instruction if VX doesn't equal NN.
		// (Usually the next instruction is a jump to skip a code block)
		c.do4XNN(currentOperationCode)
	case 0x5000:
		// 5XY0: Skips the next instruction if VX equals VY.
		// (Usually the next instruction is a jump to skip a code block)
		c.do5XY0(currentOperationCode)
	case 0x6000:
		// 6XNN: Sets VX to NN.
		c.do6XNN(currentOperationCode)
	case 0x7000:
		// 7XNN: Adds NN to VX. (Carry flag is not changed)
		c.do7XNN(currentOperationCode)
	case 0x8000:
		switch currentOperationCode & 0x000F {
		case 0x0000:
			// 8XY0: Sets VX to the value of VY.
			c.do8XY0(currentOperationCode)
		case 0x0001:
			// 8XY1: Sets VX to VX or VY. (Bitwise OR operation)
			c.do8XY1(currentOperationCode)
		case 0x0002:
			// 8XY2: Sets VX to VX and VY. (Bitwise AND operation)
			c.do8XY2(currentOperationCode)
		case 0x0003:
			// 8XY3: Sets VX to VX xor VY.
			c.do8XY3(currentOperationCode)
		case 0x0004:
			// 8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
			c.do8XY4(currentOperationCode)
		case 0x0005:
			// 8XY5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
			c.do8XY5(currentOperationCode)
		case 0x0006:
			// 8XY6: Stores the least significant bit of VX in VF and then shifts VX to the right by 1.
			c.do8XY6(currentOperationCode)
		case 0x0007:
			// 8XY7: Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
			c.do8XY7(currentOperationCode)
		case 0x000E:
			// 8XYE: Stores the most significant bit of VX in VF and then shifts VX to the left by 1.
			c.do8XYE(currentOperationCode)
		default:
			panic(fmt.Sprintf("error: Unknown operationCode (0x%x)", currentOperationCode))
		}

	case 0x9000:
		// 9XY0: Skips the next instruction if VX doesn't equal VY.
		// (Usually the next instruction is a jump to skip a code block)
		c.do9XY0(currentOperationCode)
	case 0xA000:
		// ANNN: Sets I to the address NNN.
		c.doANNN(currentOperationCode)

	case 0xB000:
		// BNNN: Jumps to the address NNN plus V0.
		c.doBNNN(currentOperationCode)
	case 0xC000:
		// CXNN: Sets VX to the result of a bitwise and operation on a random number.
		// (Typically: 0 to 255) and NN.
		c.doCXNN(currentOperationCode)
	case 0xD000:
		// DXYN: Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N+1 pixels.
		// Each row of 8 pixels is read as bit-coded starting from memory location I;
		// I value doesn???t change after the execution of this instruction. As described above,
		// VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn,
		// and to 0 if that doesn???t happen
		c.doDXYN(currentOperationCode)
	case 0xE000:
		switch currentOperationCode & 0x000F {
		case 0x000E:
			// EX9E: Skips the next instruction if the key stored in VX is pressed.
			// (Usually the next instruction is a jump to skip a code block)
			c.doEX9E(currentOperationCode)
		case 0x0001:
			// EXA1: Skips the next instruction if the key stored in VX isn't pressed.
			// (Usually the next instruction is a jump to skip a code block)
			c.doEXA1(currentOperationCode)
		default:
			panic(fmt.Sprintf("error: Unknown operationCode (0x%x)", currentOperationCode))
		}
	case 0xF000:
		switch currentOperationCode & 0x000F {
		case 0x0007:
			// FX07: Sets VX to the value of the delay timer.
			c.doFX07(currentOperationCode)
		case 0x000A:
			// FX0A: A key press is awaited, and then stored in VX.
			// (Blocking Operation. All instruction halted until next key event)
			c.doFX0A(currentOperationCode)
		case 0x0005:
			switch currentOperationCode & 0x00F0 {
			case 0x0010:
				// FX15: Sets the delay timer to VX.
				c.doFX15(currentOperationCode)
			case 0x0050:
				// FX55: Stores V0 to VX (including VX) in memory starting at address I.
				// The offset from I is increased by 1 for each value written, but I itself is left unmodified.
				c.doFX55(currentOperationCode)
			case 0x0060:
				// FX65: Fills V0 to VX (including VX) with values from memory starting at address I.
				// The offset from I is increased by 1 for each value written, but I itself is left unmodified.
				c.doFX65(currentOperationCode)
			default:
				panic(fmt.Sprintf("error: Unknown operationCode (0x%x)", currentOperationCode))
			}
		case 0x0008:
			// FX18: Sets the sound timer to VX.
			c.doFX18(currentOperationCode)
		case 0x000E:
			// FX1E: Adds VX to I. VF is not affected.
			c.doFX1E(currentOperationCode)
		case 0x0009:
			// FX29: Sets I to the location of the sprite for the character in VX.
			//  Characters 0-F (in hexadecimal) are represented by a 4x5 font.
			c.doFX29(currentOperationCode)
		case 0x0003:
			// FX33: Stores the binary-coded decimal representation of VX,
			// with the most significant of three digits at the address in I,
			// the middle digit at I plus 1, and the least significant digit at I plus 2.
			// (In other words, take the decimal representation of VX,
			// place the hundreds digit in memory at location in I, the tens digit at location I+1,
			// and the ones digit at location I+2.)
			c.doFX33(currentOperationCode)
		default:
			panic(fmt.Sprintf("error: Unknown operationCode (0x%x)", currentOperationCode))

		}
	default:
		panic(fmt.Sprintf("error: Unknown operationCode (0x%x)", currentOperationCode))
	}

	// Sound timer is updated on the rendering side

	if c.DelayTimer > 0 {
		c.DelayTimer--
	}
}
