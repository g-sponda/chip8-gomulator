# Chip 8 Gomulator

## Introduction

This is a simple CHIP-8 emulator written in Go.

CHIP-8 is an interpreted programming language developed in the 1970s to run on 8-bit microcomputers. It provides a simple architecture designed for running games.

### CHIP-8 Architecture

The CHIP-8 system consists of the following components:

- **Memory**: 4 KB of RAM, where the first 512 bytes are reserved for the interpreter.
- **Registers**: 16 general-purpose 8-bit registers (V0 - VF), with VF used for flags.
- **Index Register (I)**: A 16-bit register used for memory operations.
- **Program Counter (PC)**: Points to the current instruction in memory.
- **Stack and Stack Pointer (SP)**: Used for handling subroutine calls.
- **Timers**: Delay and sound timers, which decrement at 60 Hz.
- **Input**: A hexadecimal keypad (0-F).
- **Graphics**: A 64x32 monochrome display.

CHIP-8 programs consist of 35 opcodes that control the system's memory, registers, display, and input.

For more details on the CHIP-8 architecture, refer to:

- [Wikipedia: CHIP-8](https://en.wikipedia.org/wiki/CHIP-8)
- [Cowgod's Chip-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)

## Implemented Opcodes

Currently, the emulator supports the following CHIP-8 opcodes:

- **System Category (0x0)**

  - `0x00E0`: Clear the screen.
  - `0x00EE`: Return from a subroutine.

- **Flow Control Category (0x1 - 0x2)**

  - `0x1NNN`: Jump to address NNN.
  - `0x2NNN`: Call subroutine at address NNN.

- **Registers Category (0x6, 0x7, 0xA)**

  - `0x6XNN`: Load value NN into register VX.
  - `0x7XNN`: Add NN to register VX.
  - `0xANNN`: Set index register I to address NNN.

- **Graphics Category (0xD)**

  - `0xDXYN`: Draw sprite at coordinates (VX, VY).

- **Input Category (0xE)**
  - `0xEX9E`: Skip next instruction if key VX is pressed.
  - `0xEXA1`: Skip next instruction if key VX is not pressed.

## Features

- Full CHIP-8 CPU emulation
- Basic input handling
- Display rendering
- ROM loading support

## Requirements

- Go (>=1.18)
- SDL2 (for graphics and input handling)

## Installation

Clone the repository:

```sh
git clone https://github.com/g-sponda/chip8-gomulator.git
cd chip8-emulator
```

## Running the Emulator

To build and run the emulator:

```sh
go run main.go <path_to_rom>
```

Or compile it into a binary:

```sh
go build -o chip8 .
./chip8 <path_to_rom>
```

<!--

## Controls

The CHIP-8 uses a hex-based keypad mapped to standard keyboard keys:

```
1 2 3 4      ->  1 2 3 C
Q W E R      ->  4 5 6 D
A S D F      ->  7 8 9 E
Z X C V      ->  A 0 B F
```

## TODO

- Sound support
- Super CHIP-8 compatibility
- Debugger and disassembler
-->

## License

This project is licensed under the GNU GPLv3 License. See `LICENSE` for details.

## Acknowledgments

- [Wikipedia: CHIP-8](https://en.wikipedia.org/wiki/CHIP-8)
- [Cowgod's CHIP-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
