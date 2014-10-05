package assemble

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
)

func assemble(in io.Reader) ([]byte, error) {
	_ = bufio.NewReader(in)
	return []byte{}, nil
}

func assembleJmp(i string) instruction {
	j := jmp{}
	re := regexp.MustCompile("JMP R([0-3]) (N|)(Z|)(P|)")
	matches := re.FindStringSubmatch(i)
	if matches == nil {
		panic(errors.New("not a JMP"))
	}

	j.register = convertRegisterNum(matches[1])

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

func assembleTwoRegister(i string) instruction {
	in := tworeg{}
	re := regexp.MustCompile("(ADD|DIV|AND|XOR|LOAD|STORE) R([0-3]), R([0-3])")
	matches := re.FindStringSubmatch(i)
	if matches == nil {
		panic(errors.New("not a (ADD|DIV|AND|XOR|LOAD|STORE)"))
	}
	in.opcode = strToOpcode[matches[1]]
	in.dest = convertRegisterNum(matches[2])
	in.src = convertRegisterNum(matches[3])
	return in
}

func assembleDeviceIO(i string) instruction {
	in := devio{}
	re := regexp.MustCompile("(IN|OUT) R([0-3]), ([0-7])")
	matches := re.FindStringSubmatch(i)
	if matches == nil {
		panic(errors.New("not a (IN|OUT)"))
	}
	in.opcode = strToDevioOpcode[matches[1]]
	in.register = convertRegisterNum(matches[2])
	in.device = convertDeviceNum(matches[3])
	return in
}
