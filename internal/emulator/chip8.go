package emulator

import (
	"fmt"
	"os"

	"github.com/g-sponda/chip8-gomulator/internal/utils"
)

const (
	col_size = 64
	row_size = 32
)

// Struct to hold Chip8 emulator state
type Chip8 struct {
	// 0x000 to 0x1FF will be reserved for the interpreter(program code, and fonts)
	// 0x200 to 0xFFF program space (where the program runs)
	memory [4096]byte // Memory 4Kb
	// Registers
	v           [16]byte   // this is an 16 8-bit(1 byte) data registers (V0 - VF)
	i           uint16     // this is a 16-bit index registers
	pc          uint16     // Program Counter, 16-bit. It points to the current instruction
	sp          uint16     // Points to the current position in the stack
	stack       [16]uint16 // 16 level deep stack, this is used to store return address
	delay_timer byte       // counts down to 0
	sound_timer byte       // counts down to 0
	// IOs
	keys   [16]byte                 // 16 keys (0-9, A-F), used for keypad state (0 = not pressed; 1 = pressed)
	screen [row_size][col_size]byte // Screen 64x32 grid
}

func (c *Chip8) LoadRom(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error. Could not read file: %s %w", filename, err)
	}

	for i := 0; i < len(data); i++ {
		// Program space starts at 0x200 address
		c.memory[0x200+i] = data[i]
	}

	return nil
}

// --- Opcode functions --- //
// opcode (ak. Operation code). Part of machine instruction to tell the processor what action to perform
// It's a binary or hexadecimal represatation of an specific instruction

// Fetches and decodes instructions to opcode
//
// each instruction in Chip-8 is 2 bytes. Memory is byte-addressable (1 byte)
// we combine 2 consecutive memory locations to make the opcode
func (c *Chip8) Fetch() uint16 {
	return uint16(c.memory[c.pc]<<8 | c.memory[c.pc+1])
}

// Set all screen pixels back to 0
func (c *Chip8) ClearScreen() {
	for row := range c.screen {
		for col := range c.screen[row] {
			c.screen[row][col] = 0
		}
	}
}

// Move stack pointer 1 step down and set program counter to the current instruction (new stack pointer)
func (c *Chip8) ReturnFromSubroutine() {
	c.sp--               // move stack pointer 1 step down
	c.pc = c.stack[c.sp] // assign to program counter the current instruction in the stack
}

// Jump to address
// receives opcode
func (c *Chip8) JumpToAddr(opcode uint16) {
	// Use bitwise AND operation to get the address that we want the program counter to jump
	addr := opcode & 0x0FFF
	c.pc = addr // set program counter to the receive address
}

// Call Subroutine
func (c *Chip8) CallSubroutine() {
	c.sp++               // move stackj pointer 1 step up
	c.pc = c.stack[c.sp] // assignto program counter the current instruction in the stack
}

// Load value into register
// receives opcode
func (c *Chip8) LoadValueIntoRegister(opcode uint16) {
	x := (opcode & 0x0F00) >> 8 // Get X register index value
	nn := byte(opcode & 0x00FF) // Get NN value of opcode
	c.v[x] = nn
}

// Add to register
// receives opcode
func (c *Chip8) AddtoRegister(opcode uint16) {
	x := (opcode & 0x0F00) >> 8 // Get X register index value
	nn := byte(opcode & 0x00FF) // Get NN value of opcode
	c.v[x] += nn
}

// Set index to address
// receives opcode
func (c *Chip8) SetIndexToAddr(opcode uint16) {
	c.i = (opcode & 0x0FFF) // 0xANNN - Set I = NNN
}

