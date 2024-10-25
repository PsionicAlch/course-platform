package main

import (
	"fmt"
	"log"

	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
)

func main() {
	key, err := gatekeeper.GenerateGatekeeperKey()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Gatekeeper Key: %s\n", key)
}
