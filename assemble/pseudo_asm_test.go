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
