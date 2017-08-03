package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"syscall"

	"github.com/andriikushch/clipboard"
	"github.com/andriikushch/password-storage/repository"
	"golang.org/x/crypto/ssh/terminal"
)

var accountList map[int]string

func main() {
	var l, a, g, d bool

	flag.BoolVar(&l, "l", false, "list of stored accounts")
	flag.BoolVar(&a, "a", false, "add new username:password")
	flag.BoolVar(&g, "g", false, "copy to clip board password for account")
	flag.BoolVar(&d, "d", false, "delete password for account")

	flag.Parse()

	fmt.Print("Enter master password: ")
	masterPassword, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		panic("Can't read password input")
	}
	fmt.Println("")
	key := sha256.Sum256(masterPassword)
	switch true {
	case a:
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
	case g:
		var account string
		fmt.Print("Enter account name: ")
		fmt.Scanln(&account)
		password, err := repository.FindPassword(key[:], account)
		if err != nil {
			fmt.Errorf("%v", err)
		} else {
			clipboard.WriteAll(password)
		}
	case l:
		printAccountsList(key[:])
	case d:
		accountToDelete := -1
		printAccountsList(key[:])
		fmt.Println("Select account to delete (number): ")
		fmt.Scanln(&accountToDelete)
		repository.DeleteCredentials(key[:], accountList[accountToDelete])
	}
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
