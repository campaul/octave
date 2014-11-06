package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Initializing Octave CPU...")
	cpu := &CPU{running: true, sp: 65535}

	file, err := os.Open(os.Args[1])

	if err != nil {
		return
	}

	defer file.Close()

	cpu.memory, err = ioutil.ReadAll(file)

	cpu.devices[0] = stack{cpu}
	cpu.devices[1] = tty{bufio.NewReader(os.Stdin)}

	if err != nil {
		return
	}

	for cpu.running {
		inst_byte := fetch(cpu)
		inst_func := decode(inst_byte)
		inst_func(inst_byte, cpu)
	}
}

type CPU struct {
	memory    []uint8
	registers [4]uint8
	pc        uint16
	sp        uint16
	running   bool
	result    uint8
	devices   [8]Device
}

type instruction func(uint8, *CPU)

type Device interface {
	read() uint8
	write(uint8)
}

type tty struct {
	reader *bufio.Reader
}

func (t tty) read() uint8 {
	b, _ := t.reader.ReadByte()
	return b
}

func (t tty) write(char uint8) {
	fmt.Printf("%c", char)
}

type stack struct {
	cpu *CPU
}

func (s stack) read() uint8 {
	value := s.cpu.memory[s.cpu.sp]
	s.cpu.sp++
	return value
}

func (s stack) write(char uint8) {
	s.cpu.memory[s.cpu.sp] = char
	s.cpu.sp--
}

func fetch(cpu *CPU) uint8 {
	inst := cpu.memory[cpu.pc]
	cpu.pc = cpu.pc + 1
	return inst
}

func decode(i uint8) instruction {
	inst := illegal

	switch i >> 5 {
	case 0:
		fmt.Fprint(os.Stderr, "jmp\n")
		inst = jmp
	case 1:
		fmt.Fprint(os.Stderr, "loadi\n")
		inst = loadi
	case 2:
		fmt.Fprint(os.Stderr, "math\n")
		inst = math
	case 3:
		fmt.Fprint(os.Stderr, "logic\n")
		inst = logic
	case 4:
		fmt.Fprint(os.Stderr, "mem\n")
		inst = mem
	case 5:
		fmt.Fprint(os.Stderr, "stack\n")
		inst = stacki
	case 6:
		fmt.Fprint(os.Stderr, "in\n")
		inst = in
	case 7:
		fmt.Fprint(os.Stderr, "out\n")
		inst = out
	}

	return inst
}

func jmp(i uint8, cpu *CPU) {
	if i == 0 {
		cpu.running = false
	}

	register := i << 3 >> 6
	n := i << 5 >> 7
	z := i << 6 >> 7
	p := i << 7 >> 7

	if (n==1 && cpu.result < 0) || (z==1 && cpu.result == 0) || (p==1 && cpu.result > 0) {
		offset := int8(cpu.registers[register])
		fmt.Fprintf(os.Stderr, "Taking jump to %v\n", offset)
		cpu.pc = uint16(int32(cpu.pc) + int32(offset))
	}
}

func loadi(i uint8, cpu *CPU) {
	location := i << 3 >> 7

	if location == 0 {
		cpu.registers[0] = (i << 4) | (cpu.registers[0] << 4 >> 4)
	} else {
		cpu.registers[0] = (i << 4 >> 4) | (cpu.registers[0] >> 4 << 4)
	}
}

func math(i uint8, cpu *CPU) {
	operation := i << 3 >> 7
	destination := i << 4 >> 6
	source := i << 6 >> 6

	if operation == 0 {
		cpu.registers[destination] = cpu.registers[destination] + cpu.registers[source]
	} else {
		cpu.registers[destination] = cpu.registers[destination] / cpu.registers[source]
	}

	cpu.result = cpu.registers[destination]
}

func logic(i uint8, cpu *CPU) {
	operation := i << 3 >> 7
	destination := i << 4 >> 6
	source := i << 6 >> 6

	if operation == 0 {
		cpu.registers[destination] = cpu.registers[destination] & cpu.registers[source]
	} else {
		cpu.registers[destination] = cpu.registers[destination] ^ cpu.registers[source]
	}

	cpu.result = cpu.registers[destination]
}

func mem(i uint8, cpu *CPU) {
	operation := i << 3 >> 7
	address_high := i << 4 >> 6
	address_low := i << 6 >> 6
	address := uint16(cpu.registers[address_high]) << 8 + uint16(cpu.registers[address_low])

	if operation == 0 {
		// LOAD
		fmt.Fprintf(os.Stderr, "Loading %v\n to R0", address)
		cpu.registers[0] = cpu.memory[address]
	} else {
		// STORE
		fmt.Fprintf(os.Stderr, "Storing R0 to %v\n", address)
		cpu.memory[address] = cpu.registers[0]
	}
}

func stacki(i uint8, cpu *CPU) {
	stacki := i << 3 >> 3

	switch stacki {
		case 0:
			// add16
		case 1:
			// sub16
		case 2:
			// mul16
		case 3:
			// div16
		case 4:
			// mod16
		case 5:
			// neg16
		case 6:
			// and16
		case 7:
			// or16
		case 8:
			// xor16
		case 9:
			// not16
		case 10:
			// add32
		case 11:
			// sub32
		case 12:
			// mul32
		case 13:
			// div32
		case 14:
			// mod32
		case 15:
			// neg32
		case 16:
			// and32
		case 17:
			// or32
		case 18:
			// xor32
		case 19:
			// not32
		case 20:
			// call
		case 21:
			// trap
		case 22:
			// ret
		case 23:
			// iret
		default:
			// device := stacki - 24
			// TODO: enable device
	}
}

func in(i uint8, cpu *CPU) {
	device := i << 5 >> 5
	destination := i << 3 >> 6
	cpu.registers[destination] = cpu.devices[device].read()
}

func out(i uint8, cpu *CPU) {
	device := i << 3 >> 5
	source := i << 6 >> 6
	cpu.devices[device].write(cpu.registers[source])
}

func illegal(i uint8, cpu *CPU) {
}
