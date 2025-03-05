package emulator

import (
	"testing"
)

func TestDrawSprite(t *testing.T) {
	// Step 1: Initialize Emulator
	chip8 := Chip8{}
	chip8.ClearScreen() // Ensure screen starts empty

	// Step 2: Load a sprite into memory (a 3-row sprite)
	chip8.memory[0x300] = 0xF0 // 11110000
	chip8.memory[0x301] = 0x90 // 10010000
	chip8.memory[0x302] = 0xF0 // 11110000
	chip8.i = 0x300            // Set index register to sprite location

	// Step 3: Set Registers (VX = 5, VY = 3)
	chip8.v[0] = 5           // X position
	chip8.v[1] = 3           // Y position
	opcode := uint16(0xD013) // Draw 3-row sprite at (5,3)

	// Step 4: Call DrawSprite
	chip8.DrawSprite(opcode)

	// Step 5: Validate pixels (Check expected coordinates)
	expected := [][]int{
		{5, 3}, {6, 3}, {7, 3}, {8, 3}, // 1111 at row 3
		{5, 4}, {8, 4}, // 1  1 at row 4
		{5, 5}, {6, 5}, {7, 5}, {8, 5}, // 1111 at row 5
	}

	for _, pixel := range expected {
		x, y := pixel[0], pixel[1]
		if chip8.screen[y][x] != 1 {
			t.Errorf("Expected pixel (%d, %d) to be ON", x, y)
		}
	}

	// Step 6: Draw the same sprite again (XOR should turn pixels OFF)
	chip8.DrawSprite(opcode)

	for _, pixel := range expected {
		x, y := pixel[0], pixel[1]
		if chip8.screen[y][x] != 0 {
			t.Errorf("Expected pixel (%d, %d) to be OFF after second draw", x, y)
		}
	}

	// Step 7: Check if VF flag is set (collision detection)
	if chip8.v[0xF] != 1 {
		t.Errorf("Expected VF flag to be set (1) on collision, but got %d", chip8.v[0xF])
	}
}
