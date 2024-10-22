package main

import (
	"fmt"
	"log"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

func main() {
	hashSlice, err := utils.RandomByteSlice(64)
	if err != nil {
		log.Fatal(err)
	}

	blockSlice, err := utils.RandomByteSlice(32)
	if err != nil {
		log.Fatal(err)
	}

	hashKey := utils.EncodeString(hashSlice)
	blockKey := utils.EncodeString(blockSlice)

	fmt.Printf("Hash Key: %s\nBlock Key: %s\n", hashKey, blockKey)
}
