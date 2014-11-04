package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Initializing Octave CPU...")

	cpu := new(CPU)
	cpu.running = true

	file, err := os.Open(os.Args[1])

	if err != nil {
		return
	}

	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		return
	}

	cpu.memory = make([]uint8, stat.Size())
	_, err = file.Read(cpu.memory)

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
	memory  []uint8
	r0      uint8
	r1      uint8
	r2      uint8
	r3      uint8
	pc      uint16
	running bool
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
		inst = mem
	case 1:
		inst = loadi
	case 2:
		if (i << 3) < 25 {
			inst = stack
		} else {
			inst = inte
		}
	case 3:
		inst = jmp
	case 4:
		inst = math
	case 5:
		inst = logic
	case 6:
		inst = in
	case 7:
		inst = out
	}

	return inst
}

func mem(i uint8, cpu *CPU) {
	fmt.Println("mem")
}

func loadi(i uint8, cpu *CPU) {
	fmt.Println("loadi")
}

func stack(i uint8, cpu *CPU) {
	fmt.Println("stack")
}

func inte(i uint8, cpu *CPU) {
	fmt.Println("inte")
}

func jmp(i uint8, cpu *CPU) {
	if i == 96 {
		cpu.running = false
	}

	fmt.Println("jmp")
}

func math(i uint8, cpu *CPU) {
	fmt.Println("math")
}

func logic(i uint8, cpu *CPU) {
	fmt.Println("logic")
}

func in(i uint8, cpu *CPU) {
	fmt.Println("in")
}

func out(i uint8, cpu *CPU) {
	fmt.Println("out")
}

func illegal(i uint8, cpu *CPU) {
	fmt.Println("illegal")
}
