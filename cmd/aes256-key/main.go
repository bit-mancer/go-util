package main

import (
	"fmt"

	"github.com/bit-mancer/go-util/crypto"
)

func main() {
	key := crypto.NewRandomAESKey()
	fmt.Println(key.ToBase64())
}
