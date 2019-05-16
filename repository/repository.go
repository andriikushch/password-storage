package repository

import (
	"encoding/gob"
	"errors"
	"log"
	"os"

	"encoding/base64"
	"github.com/andriikushch/password-storage/crypt"
)

var (
	ErrOpenDatabase                            = errors.New("error open database")
	ErrDecodeString                            = errors.New("error decode string")
	ErrPasswordNotFound                        = errors.New("error password not found")
	ErrPasswordAndPasswordConfirmationMatching = errors.New("error password and password confirmation not matching")
	ErrFileCreation                            = errors.New("error create new file")
	ErrEncodeDB                                = errors.New("error encode db")
	ErrDecodeDB                                = errors.New("error decode db")
)

// interface for credential repository
type Repository interface {
	FindPassword(key []byte, account string) (string, error)
	GetAccountsList(key []byte) ([]string, error)
	AddNewCredentials(key, bytePassword, bytePasswordConfirmation []byte, account string) error
	DeleteCredentials(key []byte, account string) error
}

func NewPasswordRepository(dbFile string) *PasswordRepository {
	return &PasswordRepository{
		db:     make(map[string][]byte),
		dbFile: dbFile,
	}
}

type PasswordRepository struct {
	db     map[string][]byte
	dbFile string
}

// Load encrypted password
func (p *PasswordRepository) FindPassword(key []byte, account string) (string, error) {
	// Open a RO file
	decodeFile, err := os.Open(p.dbFile)
	if err != nil {
		log.Println(err.Error())
		return "", ErrOpenDatabase
	}
	defer func() {
		err := decodeFile.Close()
		log.Println(err.Error())
	}()

	decoder := gob.NewDecoder(decodeFile)
	decoder.Decode(&p.db)

	for acc, password := range p.db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			log.Println(err.Error())
			return "", ErrDecodeString
		}

		decryptedAccount, err := crypt.Decrypt(key, encryptedAccount)

		if err != nil {
			return "", err
		}

		if account == decryptedAccount {
			return crypt.Decrypt(key, password)
		}
	}

	return "", ErrPasswordNotFound
}

// Returns array of accounts
func (p *PasswordRepository) GetAccountsList(key []byte) ([]string, error) {
	if err := p.loadDB(); err != nil {
		return nil, err
	}

	var result []string

	for acc := range p.db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			return nil, ErrDecodeString
		}
		decryptedAccount, err := crypt.Decrypt(key, encryptedAccount)
		if err != nil {
			return nil, err
		}
		result = append(result, decryptedAccount)
	}

	return result, nil
}

// Add store new credentials
func (p *PasswordRepository) AddNewCredentials(key, bytePassword, bytePasswordConfirmation []byte, account string) error {
	password := string(bytePassword)
	passwordConfirmation := string(bytePasswordConfirmation)

	if password == passwordConfirmation {
		return p.storeAccountPasswordPair(key, account, password)
	}

	return ErrPasswordAndPasswordConfirmationMatching
}

// Remove credential
func (p *PasswordRepository) DeleteCredentials(key []byte, account string) error {
	if err := p.loadDB(); err != nil {
		return err
	}

	for acc := range p.db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			log.Println(err.Error())
			return ErrDecodeString
		}

		decryptedAccount, err := crypt.Decrypt(key, encryptedAccount)
		if err != nil {
			return err
		}

		if account == decryptedAccount {
			delete(p.db, acc)
			if err := p.writeToFile(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Re-encrypt all credentials with new master password
func (p *PasswordRepository) ChangeMasterKey(oldKey, newKey []byte) error {
	newDb := make(map[string][]byte)
	accounts, err := p.GetAccountsList(oldKey)

	if err != nil {
		return err
	}

	for _, account := range accounts {
		passwd, err := p.FindPassword(oldKey, account)

		if err != nil {
			return err
		}

		encryptedAccount, err := crypt.Encrypt(newKey, account)
		if err != nil {
			return err
		}
		encryptedPassword, err := crypt.Encrypt(newKey, passwd)
		if err != nil {
			return err
		}
		newDb[base64.StdEncoding.EncodeToString(encryptedAccount)] = encryptedPassword
	}

	p.db = newDb

	return p.writeToFile()
}

func (p *PasswordRepository) storeAccountPasswordPair(key []byte, account string, password string) error {
	err := p.loadDB()

	if err != nil {
		return err
	}

	encryptedAccount, err := crypt.Encrypt(key, account)

	if err != nil {
		return err
	}

	for acc := range p.db {
		encryptedAccount, err := base64.StdEncoding.DecodeString(acc)
		if err != nil {
			log.Println(err.Error())
			return ErrDecodeString
		}

		decryptedAccount, err := crypt.Decrypt(key, encryptedAccount)

		if err != nil {
			return err
		}
		if account == decryptedAccount {
			delete(p.db, acc)
		}
	}

	encryptedPassword, err := crypt.Encrypt(key, password)
	if err != nil {
		return err
	}
	p.db[base64.StdEncoding.EncodeToString(encryptedAccount)] = encryptedPassword

	return p.writeToFile()
}

func (p *PasswordRepository) writeToFile() error {
	encodeFile := new(os.File)

	//recreate DB file
	encodeFile, err := os.Create(p.dbFile)
	if err != nil {
		log.Println(err.Error())
		return ErrFileCreation
	}

	defer func() {
		err := encodeFile.Close()
		log.Println(err.Error())
	}()

	encoder := gob.NewEncoder(encodeFile)
	// Write to the file
	if err := encoder.Encode(p.db); err != nil {
		log.Println(err.Error())
		return ErrEncodeDB
	}

	return nil
}

func (p *PasswordRepository) loadDB() error {
	var decodeFile *os.File
	var err error
	if _, err := os.Stat(p.dbFile); os.IsNotExist(err) {
		decodeFile, err = os.Create(p.dbFile)
	} else {
		// Open a RO file
		decodeFile, err = os.Open(p.dbFile)
	}

	if err != nil {
		log.Println(err.Error())
		return ErrOpenDatabase
	}

	defer func() {
		err := decodeFile.Close()
		log.Println(err.Error())
	}()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	fi, err := decodeFile.Stat()
	if err != nil {
		log.Println(err.Error())
		return ErrDecodeDB
	}

	if fi.Size() == 0 {
		return nil
	}

	if err := decoder.Decode(&p.db); err != nil {
		log.Println(err.Error())
		return ErrDecodeDB
	}

	return nil
}
