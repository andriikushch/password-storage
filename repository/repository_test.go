package repository

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
)

func TestAddNewCredentials(t *testing.T) {
	masterPassword := "test123"
	key := sha256.Sum256([]byte(masterPassword))
	p1 := "1"
	p2 := "!@#$%^&*()_"
	p3 := "123sad123@#!@# ADSASDA"
	p4 := "123sad123@#!@# ADSASDA___((()))"
	p5 := "t5plil"
	p6 := "a"

	a1 := "facebook"
	a2 := "gmail"
	a3 := "mail"
	a4 := "bank"
	a5 := "atm"
	a6 := "twitter"

	passwords := []string{p1, p2, p3, p4, p5, p6}
	accounts := []string{a1, a2, a3, a4, a5, a6}

	databaseFile = "db2_test"
	defer os.Remove(databaseFile)

	for i := range passwords {
		AddNewCredentials(key[:], []byte(passwords[i]), []byte(passwords[i]), accounts[i])
	}
	//to found bug with duplication in map
	for i := range passwords {
		AddNewCredentials(key[:], []byte(passwords[i]), []byte(passwords[i]), accounts[i])
	}

	list, err := getAccountsList(key[:])

	if err != nil {
		fmt.Printf("%v", err)
		t.FailNow()
	}

	for i := range list {
		for j := range accounts {
			if accounts[j] != list[i] {
				break
			}

			if j == len(accounts) {
				fmt.Print("Credentials in not stored")
				t.FailNow()
			}
		}
	}

	if len(accounts) != len(list) {
		fmt.Print("Accounts duplicated")
		t.FailNow()
	}
}
