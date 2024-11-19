package main

import (
	"fmt"
	"log"

	"github.com/PsionicAlch/psionicalch-home/website/content"
)

func main() {
	key, err := content.GenerateFileKey()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(key)
}
