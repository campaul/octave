package assemble

import (
	"errors"
	"fmt"
)

type pseudoinst interface {
	translate(labels map[string]uint, pc uint) []instruction
	size() uint
}

type dummy struct{}

func (d dummy) translate(labels map[string]uint, pc uint) []instruction {
	return []instruction{}
}

func (d dummy) size() uint {
	return 0
}

type mov struct {
	dest uint8
	src  uint8
}

func (i mov) translate(labels map[string]uint, pc uint) []instruction {
	return []instruction{
		tworeg{xor, i.dest, i.dest},
		tworeg{add, i.dest, i.src},
	}
}

func (i mov) size() uint {
	return 2
}

type loadimm struct {
	val uint8
}

func (i loadimm) translate(labels map[string]uint, pc uint) []instruction {
	return []instruction{
		loadi{true, i.val & 0xF},
		loadi{false, (i.val & 0xF0) >> 4},
	}
}

func (i loadimm) size() uint {
	return 2
}

type lra struct {
	dest  uint8
	label string
}

func (i lra) translate(labels map[string]uint, pc uint) []instruction {
	addr, ok := labels[i.label]
	offset := uint8(int(addr) - int(pc+i.size()))
	if !ok {
		panic(errors.New(fmt.Sprint("%v label not found", i.label)))
	}
	insts := []instruction{}
	insts = append(insts, devio{out, 0, 0})
	insts = append(insts, loadimm{offset}.translate(labels, pc)...)
	insts = append(insts, mov{i.dest, 0}.translate(labels, pc)...)
	insts = append(insts, devio{in, 0, 0})
	return insts
}

func (i lra) size() uint {
	return 1 + 2 + 2 + 1
}
