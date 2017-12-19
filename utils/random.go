package utils

// random
import (
	"crypto/rand"
	mstring "github.com/weihualiu/toolkit/string"
	"time"
)

func Random() []byte {
	// 当前时间
	t := mstring.IntToBytes(time.Now().Nanosecond())

	c := 28
	b := make([]byte, c)
	rand.Read(b)

	r := make([]byte, 32)
	copy(r[0:4], t)
	copy(r[4:], b)

	return r
}

func Random2(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)

	return b
}
