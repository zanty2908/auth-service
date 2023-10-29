package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var encryptBytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
var mySecretKey = "qBq3AADtKXVrMg6I"

// Encrypt method is to encrypt or hide any classified text
func SimpleEncrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(mySecretKey))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, encryptBytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func SimpleDecrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(mySecretKey))
	if err != nil {
		return "", err
	}
	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, encryptBytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
