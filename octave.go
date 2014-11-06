package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Initializing Octave CPU...")
	cpu := &CPU{running: true}

	file, err := os.Open(os.Args[1])

	if err != nil {
		return
	}

	defer file.Close()

	cpu.memory, err = ioutil.ReadAll(file)

	cpu.devices[1] = tty{}

	if err != nil {
		return
	}

	for cpu.running {
		inst_byte := fetch(cpu)
		inst_func := decode(inst_byte)
		inst_func(inst_byte, cpu)

		//fmt.Printf("R0: %v\n", cpu.registers[0])
		//fmt.Printf("R1: %v\n", cpu.registers[1])
		//fmt.Printf("R2: %v\n", cpu.registers[2])
		//fmt.Printf("R3: %v\n", cpu.registers[3])
		//fmt.Println("")
	}
}

type CPU struct {
	memory    []uint8
	registers [4]uint8
	pc        uint16
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
}

func (t tty) read() uint8 {
	return 0;
}

func (t tty) write(char uint8) {
	fmt.Printf("%c", char)
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
		inst = jmp
	case 1:
		inst = loadi
	case 2:
		inst = math
	case 3:
		inst = logic
	case 4:
		inst = mem
	case 5:
		inst = stack
	case 6:
		inst = in
	case 7:
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
	destination := i << 4 >> 6
	source := i << 6 >> 6

	if operation == 0 {
		// LOAD
		address := uint16(cpu.registers[0]) << 8 + uint16(cpu.registers[source])
		cpu.registers[destination] = cpu.memory[address]
	} else {
		// STORE
		address := uint16(cpu.registers[0]) << 8 + uint16(cpu.registers[destination])
		cpu.memory[address] = cpu.registers[source]
	}
}

func stack(i uint8, cpu *CPU) {
}

func in(i uint8, cpu *CPU) {
}

func out(i uint8, cpu *CPU) {
	device := i << 3 >> 5
	source := i << 6 >> 6
	cpu.devices[device].write(cpu.registers[source])
}

func illegal(i uint8, cpu *CPU) {
}
