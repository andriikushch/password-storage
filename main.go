package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/andriikushch/clipboard"
	"github.com/andriikushch/password-storage/repository"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	masterPassword = "example key 1234"
)

func main() {
	key := sha256.Sum256([]byte(masterPassword))

	command := flag.String("command", "add-new-credentials", "a string")
	flag.Parse()

	switch *command {
	case "add-new-credentials":
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
	case "load-password":
		var account string
		fmt.Print("Enter account name: ")
		fmt.Scanln(&account)
		password, err := repository.FindPassword(key[:], account)
		if err != nil {
			fmt.Errorf("%v", err)
		} else {
			clipboard.WriteAll(password)
		}
	case "accounts":
		repository.ShowAccountsList(key[:])
	}
}
