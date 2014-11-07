package assemble

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var pseudoMap = map[string]func(string) pseudoinst{
	"LOADI": generateLoadi,
	"LRA":   generateLra,
	"LAA":   generateLaa,
	"BYTES": generateBytes,
	"FILL":  generateFill,
}

func tryPseudo(line string) pseudoinst {
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

func generateLoadi(line string) pseudoinst {
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

func generateLra(line string) pseudoinst {
	l := lra{}
	re := regexp.MustCompile("LRA ([A-Za-z]+)")
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		panic(errors.New("not an LRA"))
	}
	l.label = matches[1]
	return l
}

func generateLaa(line string) pseudoinst {
	l := laa{}
	re := regexp.MustCompile("LAA R([0-3]), R([0-3]), ([A-Za-z]+)")
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		panic(errors.New("not an LAA"))
	}
	l.highreg = convertRegisterNum(matches[1])
	l.lowreg = convertRegisterNum(matches[2])
	l.label = matches[3]
	return l
}

func generateBytes(line string) pseudoinst {
	re := regexp.MustCompile("BYTES (\".*\")")
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		panic(errors.New("not a BYTES"))
	}
	s, err := strconv.Unquote(matches[1])
	if err != nil {
		panic(err)
	}
	return rawbytes{[]byte(s)}
}

func generateFill(line string) pseudoinst {
	re := regexp.MustCompile("FILL[ \\t]+([0-9]+)")
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		panic(errors.New("not a FILL"))
	}
	num, err := strconv.ParseUint(matches[1], 10, 16)
	if err != nil {
		panic(err)
	}
	return fill{uint16(num)}
}
