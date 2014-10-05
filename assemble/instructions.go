package assemble

const (
	add   = 0x4 // 0100
	div   = 0x5 // 0101
	and   = 0x6 // 0110
	xor   = 0x7 // 0111
	load  = 0x8 // 1000
	store = 0x9 // 1001
)

var strToOpcode = map[string]uint8{
	"ADD":   add,
	"DIV":   div,
	"AND":   and,
	"XOR":   xor,
	"LOAD":  load,
	"STORE": store,
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

func (j jmp) assemble() (b byte) {
	b |= 0x0 << 5
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

func (i loadi) assemble() (b byte) {
	b |= binaryHelperByte("001") << 5
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
