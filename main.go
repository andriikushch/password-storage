package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/andriikushch/clipboard"
	"github.com/andriikushch/password-storage/repository"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

var masterPassword string

func main() {
	var l, a, g bool
	key := sha256.Sum256([]byte(masterPassword))

	flag.BoolVar(&l, "l", false, "list of stored accounts")
	flag.BoolVar(&a, "a", false, "add new username:password")
	flag.BoolVar(&g, "g", false, "copy to clip board password for account")

	flag.Parse()

	fmt.Print("Enter master password: ")
	fmt.Scanln(&masterPassword)
	fmt.Println("")

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
		repository.PrintAccountsList(key[:])
	}
}
