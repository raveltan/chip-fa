package cpu

import (
	"io/ioutil"
	"log"
)

type CPU struct {
	// Chip8's opCode is 2 bytes long (16 bits)
	operationCode uint16

	// Most common implementation of Chip8 uses 4k of memory
	// https://en.wikipedia.org/wiki/CHIP-8#Memory
	memory [4096]uint8

	// Chip8's register is a 1 byte general purpose register V0,V1...VE.
	register [16]uint8

	// Index Register
	indexRegister uint16

	// Program Counter
	programCounter uint16

	// Chip8's screen is 64 x 32 black and white screen.
	screen [64 * 32]uint8

	// Timers (60Hz)
	delayTimer uint8
	// Will buzz the system when it reaches 0
	soundTimer uint8

	// Chip8's Stack
	stack [16]uint16
	// Stack pointer
	stackPointer uint16

	// Keypad states tracker
	keypadStates [16]uint8

	// Draw Flag
	// Should be exported to be used on the Emulator struct
	ShouldDraw bool
}

func (c *CPU) Boot() {
	// Chip8's application entry point is 0x200
	// Program counter should start at application entry point
	c.programCounter = 0x200

	// Load fontset to memory
	for i, v := range fontset {
		c.memory[i] = v
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
		c.memory[0x200+i] = rom[i]
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
	currentOperationCode := uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])
	log.Println(currentOperationCode)

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
	case 0x000:
		// TODO: 0NNN: Call machine code routine (not used in most ROMS).

		// TODO: 00E0: Clear Screen.

		// TODO: 00EE: Retrun from subroutine.
		break
	case 0x100:
		// TODO: 1NNN: Jumps to address NNN.
		break
	case 0x200:
		// TODO: 2NNN: Call subroutine at NNN.
		break
	case 0x300:
		// TODO: 3XNN: Skips the next instruction if VX equals NN.
		// (Usually the next instruction is a jump to skip a code block)
		break
	case 0x400:
		// TODO: 4XNN: Skips the next instruction if VX doesn't equal NN.
		// (Usually the next instruction is a jump to skip a code block)
		break
	case 0x500:
		// TODO: 5XY0: Skips the next instruction if VX equals VY.
		// (Usually the next instruction is a jump to skip a code block)
		break
	case 0x600:
		// TODO: 6XNN: Sets VX to NN.
		break
	case 0x700:
		// TODO: 7XNN: Adds NN to VX. (Carry flag is not changed)
		break
	case 0x800:
		// TODO: 8XY0: Sets VX to the value of VY.

		// TODO: 8XY1: Sets VX to VX or VY. (Bitwise OR operation)

		// TODO: 8XY2: Sets VX to VX and VY. (Bitwise AND operation)

		// TODO: 8XY3: Sets VX to VX xor VY.

		// TODO: 8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.

		// TODO: 8XY5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.

		// TODO: 8XY6: Stores the least significant bit of VX in VF and then shifts VX to the right by 1.

		// TODO: 8XY7: Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.

		// TODO: 8XYE: Stores the most significant bit of VX in VF and then shifts VX to the left by 1.
		break
	case 0x900:
		// TODO: 9XY0: Skips the next instruction if VX doesn't equal VY.
		// (Usually the next instruction is a jump to skip a code block)
		break
	case 0xA00:
		// ANNN: Sets I to the address NNN.
		c.doANNN(currentOperationCode)

	case 0xB00:
		// TODO: BNNN: Jumps to the address NNN plus V0.
		break
	case 0xC00:
		// TODO: CXNN: Sets VX to the result of a bitwise and operation on a random number.
		// (Typically: 0 to 255) and NN.
		break
	case 0xD00:
		// TODO: DXYN: Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N+1 pixels.
		// Each row of 8 pixels is read as bit-coded starting from memory location I;
		// I value doesn’t change after the execution of this instruction. As described above,
		// VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn,
		// and to 0 if that doesn’t happen
		break
	case 0xE00:
		// TODO: EX9E: Skips the next instruction if the key stored in VX is pressed.
		// (Usually the next instruction is a jump to skip a code block)

		// TODO: EXA1: Skips the next instruction if the key stored in VX isn't pressed.
		// (Usually the next instruction is a jump to skip a code block)
		break
	case 0xF00:
		// TODO: FX07: Sets VX to the value of the delay timer.

		// TODO: FX0A: A key press is awaited, and then stored in VX.
		// (Blocking Operation. All instruction halted until next key event)

		// TODO: FX15: Sets the delay timer to VX.

		// TODO: FX18: Sets the sound timer to VX.

		// TODO: FX1E: Adds VX to I. VF is not affected.

		// TODO: FX29: Sets I to the location of the sprite for the character in VX.
		//  Characters 0-F (in hexadecimal) are represented by a 4x5 font.

		// TODO: FX33: Stores the binary-coded decimal representation of VX,
		// with the most significant of three digits at the address in I,
		// the middle digit at I plus 1, and the least significant digit at I plus 2.
		// (In other words, take the decimal representation of VX,
		// place the hundreds digit in memory at location in I, the tens digit at location I+1,
		// and the ones digit at location I+2.)

		// TODO: FX55: Stores V0 to VX (including VX) in memory starting at address I.
		// The offset from I is increased by 1 for each value written, but I itself is left unmodified.

		// TODO: FX65: Fills V0 to VX (including VX) with values from memory starting at address I.
		// The offset from I is increased by 1 for each value written, but I itself is left unmodified.
	}

	// Execute operationCode

	// Update Timer
	if c.soundTimer > 0 {
		if c.soundTimer == 1 {
			// Beep the buzzer
			log.Println("Beep")
		}
		c.soundTimer--
	}

	if c.delayTimer > 0 {
		c.delayTimer--
	}

}
