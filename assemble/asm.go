package assemble

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

func assemble(in io.Reader) ([]byte, error) {
	_ = bufio.NewReader(in)
	return []byte{}, nil
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
	b |= boolToUint8(j.negative) << 2
	b |= boolToUint8(j.zero) << 1
	b |= boolToUint8(j.positive)
	return
}

func assembleJmp(i string) instruction {
	var err error
	j := jmp{}
	re := regexp.MustCompile("JMP R([0-3]) (N|)(Z|)(P|)")
	matches := re.FindStringSubmatch(i)
	if matches == nil {
		panic(errors.New("not a JMP"))
	}

	reg, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}
	j.register = uint8(reg)
	if j.register < 0 || j.register > 3 {
		panic(errors.New(fmt.Sprint(j.register, " is not 0-3")))
	}

	j.negative = matches[2] == "N"
	j.zero = matches[3] == "Z"
	j.positive = matches[4] == "P"

	return j
}

func binaryHelper(bitstr string) (b uint8) {
	if len(bitstr) != 8 {
		panic(errors.New("bitstr must be of length 8"))
	}

	for i, c := range bitstr {
		shift := uint8(7 - i)
		var num uint8 = 0
		if c == '1' {
			num = 1
		}
		b |= num << shift
	}
	return
}

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
