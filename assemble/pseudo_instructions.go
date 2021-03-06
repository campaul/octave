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
	if !ok {
		panic(errors.New(fmt.Sprint("%v label not found", i.label)))
	}
	offset := uint8(int(addr) - int(pc+i.size()+1))
	insts := []instruction{}
	insts = append(insts, loadimm{offset}.translate(labels, pc)...)
	return insts
}

func (i lra) size() uint {
	return 2
}

type laa struct {
	highreg uint8
	lowreg  uint8
	label   string
}

func (i laa) translate(labels map[string]uint, pc uint) []instruction {
	addr, ok := labels[i.label]
	if !ok {
		panic(errors.New(fmt.Sprint("%v label not found", i.label)))
	}
	insts := []instruction{}
	if i.hasR0() {
		if !i.highR0() {
			insts = append(insts, loadimm{uint8((addr & 0xFF00) >> 8)}.translate(labels, pc)...)
			insts = append(insts, tworeg{mov, i.highreg, 0})
		}
		insts = append(insts, loadimm{uint8(addr & 0xFF)}.translate(labels, pc)...)
		insts = append(insts, tworeg{mov, i.lowreg, 0})
		if i.highR0() {
			insts = append(insts, loadimm{uint8((addr & 0xFF00) >> 8)}.translate(labels, pc)...)
			insts = append(insts, tworeg{mov, i.highreg, 0})
		}
	} else {
		insts = append(insts, loadimm{uint8((addr & 0xFF00) >> 8)}.translate(labels, pc)...)
		insts = append(insts, tworeg{mov, i.highreg, 0})
		insts = append(insts, loadimm{uint8(addr & 0xFF)}.translate(labels, pc)...)
		insts = append(insts, tworeg{mov, i.lowreg, 0})
	}
	return insts
}

func (i laa) hasR0() bool {
	return i.highreg != 0 || i.lowreg != 0
}

func (i laa) highR0() bool {
	return i.highreg == 0
}

func (i laa) size() uint {
	return 2 + 1 + 2 + 1
}

type rawbytes struct {
	bytes []byte
}

func (i rawbytes) translate(labels map[string]uint, pc uint) (out []instruction) {
	for _, b := range i.bytes {
		out = append(out, rawbyte{b})
	}
	return
}

func (i rawbytes) size() uint {
	return uint(len(i.bytes))
}

type fill struct {
	num uint16
}

func (in fill) translate(labels map[string]uint, pc uint) (out []instruction) {
	for i := 0; i < int(in.num); i++ {
		out = append(out, rawbyte{0})
	}
	return
}

func (i fill) size() uint {
	return uint(i.num)
}

type ljump struct {
	label string
	neg   bool
	zero  bool
	pos   bool
}

func (i ljump) translate(labels map[string]uint, pc uint) (o []instruction) {
	addr, ok := labels[i.label]
	if !ok {
		panic(errors.New(fmt.Sprint("%v label not found", i.label)))
	}

	o = append(o, loadimm{7}.translate(labels, pc)...)
	if i.isUnconditional() {
		o = append(o, jmp{0, !i.neg, !i.zero, !i.pos})
	}
	o = append(o, loadimm{uint8(addr & 0xFF)}.translate(labels, pc)...)
	o = append(o, devio{out, 0, 0})
	o = append(o, loadimm{uint8((addr & 0xFF00) >> 8)}.translate(labels, pc)...)
	o = append(o, devio{out, 0, 0})
	o = append(o, stackop{22})
	return
}

func (i ljump) size() (s uint) {
	s = 2 + 2 + 1 + 2 + 1 + 1
	if i.isUnconditional() {
		s += 1
	}
	return
}

func (i ljump) isUnconditional() bool {
	return !(i.neg && i.zero && i.pos)
}
