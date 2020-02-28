# password-storage

[![Build Status](https://travis-ci.org/andriikushch/password-storage.svg?branch=master)](https://travis-ci.org/andriikushch/password-storage)
[![Go Report Card](https://goreportcard.com/badge/github.com/andriikushch/password-storage)](https://goreportcard.com/report/github.com/andriikushch/password-storage)

Simple password storage implemented with AES encryption.

Current implementation has only command line interface. As storage it uses a local file ```~/.dat2```. 

## Usage:

```
go run main.go -h
  -a	add new account with random password
  -ac  	add new username:password
  -d	delete password for account
  -g	copy to clip board password for account
  -i	interactive mode
  -l	list of stored accounts
  -v	version
```
