# password-storage

[![Build Status](https://travis-ci.org/andriikushch/password-storage.svg?branch=master)](https://travis-ci.org/andriikushch/password-storage)
[![Go Report Card](https://goreportcard.com/badge/github.com/andriikushch/password-storage)](https://goreportcard.com/report/github.com/andriikushch/password-storage)

Simple password storage implemented with using AES encryption.

Usage:

```
go run main.go -h
  -g	add new account with random password
  -ac   add new username:password
  -d	delete password for account
  -g	copy to clip board password for account
  -l	list of stored accounts
```