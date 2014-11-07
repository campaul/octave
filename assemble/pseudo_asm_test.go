package assemble

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type ptestcase struct {
	in  string
	out []instruction
}

func TestGenerateMov(t *testing.T) {
	tests := []ptestcase{
		ptestcase{"MOV R0, R3", []instruction{tworeg{xor, 0, 0}, tworeg{add, 0, 3}}},
	}
	for _, test := range tests {
		out := generateMov(test.in).translate(map[string]uint{}, 0)
		if !reflect.DeepEqual(out, test.out) {
			t.Error("Got", out, "expected", test.out)
		}
	}
}

func TestGenerateLoadi(t *testing.T) {
	tests := []ptestcase{
		ptestcase{"LOADI 0xFF", []instruction{loadi{true, 0xF}, loadi{false, 0xF}}},
		ptestcase{"LOADI 0xF0", []instruction{loadi{true, 0x0}, loadi{false, 0xF}}},
	}
	for _, test := range tests {
		out := generateLoadi(test.in).translate(map[string]uint{}, 0)
		if !reflect.DeepEqual(out, test.out) {
			t.Error("Got", out, "expected", test.out)
		}
	}
}

func TestLra(t *testing.T) {
	bytes, err := Assemble(strings.NewReader(`
	LRA R1, Thing
	Thing:
	JMP R1 NZP
	`))
	if err != nil {
		t.Error("Unexpected error", err)
	}
	expected := []byte{
		binaryHelperByte("11100000"),
		binaryHelperByte("00110000"),
		binaryHelperByte("00100000"),
		binaryHelperByte("11000000"),
		binaryHelperByte("00001111"),
	}
	if !reflect.DeepEqual(bytes, expected) {
		t.Error("Got", bytesToAsm(bytes), "expected", bytesToAsm(expected))
	}
}

func bytesToAsm(bytes []byte) string {
	str := "[ "
	for _, b := range bytes {
		str += strconv.FormatUint(uint64(b), 2) + " "
	}
	str += "]"
	return str
}

func TestBytes(t *testing.T) {
	bytes, err := Assemble(strings.NewReader(`
	BYTES "butt\n"
	`))
	if err != nil {
		t.Error("Unexpected error", err)
	}
	expected := []byte("butt\n")
	if !reflect.DeepEqual(bytes, expected) {
		t.Error("Got", bytesToAsm(bytes), "expected", bytesToAsm(expected))
	}
}
