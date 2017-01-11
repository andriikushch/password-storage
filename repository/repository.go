package repository

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github.com/andriikushch/password-storage/crypt"
	"encoding/base64"
	"runtime"
)

var (
	databaseFile = userHomeDir() + "/.dat2"
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

	for acc, password := range db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			return "", err
		}

		if account == crypt.Decrypt(key, encryptedAccount) {
			return crypt.Decrypt(key, password), nil
		}
	}

	return "", errors.New("Can't find password for account")
}

func getAccountsList(key []byte) ([]string, error) {
	loadDB()

	var result []string

	for acc := range db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			return nil, err
		}
		result = append(result, crypt.Decrypt(key, encryptedAccount))
	}

	return result, nil
}

func PrintAccountsList(key []byte) error {
	list, err := getAccountsList(key)

	if err != nil {
		return err
	}
	for _, v := range list {
		fmt.Printf("%s\n", v)
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
	loadDB()

	encryptedAccount := crypt.Encrypt(key, account)

	for acc := range db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			return err
		}

		if account == crypt.Decrypt(key, encryptedAccount) {
			delete(db, acc)
		}
	}


	db[base64.StdEncoding.EncodeToString(encryptedAccount)] = crypt.Encrypt(key, password)
	encodeFile := new(os.File)

	//recreate DB file
	encodeFile, err := os.Create(databaseFile)

	if err != nil {
		panic(err)
	}

	// Since this is a binary format large parts of it will be unreadable
	encoder := gob.NewEncoder(encodeFile)

	// Write to the file
	if err := encoder.Encode(db); err != nil {
		panic(err)
	}

	return encodeFile.Close()
}

func loadDB() error {
	// Open a RO file
	decodeFile, err := os.Open(databaseFile)
	if err != nil {
		return err
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	// Decode -- We need to pass a pointer otherwise accounts2 isn't modified
	return  decoder.Decode(&db)
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}