package assemble

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
