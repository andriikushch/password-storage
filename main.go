package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/andriikushch/clipboard"
	"github.com/andriikushch/password-storage/repository"
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
		if err := repository.AddNewCredentials(key[:]); err != nil {
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
