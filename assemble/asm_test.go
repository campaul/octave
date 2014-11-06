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
		testcase{"JMP R3 NZP", binaryHelper("00011111")},
		testcase{"JMP R2 NZ", binaryHelper("00010110")},
		testcase{"JMP R1 Z", binaryHelper("00001010")},
		testcase{"JMP R0 P", binaryHelper("00000001")},
	}
	for _, test := range tests {
		b := assembleJmp(test.in).assemble()
		if b != test.out {
			t.Error("Got", b, "expected", test.out)
		}
	}
}
