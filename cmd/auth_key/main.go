package main

import (
	"fmt"
	"log"

	"github.com/PsionicAlch/course-platform/internal/authentication"
)

func main() {
	key, err := authentication.GenerateKeyString()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(key)
}
