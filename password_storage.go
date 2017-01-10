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
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"errors"
	"flag"
	"strings"
	"github.com/andriikushch/clipboard"
)

const PASSWORD_ACCOUNT_NAME_SEPARATOR = "\x01"
const PADDING = "\x00"

var (
	masterPassword = "example key 1234"
	databaseFile = "dat2"
)

func main() {
	key := sha256.Sum256([]byte(masterPassword))

	command := flag.String("command", "add-new-credentials", "a string")
	flag.Parse()

	switch *command {
	case "add-new-credentials":
		if err := addNewCredentials(key[:]); err != nil {
			fmt.Errorf("%v", err)
		}
	case "load-password":
		var account string
		fmt.Print("Enter account name: ")
		fmt.Scanln(&account)
		loadDBAndDecryptPassword(key[:], account)
	case "accounts":
		showAccountsList(key[:])
	}
}

func addNewCredentials(key []byte) error {
	var account string

	fmt.Print("Enter account name: ")
	fmt.Scanln(&account)

	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("%v", err)
		return errors.New("Can't read password")
	}
	password := string(bytePassword)

	fmt.Print("\nEnter password confirmation: ")
	bytePasswordConfirmation, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return errors.New("Can't read password confiramtion")
	}
	passwordConfirmation := string(bytePasswordConfirmation)

	if password == passwordConfirmation {
		return storeAccountPasswordPair(encrypt(key, account + PASSWORD_ACCOUNT_NAME_SEPARATOR + password))
	}

	return errors.New("Password and Password confirmation is not equal")
}

func storeAccountPasswordPair(encryptedCredentials []byte) error {
	f, err := os.OpenFile(databaseFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)

	defer f.Close()

	if _, err = f.Write(encryptedCredentials); err != nil {
		return errors.New("Can't write credentials")
	}

	return f.Sync()
}

func loadDBAndDecryptPassword(key []byte, account string) error {
	lines, err := readLines(databaseFile)

	if err != nil {
		return errors.New("Error while reading")
	}

	for _,element := range lines {
		line := decrypt(key, []byte(element))
		accountPasswordPair := strings.Split(line, PASSWORD_ACCOUNT_NAME_SEPARATOR)

		if accountPasswordPair[0] == account {
			fmt.Println("Password found")
			clipboard.WriteAll(accountPasswordPair[1])
		}
	}

	return errors.New("Can't find password for account")
}

func showAccountsList(key []byte) error {
	lines, err := readLines(databaseFile)

	if err != nil {
		return errors.New("Error while reading")
	}

	for _,element := range lines {
		line := decrypt(key, []byte(element))
		accountPasswordPair := strings.Split(line, PASSWORD_ACCOUNT_NAME_SEPARATOR)
		fmt.Println(accountPasswordPair[0])
	}

	return nil
}

func encrypt(key []byte, password string) []byte {

	for len(password) % aes.BlockSize != 0 {
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

	return  ciphertext
}

func decrypt(key []byte, encryptedpassword []byte) string {

	ciphertext := encryptedpassword

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

	return strings.TrimRight(string(ciphertext), PADDING)
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