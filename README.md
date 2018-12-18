# password-storage

[![Build Status](https://travis-ci.org/andriikushch/password-storage.svg?branch=master)](https://travis-ci.org/andriikushch/password-storage)
[![Go Report Card](https://goreportcard.com/badge/github.com/andriikushch/password-storage)](https://goreportcard.com/report/github.com/andriikushch/password-storage)

Simple password storage implemented with AES encryption.

|            Build                                                                                                                     | Md5 sum                           | Version |  OS   |
|:-------------------------------------------------------------------------------------------------------------------------------------|:----------------------------------|:--------|:------|
|  [password-storage-mac-x64-v0.0.4(https://github.com/andriikushch/password-storage/tree/master/bin/password-storage-mac-x64-v0.0.4)] |  8caaa9171e17de46ea7b365d098a1a08 | v0.0.4  | OS X  |

Usage:

```
go run main.go -h
  -g	add new account with random password
  -ac   add new username:password
  -d	delete password for account
  -g	copy to clip board password for account
  -l	list of stored accounts
```