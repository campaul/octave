package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Initializing Octave CPU...")
	cpu := &CPU{running: true, sp: 0xFFFF}

	file, err := os.Open(os.Args[1])

	if err != nil {
		return
	}

	defer file.Close()

	mem, err := ioutil.ReadAll(file)

	for i, val := range mem {
		cpu.memory[i] = val
	}

	cpu.stack = stack{cpu}
	cpu.devices[0] = cpu.stack
	cpu.devices[1] = tty{bufio.NewReader(os.Stdin)}
	cpu.devices[4] = register{}
	cpu.devices[5] = register{}
	cpu.devices[6] = register{}
	cpu.devices[7] = register{}

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
	memory    [1 << 16]uint8
	registers [4]uint8
	pc        uint16
	sp        uint16
	running   bool
	result    uint8
	devices   [8]Device
	stack     stack
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
	s.cpu.sp++
	value := s.cpu.memory[s.cpu.sp]
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

type register struct {
	value uint8
}

func (r register) read() uint8 {
	return r.value
}

func (r register) write(char uint8) {
	r.value = char
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

	if (n == 1 && cpu.result < 0) || (z == 1 && cpu.result == 0) || (p == 1 && cpu.result > 0) {
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
		cpu.registers[destination] = cpu.registers[source]
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
	address := uint16(cpu.registers[address_high])<<8 + uint16(cpu.registers[address_low])

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

func pop16(s stack) uint16 {
	byte_1 := s.read()
	byte_2 := s.read()
	return uint16(byte_1)<<8 + uint16(byte_2)
}

func pop32(s stack) uint32 {
	byte_1 := s.read()
	byte_2 := s.read()
	byte_3 := s.read()
	byte_4 := s.read()
	return uint32(byte_1)<<24 + uint32(byte_2)<<16 + uint32(byte_3)<<8 + uint32(byte_4)
}

func push16(s stack, value uint16) {
	byte_1 := uint8(value >> 8)
	byte_2 := uint8(value << 8 >> 8)
	s.write(byte_2)
	s.write(byte_1)
}

func stacki(i uint8, cpu *CPU) {
	stacki := i << 3 >> 3

	switch stacki {
	case 0:
		// add16
		b := cpu.stack.read()
		a := cpu.stack.read()
		cpu.stack.write(a + b)
	case 1:
		// sub16
		b := cpu.stack.read()
		a := cpu.stack.read()
		cpu.stack.write(a - b)
	case 2:
		// mul16
		b := cpu.stack.read()
		a := cpu.stack.read()
		cpu.stack.write(a * b)
	case 3:
		// div16
		b := cpu.stack.read()
		a := cpu.stack.read()
		cpu.stack.write(a / b)
	case 4:
		// mod16
		b := cpu.stack.read()
		a := cpu.stack.read()
		cpu.stack.write(a % b)
	case 5:
		// neg16
		a := cpu.stack.read()
		cpu.stack.write(uint8(int8(a) * -1))
	case 6:
		// and16
		b := cpu.stack.read()
		a := cpu.stack.read()
		cpu.stack.write(a & b)
	case 7:
		// or16
		b := cpu.stack.read()
		a := cpu.stack.read()
		cpu.stack.write(a | b)
	case 8:
		// xor16
		b := cpu.stack.read()
		a := cpu.stack.read()
		cpu.stack.write(a ^ b)
	case 9:
		// not16
		a := cpu.stack.read()
		cpu.stack.write(^a)
	case 10:
	case 11:
	case 12:
	case 13:
	case 14:
	case 15:
	case 16:
	case 17:
	case 18:
	case 19:
	case 20:
		// Get jump address off the stack
		new_pc_high := cpu.devices[0].read()
		new_pc_low := cpu.devices[0].read()
		new_pc := uint16(new_pc_high)<<8 + uint16(new_pc_low)

		// Push next address to the stack
		pc_high := cpu.pc >> 8
		pc_low := cpu.pc << 8 >> 8
		cpu.devices[0].write(uint8(pc_low))
		cpu.devices[0].write(uint8(pc_high))

		// Jump
		cpu.pc = new_pc
	case 21:
		// trap
	case 22:
		// Get return address off the stack
		pc_high := cpu.devices[0].read()
		pc_low := cpu.devices[0].read()
		pc := uint16(pc_high)<<8 + uint16(pc_low)

		// Jump
		cpu.pc = pc
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
