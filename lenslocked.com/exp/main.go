package main

import (
	"fmt"
	"github.com/username/project-name/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}
