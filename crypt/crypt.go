package crypt

import (
	"VpnScrapy/util"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func PKCSNoPadding(ciphertext []byte) []byte {
	length := len(ciphertext)
	if length%16 != 0 {
		padding := 16 - length%16
		padtext := bytes.Repeat([]byte{byte(0)}, padding)
		return append(ciphertext, padtext...)
	}
	return ciphertext
}

func pkcsUnPadding(origData []byte) []byte {
	length := len(origData)
	for i := 1; i <= 16; i++ {
		if origData[length-i] != byte(0) {
			return origData[:(length - i + 1)]
		}
	}
	return origData
}

func AesCBCEncrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(util.WithConstant().GetAesKey())
	if err != nil {
		return nil, err
	}
	origData = PKCSNoPadding(origData)
	blockMode := cipher.NewCBCEncrypter(block, util.WithConstant().GetAesIv())
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesCBCDecrypt(encryptData []byte) ([]byte, error) {
	block, err := aes.NewCipher(util.WithConstant().GetAesKey())
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, util.WithConstant().GetAesIv())
	origData := make([]byte, len(encryptData))
	blockMode.CryptBlocks(origData, encryptData)
	origData = pkcsUnPadding(origData)
	return origData, nil
}
