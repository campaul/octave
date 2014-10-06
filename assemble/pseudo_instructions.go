package assemble

type psuedoinst interface {
	translate() []instruction
}

type dummy struct{}

func (d dummy) translate() []instruction {
	return []instruction{}
}

type mov struct {
	dest uint8
	src  uint8
}

func (i mov) translate() []instruction {
	return []instruction{
		tworeg{xor, i.dest, i.dest},
		tworeg{add, i.dest, i.src},
	}
}

type loadimm struct {
	val uint8
}

func (i loadimm) translate() []instruction {
	return []instruction{
		loadi{true, i.val & 0xF},
		loadi{false, (i.val & 0xF0) >> 4},
	}
}
