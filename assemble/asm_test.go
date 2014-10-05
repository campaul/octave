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

func TestAssembleTwoRegister(t *testing.T) {
	tests := []testcase{
		testcase{"ADD R0, R0", binaryHelperByte("01000000")},
		testcase{"DIV R1, R1", binaryHelperByte("01010101")},
		testcase{"AND R2, R2", binaryHelperByte("01101010")},
		testcase{"XOR R3, R3", binaryHelperByte("01111111")},
		testcase{"LOAD R0, R3", binaryHelperByte("10000011")},
		testcase{"STORE R1, R2", binaryHelperByte("10010110")},
	}
	for _, test := range tests {
		b := assembleTwoRegister(test.in).assemble()
		if b != test.out {
			t.Error("Got", b, "expected", test.out)
		}
	}
}

func TestAssembleStackop(t *testing.T) {
	tests := []testcase{
		testcase{"STACKOP 10", binaryHelperByte("10101010")},
	}
	for _, test := range tests {
		b := assembleStackop(test.in).assemble()
		if b != test.out {
			t.Error("Got", b, "expected", test.out)
		}
	}
}

func TestAssembleDeviceIO(t *testing.T) {
	tests := []testcase{
		testcase{"IN R0, 0", binaryHelperByte("11000000")},
		testcase{"OUT R3, 7", binaryHelperByte("11111111")},
	}
	for _, test := range tests {
		b := assembleDeviceIO(test.in).assemble()
		if b != test.out {
			t.Error("Got", b, "expected", test.out)
		}
	}
}
