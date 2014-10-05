package assemble

import (
	"strconv"
)

func binaryHelperByte(bitstr string) byte {
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

func convertRegisterNum(s string) uint8 {
	return convertNum(s, 2)
}

func convertDeviceNum(s string) uint8 {
	return convertNum(s, 3)
}

func convertNum(s string, bitdepth int) uint8 {
	reg, err := strconv.ParseUint(s, 10, bitdepth)
	if err != nil {
		panic(err)
	}
	return uint8(reg)
}
