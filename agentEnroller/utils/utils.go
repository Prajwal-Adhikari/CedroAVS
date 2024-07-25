package utils

import "bytes"

func StringToBytes32(s string) [32]byte {
	var b [32]byte
	copy(b[:], s)
	return b
}

// Convert bytes32 to string
func Bytes32ToString(b [32]byte) string {
	n := bytes.IndexByte(b[:], 0)
	if n == -1 {
		n = 32
	}
	return string(b[:n])
}
