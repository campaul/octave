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

	if err != nil {
		return
	}

	for cpu.running {
		inst_byte := fetch(cpu)
		inst_func := decode(inst_byte)
		inst_func(inst_byte, cpu)

		fmt.Printf("R0: %v\n", cpu.registers[0])
		fmt.Printf("R1: %v\n", cpu.registers[1])
		fmt.Printf("R2: %v\n", cpu.registers[2])
		fmt.Printf("R3: %v\n", cpu.registers[3])
		fmt.Println("")
	}
}

type CPU struct {
	memory    []uint8
	registers [4]uint8
	pc        uint16
	running   bool
}

type instruction func(uint8, *CPU)

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
}

func logic(i uint8, cpu *CPU) {
}

func mem(i uint8, cpu *CPU) {
}

func stack(i uint8, cpu *CPU) {
}

func in(i uint8, cpu *CPU) {
}

func out(i uint8, cpu *CPU) {
}

func illegal(i uint8, cpu *CPU) {
}
