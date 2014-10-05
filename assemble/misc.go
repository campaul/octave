package assemble

import (
	"errors"
	"fmt"
	"strconv"
)

func binaryHelperByte(bitstr string) (b byte) {
	return byte(binaryHelper(bitstr))
}

func binaryHelper(bitstr string) (b uint32) {
	for i, c := range bitstr {
		shift := uint32(len(bitstr) - 1 - i)
		var num uint32 = 0
		if c == '1' {
			num = 1
		}
		b |= num << shift
	}
	return
}

func convertRegisterNum(s string) (num uint8) {
	reg, err := strconv.ParseUint(s, 10, 2)
	if err != nil {
		panic(err)
	}
	num = uint8(reg)
	if num < 0 || num > 3 {
		panic(errors.New(fmt.Sprint(num, " is not 0-3")))
	}
	return
}
