package repository

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
)

func TestAddNewCredentials(t *testing.T) {
	masterPassword := "test123"
	tmpKey := sha256.Sum256([]byte(masterPassword))
	key := tmpKey[:]
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

	databaseFile = "/tmp/db2_test"
	defer os.Remove(databaseFile)

	for i := range passwords {
		if err := AddNewCredentials(key, []byte(passwords[i]), []byte(passwords[i]), accounts[i]); err != nil {
			fmt.Println("AddNewCredentials 1")
			fmt.Println(err.Error())
			t.FailNow()
		}
	}

	//to found bug with duplication in map
	for i := range passwords {
		if err := AddNewCredentials(key, []byte(passwords[i]), []byte(passwords[i]), accounts[i]); err != nil {
			fmt.Println("AddNewCredentials 2")
			fmt.Println(err.Error())
			t.FailNow()
		}
	}

	list, err := GetAccountsList(key)

	if err != nil {
		fmt.Println("GetAccountsList")
		fmt.Println(err.Error())
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

	password, err := FindPassword(key, a1)

	if err != nil {
		fmt.Print(err.Error())
		t.FailNow()
	}

	if password != p1 {
		fmt.Print("Passwords should match")
		t.FailNow()
	}

	if err := DeleteCredentials(key, a1); err != nil {
		fmt.Printf("%v", err)
		t.FailNow()
	}

	_, err = FindPassword(key, a1)

	if err.Error() != ErrPasswordNotFound.Error() {
		fmt.Printf("%v", err)
		t.FailNow()
	}
}

func TestChangeMasterKey(t *testing.T) {
	masterPassword := "test123"
	tmpKey := sha256.Sum256([]byte(masterPassword))
	key := tmpKey[:]
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

	databaseFile = "/tmp/db2_test"
	defer os.Remove(databaseFile)

	for i := range passwords {
		if err := AddNewCredentials(key, []byte(passwords[i]), []byte(passwords[i]), accounts[i]); err != nil {
			fmt.Println("AddNewCredentials 1")
			fmt.Println(err.Error())
			t.FailNow()
		}
	}

	// new key
	newMasterPassword := "123test"
	newTmpKey := sha256.Sum256([]byte(newMasterPassword))
	newKey := newTmpKey[:]

	err := ChangeMasterKey(key, newKey)
	if err != nil {
		fmt.Println("ChangeMasterKey")
		fmt.Println(err.Error())
		t.FailNow()
	}

	list, err := GetAccountsList(newKey)

	// assertions
	if err != nil {
		fmt.Println("GetAccountsList")
		fmt.Println(err.Error())
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

	password, err := FindPassword(newKey, a1)

	if err != nil {
		fmt.Print(err.Error())
		t.FailNow()
	}

	if password != p1 {
		fmt.Print("Passwords should match")
		t.FailNow()
	}
}
