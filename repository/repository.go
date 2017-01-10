package repository

import (
	"github.com/andriikushch/password-storage/crypt"
	"strings"
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"regexp"
	"errors"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
)

const PASSWORD_ACCOUNT_NAME_SEPARATOR = "\x01"

var (
	databaseFile          = "dat2"
)

func FindPassword(key []byte, account string) (string, error) {
	lines, err := readLines(databaseFile)

	if err != nil {
		return "", errors.New("Error while reading")
	}

	for _, element := range lines {
		line := crypt.Decrypt(key, []byte(element))
		accountPasswordPair := strings.Split(line, PASSWORD_ACCOUNT_NAME_SEPARATOR)

		if accountPasswordPair[0] == account {
			fmt.Println("Password found")
			return accountPasswordPair[1], nil
		}
	}

	return "", errors.New("Can't find password for account")
}

func isPasswordExist(key []byte, account string) bool {
	lines, err := readLines(databaseFile)

	if err != nil {
		return false
	}

	for _, element := range lines {
		line := crypt.Decrypt(key, []byte(element))
		accountPasswordPair := strings.Split(line, PASSWORD_ACCOUNT_NAME_SEPARATOR)

		if accountPasswordPair[0] == account {
			return true
		}
	}

	return false
}


func ShowAccountsList(key []byte) error {
	lines, err := readLines(databaseFile)

	if err != nil {
		return errors.New("Error while reading")
	}

	for _, element := range lines {
		line := crypt.Decrypt(key, []byte(element))
		accountPasswordPair := strings.Split(line, PASSWORD_ACCOUNT_NAME_SEPARATOR)
		fmt.Println(accountPasswordPair[0])
	}

	return nil
}

func updateAccountPasswordPair(encryptedCredentials []byte) error {
	data, err := ioutil.ReadFile(databaseFile)
	if err != nil {
		panic(err)
	}
	input := string(data)
	re := regexp.MustCompile(`^` + string(encryptedCredentials) + `\r?\n`)
	input = re.ReplaceAllString(input, "")

	f, err := os.OpenFile(databaseFile, os.O_RDWR|os.O_TRUNC, 0660)

	defer f.Close()

	stringToWrite := string(encryptedCredentials) + "\n"
	if _, err = f.WriteString(stringToWrite); err != nil {
		return errors.New("Can't write credentials")
	}

	return f.Sync()
}


func AddNewCredentials(key []byte) error {
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
		if isPasswordExist(key, account) {
			return updateAccountPasswordPair(crypt.Encrypt(key, account + PASSWORD_ACCOUNT_NAME_SEPARATOR + password))
		} else {
			return storeAccountPasswordPair(crypt.Encrypt(key, account + PASSWORD_ACCOUNT_NAME_SEPARATOR + password))
		}
	}

	return errors.New("Password and Password confirmation is not equal")
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

func storeAccountPasswordPair(encryptedCredentials []byte) error {
	f, err := os.OpenFile(databaseFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)

	defer f.Close()

	stringToWrite := string(encryptedCredentials) + "\n"
	if _, err = f.WriteString(stringToWrite); err != nil {
		return errors.New("Can't write credentials")
	}

	return f.Sync()
}