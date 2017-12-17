package utils

// aes
import (
	"crypto/aes"
	"crypto/cipher"
	"bytes"
)

func AesEncrypt(plainText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plainText = pkcs5Padding(plainText, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockModel.CryptBlocks(cipherText, plainText)
	return cipherText, nil
}


func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padtext...)
}

func AesDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plainText, cipherText)
	plainText = pkcs5UnPadding(plainText, block.BlockSize())
	return plainText, nil
}

func pkcs5UnPadding(plainText []byte, blockSize int) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])
	return plainText[:(length - unpadding)]
}
