package assemble

import (
	"reflect"
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
		out := generateMov(test.in).translate()
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
		out := generateLoadi(test.in).translate()
		if !reflect.DeepEqual(out, test.out) {
			t.Error("Got", out, "expected", test.out)
		}
	}
}
