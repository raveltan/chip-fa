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
	var currentOperationCode uint16
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
	currentOperationCode = uint16(c.memory[c.programCounter]<<8) | uint16(c.memory[c.programCounter+1])
	log.Println(currentOperationCode)
	// Decode operationCode
	// operationCode table: https://en.wikipedia.org/wiki/CHIP-8#Opcode_table

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
