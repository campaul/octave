package assemble

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var codeToFunc = map[string]func(string) instruction{
	"JMP":     assembleJmp,
	"LOADIH":  assembleLoadi,
	"LOADIL":  assembleLoadi,
	"ADD":     assembleTwoRegister,
	"DIV":     assembleTwoRegister,
	"AND":     assembleTwoRegister,
	"XOR":     assembleTwoRegister,
	"LOAD":    assembleTwoRegister,
	"STORE":   assembleTwoRegister,
	"STACKOP": assembleStackop,
	"IN":      assembleDeviceIO,
	"OUT":     assembleDeviceIO,
}

func Assemble(in io.Reader) (bytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	r := bufio.NewReader(in)
	for line, err := r.ReadString('\n'); err != io.EOF; line, err = r.ReadString('\n') {
		line = strings.TrimSpace(line)
		if err != nil {
			return []byte{}, err
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		asmFunc, ok := codeToFunc[fields[0]]
		if !ok {
			return []byte{}, errors.New(
				fmt.Sprintf("'%v' is not a valiad assembly instruction", fields[0]),
			)
		}
		bytes = append(bytes, asmFunc(line).assemble())
	}
	return bytes, nil
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
	re := regexp.MustCompile("LOADI(L|H) (0x[0-9A-F]|0[0-7]+|[0-9]+)")
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

func assembleStackop(i string) instruction {
	in := stackop{}
	re := regexp.MustCompile("STACKOP ([0-9]{1,2})")
	matches := re.FindStringSubmatch(i)
	if matches == nil {
		panic(errors.New("not a STACKOP"))
	}
	op, err := strconv.ParseUint(matches[1], 0, 5)
	if err != nil {
		panic(err)
	}
	in.op = uint8(op)
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
