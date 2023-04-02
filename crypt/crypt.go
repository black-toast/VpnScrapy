package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

var cryptKey = []byte("tlynet923456789k")
var cryptKIv = []byte("9987654321fedcsu")

func PKCSNoPadding(ciphertext []byte) []byte {
	length := len(ciphertext)
	if length%16 != 0 {
		padding := 16 - length%16
		padtext := bytes.Repeat([]byte{byte(0)}, padding)
		return append(ciphertext, padtext...)
	}
	return ciphertext
}

func AesCBCEncrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(cryptKey)
	if err != nil {
		return nil, err
	}
	origData = PKCSNoPadding(origData)
	blockMode := cipher.NewCBCEncrypter(block, cryptKIv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesCBCDecrypt(encryptData []byte) ([]byte, error) {
	block, err := aes.NewCipher(cryptKey)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, cryptKIv)
	origData := make([]byte, len(encryptData))
	blockMode.CryptBlocks(origData, encryptData)
	return origData, nil
}
