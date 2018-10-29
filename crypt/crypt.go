package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"strings"
)

const PADDING = "\x00"

func Encrypt(key []byte, password string) []byte {

	for len(password)%aes.BlockSize != 0 {
		password = password + PADDING
	}

	plaintext := []byte(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}

func Decrypt(key []byte, encryptedpassword []byte) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(encryptedpassword) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encryptedpassword[:aes.BlockSize]
	encryptedpassword = encryptedpassword[aes.BlockSize:]

	if len(encryptedpassword)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptedpassword, encryptedpassword)

	return strings.TrimRight(string(encryptedpassword), PADDING)
}