// Draw sprite at (VX, VY)
// receives opcode
func (c *Chip8) DrawSprite(opcode uint16) {
	// draw 8xN pixel sprite at position vX, vY with data starting at the address in I
	// To get the X and Y coordinations and make sure it's inside the limit,
	// we can get the value mod the collumn/row size, so we ensure it's in the limit
	// if the value is smaller than the divisor, the remainder is the value itself.
	// This way of implement guarantee that we don't draw outside of the screen, but this will wrap around
	x_pos := c.v[(opcode&0x0F00)>>8] % col_size // get X coordination from VX register
	y_pos := c.v[(opcode&0x00F0)>>4] % row_size // get Y coordination from VY register
	height_n := (opcode & 0x000F)               // get height, the width is always 8 pixels wide
	c.v[0xF] = 0

	// implement loop to set DrawSprite and set pixels ON or OFF
	// Also need to set VF register if a pisxel is erased (collision detection)
	// we need to draw 8xN pixel into the sprite where n is the height.
	// let's loop to our screen, starting from poisition
	for row := uint16(0); row < height_n; row++ {
		sprite_byte := c.memory[c.i+row] // gets one row of the sprite, that needs to be draw. this is a 8-bit value.
		for col := 0; col < 8; col++ {
			// 0x80 is 8 bits, where only the left most is 1. >> col we shift this bit col positions
			// since we only want to draw or "activate" the pixel, we can do an bitwise AND operation,
			// to compare if the bit of the sprite_byte and the current collumn bit result in 1, which means we should
			// turn the pixel on, or validate collision
			if (sprite_byte & (0x80 >> col)) != 0 {
				if c.screen[row+uint16(y_pos)][col+int(x_pos)] == 1 { // A collision happened
					c.v[0xF] = 1 // VF is the data register used to check collision, we set to 1, since a collision happened
				}
				c.screen[row+uint16(y_pos)][col+int(x_pos)] ^= 1 // bitwise XOR operation, only results in true when values differ
				// 0 XOR 0 = 0
				// 1 XOR 0 = 1
				// 0 XOR 1 = 1
				// 1 XOR 1 = 0
			}
		}
	}
}

func (c *Chip8) ExecuteOpcode() {
	opcode := c.Fetch()
	c.pc += 2 // move 2 bytes to go to the next instruction

	/* We will do a bitwise AND (&) operation. For that we will be applying a mask to extract the first nibble,
	 * which determines the type of the instruction
	 * e.G   1010 0010 1111 0000   (opcode = 0xA2F0)
	 *		 & 1111 0000 0000 0000   (mask = 0xF000)
	 *			 --------------------
	 *			 1010 0000 0000 0000   (result = 0xA000)
	 */
	switch opcode & 0xF000 {
	// System Category cases - 0x0
	case 0x0000: // Handle 0x0000 opcodes - e.G clear screen, return from a subroutine
		switch opcode {
		case 0x00E0: // Opcode to Clear the screen
			c.ClearScreen()
		case 0x00EE: // opcode to returns from the subroutine
			c.ReturnFromSubroutine()
		}
	// Flow Control category cases - 0x1 - 0x2
	case 0x1000: // Handle 0x1NNN - Jump to address NNN
		c.JumpToAddr(opcode) // bitwise AND operation, get the values.
	case 0x2000: // Handle 0x2NNN - Call subroutine
		c.CallSubroutine()
	// Registers Category cases - 0x6, 0x7 and 0xA
	case 0x6000: // Handle 0x6NNN - Load value into register Vx. Set register VX = NN.
		c.LoadValueIntoRegister(opcode)
	case 0x7000: // Handle 0x7NNN - Add NN to VX
		c.AddtoRegister(opcode)
	case 0xA000: // Handle 0xANNN - Set I = NNN
		c.SetIndexToAddr(opcode)
	// Graphics category cases - 0xD
	case 0xD000: // Handle 0xDXYN - Draw Sprite at (VX,VY)
		c.DrawSprite(opcode)
	// Input category cases - 0xE
	case 0xE00: // Handle 0xENN - Skip next instruction when key not pressed
	// Conditionals category cases - 0x3 - 0x4
	case 0x300: // Handle 0x3NNN - Skip if VX == NN
		utils.PrintOpcodeNeedsImplementation(opcode)
	case 0x400: // Handle 0x4NNN - Skip if VX != NN
		utils.PrintOpcodeNeedsImplementation(opcode)
	// Timers category cases - 0xF
	case 0xF00:
		switch opcode {
		case 0xF07: // Handle 0xF07 - Read delay timer
			utils.PrintOpcodeNeedsImplementation(opcode)
		case 0xF15: // Handle 0xF15 - Set delay timer
			utils.PrintOpcodeNeedsImplementation(opcode)
		}
	default:
		fmt.Printf("Unkown opcode 0x%X\n", opcode)
		// Add cases for other opcodes here
	}
}
