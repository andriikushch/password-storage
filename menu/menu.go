package menu

import (
	"crypto/rand"
	"fmt"
	"github.com/andriikushch/password-storage/repository"
	"github.com/atotto/clipboard"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"syscall"
)

type Menu struct {
	Repository repository.Repository
}

func NewMenu(repository repository.Repository) *Menu {
	return &Menu{Repository: repository}
}

func (m *Menu) PrintAccountMenuItem(key []byte) {
	m.printAccountsList(key)
}

func (m *Menu) DeleteAccountMenuItem(key []byte) error {
	accountToDelete := -1
	m.printAccountsList(key)
	fmt.Println("Select account to delete (number): ")
	_, err := fmt.Scanln(&accountToDelete)
	if err != nil {
		return err
	}
	return m.Repository.DeleteCredentials(key, m.generateAccountList(key[:])[accountToDelete])
}

func (m *Menu) GetPasswordForAccountMenuItem(key []byte) error {
	var account string
	fmt.Print("Enter account name: ")
	_, err := fmt.Scanln(&account)
	if err != nil {
		return err
	}
	password, err := m.Repository.FindPassword(key, account)
	if err != nil {
		return err
	} else {
		return clipboard.WriteAll(password)
	}
}

func (m *Menu) AddAccountWithRandomPasswordMenuItem(key []byte) error {
	var account string
	fmt.Print("Enter account name: ")
	_, err := fmt.Scanln(&account)
	if err != nil {
		return err
	}

	password := m.randChar(16)
	return m.Repository.AddNewCredentials(key[:], []byte(password), []byte(password), account)
}

func (m *Menu) AddCredentialsMenuItem(key []byte) error {
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

	return m.Repository.AddNewCredentials(key[:], bytePassword, bytePasswordConfirmation, account)
}

func (m *Menu) printAccountsList(key []byte) {

	for k, v := range m.generateAccountList(key[:]) {
		fmt.Printf("%d : %s \n", k, v)
	}
}

func (m *Menu) generateAccountList(key []byte) []string {
	list, err := m.Repository.GetAccountsList(key)

	if err != nil {
		panic("Some error" + err.Error())
	}

	return list
}

func (m *Menu) randChar(length int) string {
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
