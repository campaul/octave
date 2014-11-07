package assemble

const (
	add   = 0x4 // 0100
	div   = 0x5 // 0101
	and   = 0x6 // 0110
	xor   = 0x7 // 0111
	load  = 0x8 // 1000
	store = 0x9 // 1001
)

const (
	in  = 0x6 //110
	out = 0x7 //111
)

var strToOpcode = map[string]uint8{
	"ADD":   add,
	"DIV":   div,
	"AND":   and,
	"XOR":   xor,
	"LOAD":  load,
	"STORE": store,
}

var strToDevioOpcode = map[string]uint8{
	"IN":  in,
	"OUT": out,
}

type instruction interface {
	assemble() byte
}

type jmp struct {
	register uint8
	negative bool
	zero     bool
	positive bool
}

const jmpOpcode = 0x0 // 000
func (j jmp) assemble() (b byte) {
	b |= jmpOpcode << 5
	b |= j.register << 3
	b |= bootToUint8(j.negative) << 2
	b |= bootToUint8(j.zero) << 1
	b |= bootToUint8(j.positive)
	return
}

func bootToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

type loadi struct {
	low    bool
	nibble uint8
}

const loadiOpcode = 0x1 // 001
func (i loadi) assemble() (b byte) {
	b |= loadiOpcode << 5
	b |= bootToUint8(i.low) << 4
	b |= i.nibble
	return
}

type tworeg struct {
	opcode uint8
	dest   uint8
	src    uint8
}

func (i tworeg) assemble() (b byte) {
	b |= i.opcode << 4
	b |= i.dest << 2
	b |= i.src
	return
}

type stackop struct {
	op uint8
}

const stackopOpcode = 0x5 // 101
func (i stackop) assemble() (b byte) {
	b |= stackopOpcode << 5
	b |= i.op
	return
}

type devio struct {
	opcode   uint8
	register uint8
	device   uint8
}

func (i devio) assemble() (b byte) {
	b |= i.opcode << 5
	switch i.opcode {
	case in:
		b |= i.register << 3
		b |= i.device
	case out:
		b |= i.device << 2
		b |= i.register
	}
	return
}

type rawbyte struct {
	value uint8
}

func (i rawbyte) assemble() byte {
	return i.value
}
