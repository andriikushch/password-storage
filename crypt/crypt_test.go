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
		c := Encrypt(key[:], p)
		d := string(Decrypt(key[:], c))
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

	if string(Encrypt(key[:], text)) == string(Encrypt(key[:], text)) {
		fmt.Printf("%v", "encrypted text is equal")
		t.FailNow()
	}
}

func TestDecryptPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	Encrypt([]byte(""), "facebook")
}
