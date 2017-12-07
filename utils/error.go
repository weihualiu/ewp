package utils

// error struct

import (
	"bytes"
)

type Error struct {
	Buffer []byte
	Str string
}

func NewError(str string) *Error {
	err := new(Error)
	err.Str = str
	return err
}

func NewErrorByte(buf []byte) *Error {
	err := new(Error)
	err.Buffer = buf
	return err
}

func (this Error)Error() []byte {
	if this.Buffer != nil {
		return this.Buffer
	}else if this.Str == "" {
		return bytes.NewBufferString(this.Str).Bytes()
	}else{
		return make([]byte, 0)
	}
}

