package string

import (
	"encoding/binary"
	"bytes"
)

func BytesToUInt32(buf []byte) uint32 {
	return uint32(binary.BigEndian.Uint32(buf))
}

func BytesToUInt16(buf []byte) uint16 {
	return uint16(binary.BigEndian.Uint16(buf))
}

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func UInt16ToBytes(n uint16) []byte {
	b := bytes.NewBuffer([]byte{})
	binary.Write(b, binary.BigEndian, n)
	return b.Bytes()
}

func BytesToString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func BytesTrim(c []byte) []byte {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return c[:n+1]
}

func IsExistUpper(s string) bool {
	f := false
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x41 && s[i] <= 0x5A {
			f = true
			break
		}
	}
	return f
}
