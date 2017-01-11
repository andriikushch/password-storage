package repository

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github.com/andriikushch/password-storage/crypt"
)

var (
	databaseFile = "dat2"
	db           = make(map[string][]byte)
)

func FindPassword(key []byte, account string) (string, error) {
	// Open a RO file
	decodeFile, err := os.Open(databaseFile)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	// Decode -- We need to pass a pointer otherwise accounts2 isn't modified
	decoder.Decode(&db)

	password, ok := db[account]

	if !ok {
		return "", errors.New("Can't find password for account")
	}

	return crypt.Decrypt(key, password), nil
}

func ShowAccountsList(key []byte) error {
	// Open a RO file
	decodeFile, err := os.Open(databaseFile)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	// Decode -- We need to pass a pointer otherwise accounts2 isn't modified
	decoder.Decode(&db)

	for acc := range db {
		fmt.Printf("%s\n", acc)
	}

	return nil
}

func AddNewCredentials(key, bytePassword, bytePasswordConfirmation []byte, account string) error {
	password := string(bytePassword)
	passwordConfirmation := string(bytePasswordConfirmation)

	if password == passwordConfirmation {
		return storeAccountPasswordPair(key, account, password)
	}

	return errors.New("Password and Password confirmation is not equal")
}

func storeAccountPasswordPair(key []byte, account string, password string) error {
	db[account] = crypt.Encrypt(key, password)
	encodeFile := new(os.File)

	if _, err := os.Stat(databaseFile); os.IsNotExist(err) {
		encodeFile, err = os.Create(databaseFile)

		if err != nil {
			panic(err)
		}
	}

	// Since this is a binary format large parts of it will be unreadable
	encoder := gob.NewEncoder(encodeFile)

	// Write to the file
	if err := encoder.Encode(db); err != nil {
		panic(err)
	}

	return encodeFile.Close()
}
