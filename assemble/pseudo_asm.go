package assemble

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var pseudoMap = map[string]func(string) psuedoinst{
	"MOV": generateMov,
}

func tryPseudo(line string) psuedoinst {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return dummy{}
	}
	f, ok := pseudoMap[fields[0]]
	if !ok {
		return dummy{}
	}
	return f(line)
}

func generateMov(line string) psuedoinst {
	m := mov{}
	re := regexp.MustCompile("MOV R([0-3]), R([0-3])")
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		panic(errors.New("not a MOV"))
	}
	m.dest = convertRegisterNum(matches[1])
	m.src = convertRegisterNum(matches[2])
	return m
}

func generateLoadi(line string) psuedoinst {
	l := loadimm{}
	re := regexp.MustCompile("LOADI (0x[0-9A-F]{1,2})")
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		panic(errors.New("not a LOADI"))
	}
	val, err := strconv.ParseUint(matches[1], 0, 8)
	if err != nil {
		panic(err)
	}
	l.val = uint8(val)
	return l
}
