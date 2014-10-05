package assemble

import (
	"testing"
)

type testcase struct {
	in  string
	out byte
}

func TestAssembleJmp(t *testing.T) {
	tests := []testcase{
		testcase{"JMP R3 NZP", binaryHelperByte("00011111")},
		testcase{"JMP R2 NZ", binaryHelperByte("00010110")},
		testcase{"JMP R1 Z", binaryHelperByte("00001010")},
		testcase{"JMP R0 P", binaryHelperByte("00000001")},
	}
	for _, test := range tests {
		b := assembleJmp(test.in).assemble()
		if b != test.out {
			t.Error("Got", b, "expected", test.out)
		}
	}
}

func TestAssembleLoadi(t *testing.T) {
	tests := []testcase{
		testcase{"LOADIL 0xF", binaryHelperByte("00111111")},
		testcase{"LOADIH 017", binaryHelperByte("00101111")},
		testcase{"LOADIH 15", binaryHelperByte("00101111")},
		testcase{"LOADIL 0", binaryHelperByte("00110000")},
	}
	for _, test := range tests {
		b := assembleLoadi(test.in).assemble()
		if b != test.out {
			t.Error("Got", b, "expected", test.out)
		}
	}
}
