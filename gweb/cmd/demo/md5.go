package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	//pass := "admin"
	//salt := fmt.Sprintf("%x", md5.Sum([]byte(uuid.NewV4().Bytes())))
	//password := fmt.Sprintf("%x", md5.Sum([]byte(pass + salt)))
	//
	//fmt.Printf("salt: %s\n", salt)
	//fmt.Printf("password: %s\n", password)

	// verify
	ts := "b0d2c14f091a9f895392bb67aca06ba6"
	tp := "64ddd4e4850955cb4322dc06b316253e"
	pa := fmt.Sprintf("%x", md5.Sum([]byte("admin"+ts)))
	if pa == tp {
		fmt.Printf("verify password\n")
	} else {
		fmt.Printf("pa: %s", pa)
		fmt.Printf("tp: %s", tp)
	}
}
