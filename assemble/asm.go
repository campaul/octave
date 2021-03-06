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
	"MOV":     assembleTwoRegister,
	"AND":     assembleTwoRegister,
	"XOR":     assembleTwoRegister,
	"LOAD":    assembleTwoRegister,
	"STORE":   assembleTwoRegister,
	"STACKOP": assembleStackop,
	"IN":      assembleDeviceIO,
	"OUT":     assembleDeviceIO,
	"BYTE":    assembleByte,
	"HALT":    assembleHalt,
	"ADD8":    assembleStackopAlias,
	"SUB8":    assembleStackopAlias,
	"MUL8":    assembleStackopAlias,
	"DIV8":    assembleStackopAlias,
	"MOD8":    assembleStackopAlias,
	"NEG8":    assembleStackopAlias,
	"AND8":    assembleStackopAlias,
	"OR8":     assembleStackopAlias,
	"XOR8":    assembleStackopAlias,
	"NOT8":    assembleStackopAlias,
	"ADD16":   assembleStackopAlias,
	"SUB16":   assembleStackopAlias,
	"MUL16":   assembleStackopAlias,
	"DIV16":   assembleStackopAlias,
	"MOD16":   assembleStackopAlias,
	"NEG16":   assembleStackopAlias,
	"AND16":   assembleStackopAlias,
	"OR16":    assembleStackopAlias,
	"XOR16":   assembleStackopAlias,
	"NOT16":   assembleStackopAlias,
	"CALL":    assembleStackopAlias,
	"TRAP":    assembleStackopAlias,
	"RET":     assembleStackopAlias,
	"IRET":    assembleStackopAlias,
	"INT0E":   assembleStackopAlias,
	"INT1E":   assembleStackopAlias,
	"INT2E":   assembleStackopAlias,
	"INT3E":   assembleStackopAlias,
	"INT4E":   assembleStackopAlias,
	"INT5E":   assembleStackopAlias,
	"INT6E":   assembleStackopAlias,
	"INT7E":   assembleStackopAlias,
	"PUSH":    assemblePushPop,
	"POP":     assemblePushPop,
}

func Assemble(in io.Reader) (bytes []byte, err error) {
	/*defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()*/
	r := bufio.NewReader(in)
	instructions, _, labels := convertToInstructions(r)
	pc := uint(0)
	for _, i := range instructions {
		switch i := i.(type) {
		case instruction:
			bytes = append(bytes, i.assemble())
			pc++
		case pseudoinst:
			for _, inst := range i.translate(labels, pc) {
				bytes = append(bytes, inst.assemble())
			}
			pc += i.size()
		}
	}
	return bytes, nil
}

func convertToInstructions(in *bufio.Reader) (instructions []interface{}, size uint, labels map[string]uint) {
	labels = make(map[string]uint)
	for line, err := in.ReadString('\n'); err != io.EOF; line, err = in.ReadString('\n') {
		line = strings.TrimSpace(line)
		if err != nil {
			panic(err)
		}

		if len(line) == 0 || line[0] == ';' {
			continue
		}
		fields := strings.Fields(line)

		label := tryLabel(line)
		if label != "" {
			labels[label] = size
			continue
		}

		pinst := tryPseudo(line)
		if pinst.size() != 0 {
			size += pinst.size()
			instructions = append(instructions, pinst)
			continue
		}

		asmFunc, ok := codeToFunc[fields[0]]
		if !ok {
			panic(errors.New(fmt.Sprintf("'%v' is not a valid assembly instruction", fields[0])))
		}
		instructions = append(instructions, asmFunc(line))
		size++
	}
	return
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
	re := regexp.MustCompile("(ADD|MOV|AND|XOR|LOAD|STORE) R([0-3]), R([0-3])")
	matches := re.FindStringSubmatch(i)
	if matches == nil {
		panic(errors.New(i + " not a (ADD|MOV|AND|XOR|LOAD|STORE)"))
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

func assembleByte(i string) instruction {
	re := regexp.MustCompile("BYTE (.+)")
	matches := re.FindStringSubmatch(i)
	if matches == nil {
		panic(errors.New("not a BYTE"))
	}
	val, err := strconv.ParseUint(matches[1], 0, 8)
	if err != nil {
		panic(err)
	}
	return rawbyte{uint8(val)}
}

func assembleHalt(i string) instruction {
	return rawbyte{0}
}

var stackopAliases = map[string]uint8{
	"ADD8":  0,
	"SUB8":  1,
	"MUL8":  2,
	"DIV8":  3,
	"MOD8":  4,
	"NEG8":  5,
	"AND8":  6,
	"OR8":   7,
	"XOR8":  8,
	"NOT8":  9,
	"ADD16": 10,
	"SUB16": 11,
	"MUL16": 12,
	"DIV16": 13,
	"MOD16": 14,
	"NEG16": 15,
	"AND16": 16,
	"OR16":  17,
	"XOR16": 18,
	"NOT16": 19,
	"CALL":  20,
	"TRAP":  21,
	"RET":   22,
	"IRET":  23,
	"INT0E": 24,
	"INT1E": 25,
	"INT2E": 26,
	"INT3E": 27,
	"INT4E": 28,
	"INT5E": 29,
	"INT6E": 30,
	"INT7E": 31,
}

func assembleStackopAlias(line string) instruction {
	return stackop{stackopAliases[line]}
}

func assemblePushPop(line string) instruction {
	re := regexp.MustCompile("(PUSH|POP)[ \\t]R([0-3])")
	matches := re.FindStringSubmatch(line)

	if matches == nil {
		panic(errors.New("not a (PUSH|POP)"))
	}

	reg := convertRegisterNum(matches[2])
	if matches[1] == "PUSH" {
		return devio{out, reg, 0}
	}
	return devio{in, reg, 0}
}
