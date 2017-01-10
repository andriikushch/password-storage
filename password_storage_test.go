package main

import (
	"testing"
	"crypto/sha256"
	"fmt"
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

	passwords := []string {p1, p2, p3, p4, p5, p6}

	for _, p := range passwords {
		c := encrypt(key[:], p)
		d := string(decrypt(key[:], c))
		if d != p {
			fmt.Printf("%v \n %v \n", []byte(d), []byte(p))
			fmt.Printf("%s \n", d + " != " + p)
			t.Fail()
		}
	}
}
