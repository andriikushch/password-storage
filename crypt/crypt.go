package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"strings"
)

const padding = "\x00"

var (
	ErrCipherTextRead      = errors.New("read chipher text error")
	ErrCipherCreation      = errors.New("new chipher creation error")
	ErrCipherTooShort      = errors.New("chipher too shortr")
	ErrCipherTextWrongSize = errors.New("ciphertext is not a multiple of the block size")
)

// Encrypt adding padding to plain message, ensuring that it is matching length required by AES.
// Returns encrypted message
func Encrypt(key []byte, message string) ([]byte, error) {
	for len(message)%aes.BlockSize != 0 {
		message = message + padding
	}

	plaintext := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err.Error())
		return nil, ErrCipherCreation
	}

	encryptedText := make([]byte, aes.BlockSize+len(plaintext))
	iv := encryptedText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Println(err.Error())
		return nil, ErrCipherTextRead
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptedText[aes.BlockSize:], plaintext)

	return encryptedText, nil
}

// Returns encrypted message and trims padding
func Decrypt(key []byte, encryptedMessage []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err.Error())
		return "", ErrCipherCreation
	}

	if len(encryptedMessage) < aes.BlockSize {
		log.Println(err.Error())
		return "", ErrCipherTooShort
	}
	iv := encryptedMessage[:aes.BlockSize]
	encryptedMessage = encryptedMessage[aes.BlockSize:]

	if len(encryptedMessage)%aes.BlockSize != 0 {
		log.Println(err.Error())
		return "", ErrCipherTextWrongSize
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptedMessage, encryptedMessage)

	return strings.TrimRight(string(encryptedMessage), padding), nil
}
