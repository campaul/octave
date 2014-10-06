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
