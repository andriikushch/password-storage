package main

import (
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/atotto/clipboard"

	"io"

	"github.com/andriikushch/password-storage/repository"
	"golang.org/x/crypto/ssh/terminal"
)

const version = "0.0.6"

var r = repository.NewPasswordRepository(userHomeDir() + "/.dat2")

func main() {
	var l, ac, a, g, d, v, i bool

	flag.BoolVar(&l, "l", false, "list of stored accounts")
	flag.BoolVar(&a, "a", false, "add new account with random password")
	flag.BoolVar(&ac, "ac", false, "add new username:password")
	flag.BoolVar(&g, "g", false, "copy to clip board password for account")
	flag.BoolVar(&d, "d", false, "delete password for account")
	flag.BoolVar(&v, "v", false, "version")
	flag.BoolVar(&i, "i", false, "interactive mode")

	flag.Parse()

	if v {
		fmt.Println(version)
		return
	}

	fmt.Print("Enter master password: ")
	masterPassword, err := terminal.ReadPassword(syscall.Stdin)
	tmpKey := sha256.Sum256(masterPassword)
	key := tmpKey[:]

	if err != nil {
		panic("Can't read password input")
	}
	fmt.Println("")
	if !i {
		var err error
		switch true {
		case a:
			err = addAccountWithRandomPasswordMenuItem(key)
		case ac:
			err = addCredentialsMenuItem(key)
		case g:
			err = getPasswordForAccountMenuItem(key)
		case l:
			printAccountMenuItem(key)
		case d:
			err = deleteAccountMenuItem(key)
		}

		if err != nil {
			log.Fatalln(err)
		}
	} else {
		var quit bool
		for !quit {
			var err error
			fmt.Println("Print command: ")
			var command string
			_, err = fmt.Scanln(&command)

			if err != nil {
				panic("Can't read command input")
			}

			switch command {
			case "a":
				err = addAccountWithRandomPasswordMenuItem(key)
			case "ac":
				err = addCredentialsMenuItem(key)
			case "g":
				err = getPasswordForAccountMenuItem(key)
			case "l":
				printAccountMenuItem(key)
			case "d":
				err = deleteAccountMenuItem(key)
			case "q":
				quit = true
			}

			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
func printAccountMenuItem(key []byte) {
	printAccountsList(key)
}

func deleteAccountMenuItem(key []byte) error {
	accountToDelete := -1
	printAccountsList(key)
	fmt.Println("Select account to delete (number): ")
	_, err := fmt.Scanln(&accountToDelete)
	if err != nil {
		return err
	}
	return r.DeleteCredentials(key, generateAccountList(key[:])[accountToDelete])
}

func getPasswordForAccountMenuItem(key []byte) error {
	var account string
	fmt.Print("Enter account name: ")
	_, err := fmt.Scanln(&account)
	if err != nil {
		return err
	}
	password, err := r.FindPassword(key, account)
	if err != nil {
		return err
	} else {
		return clipboard.WriteAll(password)
	}
}

func addAccountWithRandomPasswordMenuItem(key []byte) error {
	var account string
	fmt.Print("Enter account name: ")
	_, err := fmt.Scanln(&account)
	if err != nil {
		return err
	}

	password := randChar(16)
	return r.AddNewCredentials(key[:], []byte(password), []byte(password), account)
}

func addCredentialsMenuItem(key []byte) error {
	var account string
	fmt.Print("Enter account name: ")
	_, err := fmt.Scanln(&account)
	if err != nil {
		return err
	}
	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Print("\nEnter password confirmation: ")
	bytePasswordConfirmation, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Can't read password confirmation")
		return err
	}

	return r.AddNewCredentials(key[:], bytePassword, bytePasswordConfirmation, account)
}

func printAccountsList(key []byte) {

	for k, v := range generateAccountList(key[:]) {
		fmt.Printf("%d : %s \n", k, v)
	}
}

func generateAccountList(key []byte) []string {
	list, err := r.GetAccountsList(key)

	if err != nil {
		panic("Some error" + err.Error())
	}

	return list
}

func randChar(length int) string {
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%^&*()-_=+{}[]")
	password := make([]byte, length)
	randomData := make([]byte, length+(length/4)) // storage for random bytes.
	charsLength := byte(len(chars))
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
			password[i] = chars[c%charsLength]
			i++
			if i == length {
				return string(password)
			}
		}
	}
}

func userHomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	panic("$HOME is not set")
}
