package utils

// crypto
import (
	"encoding/base64"
)

func Base64Decode(src []byte) ([]byte, error) {
	buf := make([]byte, len(src))
	dstlen, err := base64.StdEncoding.Decode(buf, src)
	if err != nil {
		return nil, err
	}
	return buf[0:dstlen], nil
}

func Base64Encode(src []byte) []byte {
	buf := make([]byte, len(src)*2)
	base64.StdEncoding.Encode(buf, src)
	dstlen := base64.StdEncoding.EncodedLen(len(src))
	return buf[0:dstlen]
}
