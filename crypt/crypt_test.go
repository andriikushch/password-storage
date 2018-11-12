package crypt

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	masterPassword := "test123"
	key := sha256.Sum256([]byte(masterPassword))
	p1 := "1"
	p2 := "!@#$%^&*()_"
	p3 := "123sad123@#!@# ADSASDA"
	p4 := "123sad123@#!@# ADSASDA___((()))"
	p5 := "t5plil"
	p6 := "a"

	passwords := []string{p1, p2, p3, p4, p5, p6}

	for _, p := range passwords {
		c, err1 := Encrypt(key[:], p)
		if err1 != nil {
			fmt.Println(err1.Error())
			t.Fail()
		}

		d, err2 := Decrypt(key[:], c)
		if err2 != nil {
			fmt.Println(err2.Error())
			t.Fail()
		}
		if d != p {
			fmt.Printf("%v \n %v \n", []byte(d), []byte(p))
			fmt.Printf("%s \n", d+" != "+p)
			t.Fail()
		}
	}
}

func TestEncrypt(t *testing.T) {
	masterPassword := "test123"
	key := sha256.Sum256([]byte(masterPassword))
	text := "facebook"

	a, err1 := Encrypt(key[:], text)

	b, err2 := Encrypt(key[:], text)

	if err1 != nil {
		fmt.Println(err1.Error())
		t.FailNow()
	}

	if err2 != nil {
		fmt.Println(err2.Error())
		t.FailNow()
	}

	if string(a) == string(b) {
		fmt.Println("encrypted text is equal")
		t.FailNow()
	}
}

func TestDecryptError(t *testing.T) {
	_, err := Encrypt([]byte(""), "facebook")

	if err != ErrCipherCreation {
		fmt.Println("Wrong error ", err.Error())
		t.FailNow()
	}
}
