package main

import (
	"fmt"
)

func (c8 *Chip8) execInstr() {
	nextInstr := c8.programPtr + 2
	instr := c8.memory[c8.programPtr:nextInstr]

	switch {
	// 00E0 - CLS - clear the display
	case instr[0] == 0x00 && instr[1] == 0xE0:
		c8.display.clear()
	// 00EE RET returns from a subroutine
	case instr[0] == 0x00 && instr[1] == 0xEE:
		cs := c8.callStack
		nextInstr, c8.callStack = cs[len(cs)-1], cs[:len(cs)-1]
	// 1nnn - JP addr
	case lNib(instr[0]) == 0x1:
		addr := getAddr(instr)
		nextInstr = addr
		// 2nnn - JP addr - pushes program counter +2 to the call stack and makes program counter = nnn
	case lNib(instr[0]) == 0x2:
		c8.callStack = append(c8.callStack, c8.programPtr+2)
		addr := getAddr(instr)
		nextInstr = addr
	// 3xkk - SE Vx, byte - Skip next instruction if Vx = kk
	case lNib(instr[0]) == 0x3:
		if c8.registers[rNib(instr[0])] == instr[1] {
			nextInstr += 2
		}
	// 4xkk SNE Vx, byte - Skip next instruction if Vx != kk
	case lNib(instr[0]) == 0x4:
		if c8.registers[rNib(instr[0])] != instr[1] {
			nextInstr += 2
		}
	// 6xkk - LD Vx, byte - Load the byte value into the register specified by x
	case lNib(instr[0]) == 0x6:
		c8.registers[rNib(instr[0])] = instr[1]
	// 7xkk - ADD Vx, byte
	case lNib(instr[0]) == 0x7:
		c8.registers[rNib(instr[0])] += instr[1]
	// Annn - LD I, addr - Load the int16 addr specified by nnn into the I register
	case lNib(instr[0]) == 0xA:
		addr := getAddr(instr)
		c8.regI = addr
	// Dxyn - DRW Vx, Vy, nibble - grab an nibble length byte from I and draw it at the
	// values of Vx and Vy. If at least one pixel is erased set VF to 1 otherwise to 0
	// if a part of the sprite is located off screen - wrap it.
	case lNib(instr[0]) == 0xD:
		xOffset := 56 - int(c8.registers[rNib(instr[0])])
		yOffset := c8.registers[lNib(instr[1])]
		length := rNib(instr[1])

		sprite := c8.memory[c8.regI : c8.regI+uint16(length)]
		//for _, line := range sprite {
		//	fmt.Printf("%08b\n", line)
		//}

		c8.registers[0xF] = 0
		for i := 0; i < int(length); i++ {
			var spriteRow uint64
			// if we need to wrap
			if xOffset < 0 {
				unWrappedOffset := xOffset * -1
				unWrapped := uint64(sprite[i]) >> (uint(unWrappedOffset))
				wrapped := uint64(sprite[i]) << uint(64+xOffset)
				spriteRow = unWrapped ^ wrapped

			} else {
				spriteRow = uint64(sprite[i]) << uint(xOffset)
			}
			y := (yOffset + byte(i)) % 32
			currentRow := c8.display.screen[y]

			collisionFree := currentRow | spriteRow
			c8.display.screen[y] = currentRow ^ spriteRow
			if collisionFree^c8.display.screen[y] > 0 && c8.registers[0xF] == 0 {
				c8.registers[0xF] = 1
			}
		}
	// Fx1E = ADD I, VX - Add Vx to I and store in I
	case lNib(instr[0]) == 0xF && instr[1] == 0x1E:
		c8.regI += uint16(c8.registers[rNib(instr[0])])
	default:
		msg := fmt.Sprintf("Unknown Instruction: %X\n", instr)
		panic(msg)
	}

	c8.programPtr = nextInstr
}

func lNib(b byte) byte {
	return b >> 4
}

func rNib(b byte) byte {
	return b & 15
}

func getAddr(bs []byte) uint16 {
	right := uint16(rNib(bs[0])) << 8
	left := uint16(bs[1])
	return right + left
}
