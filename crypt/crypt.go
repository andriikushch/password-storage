package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"
)

const PADDING = "\x00"

var (
	ErrCipherTextRead      = errors.New("read chipher text error")
	ErrCipherCreation      = errors.New("new chipher creation error")
	ErrCipherTooShort      = errors.New("chipher too shortr")
	ErrCipherTextWrongSize = errors.New("ciphertext is not a multiple of the block size")
)

func Encrypt(key []byte, password string) ([]byte, error) {

	for len(password)%aes.BlockSize != 0 {
		password = password + PADDING
	}

	plaintext := []byte(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		return nil, ErrCipherCreation
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Printf("%v", err)
		return nil, ErrCipherTextRead
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func Decrypt(key []byte, encryptedpassword []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("%v", err)
		return "", ErrCipherCreation
	}

	if len(encryptedpassword) < aes.BlockSize {
		fmt.Printf("%v", err)
		return "", ErrCipherTooShort
	}
	iv := encryptedpassword[:aes.BlockSize]
	encryptedpassword = encryptedpassword[aes.BlockSize:]

	if len(encryptedpassword)%aes.BlockSize != 0 {
		fmt.Printf("%v", err)
		return "", ErrCipherTextWrongSize
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptedpassword, encryptedpassword)

	return strings.TrimRight(string(encryptedpassword), PADDING), nil
}
