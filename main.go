package main

import (
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"syscall"

	"io"

	"github.com/andriikushch/clipboard"
	"github.com/andriikushch/password-storage/repository"
	"golang.org/x/crypto/ssh/terminal"
)

var accountList map[int]string

func main() {
	var l, ac, a, g, d bool

	flag.BoolVar(&l, "l", false, "list of stored accounts")
	flag.BoolVar(&a, "a", false, "add new account with random password")
	flag.BoolVar(&ac, "ac", false, "add new username:password")
	flag.BoolVar(&g, "g", false, "copy to clip board password for account")
	flag.BoolVar(&d, "d", false, "delete password for account")

	flag.Parse()

	fmt.Print("Enter master password: ")
	masterPassword, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		panic("Can't read password input")
	}
	fmt.Println("")
	tmpKey := sha256.Sum256(masterPassword)
	key := tmpKey[:]
	switch true {
	case a:
		addAccountWithRandomPasswordMenuItem(key)
	case ac:
		addCredentialsMenuItem(key)
	case g:
		getPasswordForAccountMenuItem(key)
	case l:
		printAccountMenuItem(key)
	case d:
		deleteAccountMenuItem(key)
	}
}
func printAccountMenuItem(key []byte) {
	printAccountsList(key)
}

func deleteAccountMenuItem(key []byte) {
	accountToDelete := -1
	printAccountsList(key)
	fmt.Println("Select account to delete (number): ")
	fmt.Scanln(&accountToDelete)
	repository.DeleteCredentials(key, accountList[accountToDelete])
}

func getPasswordForAccountMenuItem(key []byte) {
	var account string
	fmt.Print("Enter account name: ")
	fmt.Scanln(&account)
	password, err := repository.FindPassword(key, account)
	if err != nil {
		fmt.Errorf("%v", err)
	} else {
		clipboard.WriteAll(password)
	}
}

func addAccountWithRandomPasswordMenuItem(key []byte) {
	var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%^&*()-_=+{}[]")
	var account string
	fmt.Print("Enter account name: ")
	fmt.Scanln(&account)

	password := randChar(16, stdChars)
	if err := repository.AddNewCredentials(key[:], []byte(password), []byte(password), account); err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Println("")
}

func addCredentialsMenuItem(key []byte) {
	var account string
	fmt.Print("Enter account name: ")
	fmt.Scanln(&account)
	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Print("\nEnter password confirmation: ")
	bytePasswordConfirmation, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Can't read password confiramtion")
	}
	if err := repository.AddNewCredentials(key[:], bytePassword, bytePasswordConfirmation, account); err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Println("")
}

func printAccountsList(key []byte) {
	generateAccountList(key[:])
	for k, v := range accountList {
		fmt.Printf("%d : %s \n", k, v)
	}
}

func generateAccountList(key []byte) {
	list, err := repository.GetAccountsList(key)

	if err != nil {
		panic("Some error" + err.Error())
	}
	accountList = make(map[int]string, len(list))
	for k, v := range list {
		accountList[k] = v
	}
}

func randChar(length int, chars []byte) string {
	newPword := make([]byte, length)
	randomData := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, randomData); err != nil {
			panic(err)
		}
		for _, c := range randomData {
			if c >= maxrb {
				continue
			}
			newPword[i] = chars[c%clen]
			i++
			if i == length {
				return string(newPword)
			}
		}
	}
	panic("unreachable")
}
