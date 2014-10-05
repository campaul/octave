package assemble

import (
	"strconv"
)

func binaryHelperByte(bitstr string) byte {
	num, err := strconv.ParseUint(bitstr, 2, 8)
	if err != nil {
		panic(err)
	}
	return byte(num)
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
