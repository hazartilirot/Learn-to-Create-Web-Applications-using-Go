package main

import (
	"fmt"
	"github.com/username/project-name/hash"
)

func main() {
	hmac := hash.NewHMAC("my-secret-key")
	sha := hmac.Hash("this-is-my-string-to-hash")
	fmt.Println(sha)
}
