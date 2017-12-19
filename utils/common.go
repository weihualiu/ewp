package utils

// bytes, string int
import (
	"bytes"
)

// 合并多个字节数组
func BytesAppend(b ...[]byte) []byte {
	tmp := bytes.NewBuffer([]byte{})
	for i := 0; i < len(b); i++ {
		tmp.Write(b[i])
	}
	return tmp.Bytes()
}
