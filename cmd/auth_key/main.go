package main

import (
	"fmt"
	"log"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
)

func main() {
	key, err := authentication.GenerateKeyString()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(key)
}
