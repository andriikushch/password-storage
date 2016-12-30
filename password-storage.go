package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"crypto/sha256"
	"os"
	"bufio"
	"reflect"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"errors"
)

func main() {
	masterPassword := "example key 1234"
	key := sha256.Sum256([]byte(masterPassword))
	var command string

	fmt.Print("Enter command: ")
	fmt.Scanln(&command)

	for {
		switch command {
		case "add-new-credentials":
			addNewCredentials(key[:])
		case "load-password":
			loadDBAndDecryptAllPassword(key[:])
		}
	}
}

func addNewCredentials(key []byte) error {
	var account string
	var password string
	var passwordConfirmation string

	fmt.Print("Enter account name: ")
	fmt.Scanln(&account)

	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return errors.New("Can't read password")
	}
	password = string(bytePassword)

	fmt.Print("Enter password confirmation: ")
	bytePasswordConfirmation, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return errors.New("Can't read password confiramtion")
	}
	passwordConfirmation = string(bytePasswordConfirmation)

	if password == passwordConfirmation {
		encryptedCredentials := encrypt(key, account + "\x01" + password)
		storeAccountPasswordPair(encryptedCredentials)
	}
}

func storeAccountPasswordPair(encryptedCredentials []byte)  {
	f, err := os.OpenFile("dat2", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)

	defer f.Close()

	_, err = f.Write(encryptedCredentials)
	check(err)

	f.Sync()
}

func loadDBAndDecryptAllPassword(key []byte) {
	var account string
	fmt.Print("Enter account name: ")
	fmt.Scanln(&account)

	lines, _ := readLines("dat2")

	for _,element := range lines {
		line := decrypt(key, []byte(element))

		var passwordFromFile []byte

		var toRecord bool = false
		for _,element := range line {
			if (element == 0) {
				toRecord = false
			}

			if (toRecord == true) {
				passwordFromFile = append(passwordFromFile, element)
			}

			if (element == 1) {
				toRecord = true
			}
		}

		if reflect.DeepEqual(passwordFromFile, []byte(account)) {
			fmt.Println("Password found")
		}
	}
}

func encrypt(key []byte, password string) []byte {

	for len(password) % aes.BlockSize != 0 {
		password = password + "\x00"
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

	return  ciphertext
}

func decrypt(key []byte, ciphertext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}