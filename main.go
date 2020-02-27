package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/andriikushch/password-storage/menu"
	"log"
	"os"
	"syscall"

	"github.com/andriikushch/password-storage/repository"
	"golang.org/x/crypto/ssh/terminal"
)

const version = "0.0.6"

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

	m := menu.NewMenu(repository.NewPasswordRepository(userHomeDir() + "/.dat2"))

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
			err = m.AddAccountWithRandomPasswordMenuItem(key)
		case ac:
			err = m.AddCredentialsMenuItem(key)
		case g:
			err = m.GetPasswordForAccountMenuItem(key)
		case l:
			m.PrintAccountMenuItem(key)
		case d:
			err = m.DeleteAccountMenuItem(key)
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
				err = m.AddAccountWithRandomPasswordMenuItem(key)
			case "ac":
				err = m.AddCredentialsMenuItem(key)
			case "g":
				err = m.GetPasswordForAccountMenuItem(key)
			case "l":
				m.PrintAccountMenuItem(key)
			case "d":
				err = m.DeleteAccountMenuItem(key)
			case "q":
				quit = true
			}

			if err != nil {
				log.Fatalln(err)
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
