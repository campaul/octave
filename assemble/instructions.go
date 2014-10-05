package assemble

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
