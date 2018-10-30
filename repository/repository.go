package repository

import (
	"encoding/gob"
	"errors"
	"os"

	"encoding/base64"
	"runtime"

	"github.com/andriikushch/password-storage/crypt"
)

var (
	databaseFile                               = userHomeDir() + "/.dat2"
	db                                         = make(map[string][]byte)
	ErrOpenDatabase                            = errors.New("error open database")
	ErrDecodeString                            = errors.New("error decode string")
	ErrPasswordNotFound                        = errors.New("error password not found")
	ErrPasswordAndPasswordConfirmationMatching = errors.New("error password and password confirmation not matching")
	ErrFileCreation                            = errors.New("error create new file")
	ErrEncodeDB                                = errors.New("error encode db")
	ErrDecodeDB                                = errors.New("error decode db")
)

type Repository interface {
	FindPassword(key []byte, account string) (string, error)
	GetAccountsList(key []byte) ([]string, error)
	AddNewCredentials(key, bytePassword, bytePasswordConfirmation []byte, account string) error
	DeleteCredentials(key []byte, account string) error
}

func FindPassword(key []byte, account string) (string, error) {
	// Open a RO file
	decodeFile, err := os.Open(databaseFile)
	if err != nil {
		return "", ErrOpenDatabase
	}
	defer decodeFile.Close()

	decoder := gob.NewDecoder(decodeFile)
	decoder.Decode(&db)

	for acc, password := range db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			return "", ErrDecodeString
		}

		decryptedAccount := crypt.Decrypt(key, encryptedAccount)
		if account == decryptedAccount {
			return crypt.Decrypt(key, password), nil
		}
	}

	return "", ErrPasswordNotFound
}

func GetAccountsList(key []byte) ([]string, error) {
	if err := loadDB(); err != nil {
		return nil, err
	}

	var result []string

	for acc := range db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			return nil, ErrDecodeString
		}
		result = append(result, crypt.Decrypt(key, encryptedAccount))
	}

	return result, nil
}

func AddNewCredentials(key, bytePassword, bytePasswordConfirmation []byte, account string) error {
	password := string(bytePassword)
	passwordConfirmation := string(bytePasswordConfirmation)

	if password == passwordConfirmation {
		return storeAccountPasswordPair(key, account, password)
	}

	return ErrPasswordAndPasswordConfirmationMatching
}

func DeleteCredentials(key []byte, account string) error {
	if err := loadDB(); err != nil {
		return err
	}

	for acc := range db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			return ErrDecodeString
		}

		if account == crypt.Decrypt(key, encryptedAccount) {
			delete(db, acc)
			if err := writeToFile(); err != nil {
				return err
			}
		}
	}

	return nil
}

func storeAccountPasswordPair(key []byte, account string, password string) error {
	loadDB()

	encryptedAccount := crypt.Encrypt(key, account)

	for acc := range db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			return ErrDecodeString
		}

		if account == crypt.Decrypt(key, encryptedAccount) {
			delete(db, acc)
		}
	}

	db[base64.StdEncoding.EncodeToString(encryptedAccount)] = crypt.Encrypt(key, password)

	writeToFile()

	return nil
}

func writeToFile() error {
	encodeFile := new(os.File)

	//recreate DB file
	encodeFile, err := os.Create(databaseFile)
	if err != nil {
		return ErrFileCreation
	}

	defer encodeFile.Close()
	// Since this is a binary format large parts of it will be unreadable
	encoder := gob.NewEncoder(encodeFile)
	// Write to the file
	if err := encoder.Encode(db); err != nil {
		return ErrEncodeDB
	}

	return nil
}

func loadDB() error {
	// Open a RO file
	decodeFile, err := os.Open(databaseFile)
	if err != nil {
		return ErrOpenDatabase
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	if err := decoder.Decode(&db); err != nil {
		return ErrDecodeDB
	}

	return nil
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
