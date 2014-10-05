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

func assembleLoadi(i string) instruction {
	in := loadi{}
	re := regexp.MustCompile("LOADI(L|H) (0x[0-9A-F]|0[0-9]+|[0-9]+)")
	matches := re.FindStringSubmatch(i)
	if matches == nil {
		panic(errors.New("not a LOADI(L|H)"))
	}
	lr := matches[1]
	in.low = (lr == "L")

	nibble, err := strconv.ParseUint(matches[2], 0, 4)
	if err != nil {
		panic(err)
	}
	in.nibble = uint8(nibble)
	return in
}
