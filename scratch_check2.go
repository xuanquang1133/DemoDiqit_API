package main

import (
	"fmt"
	"demodiqit_api/helpers/crypt"
)

func main() {
	hash := "$2a$10$.ldKiy4qDYYf5ij1C4SX8OuM9SrT96wkSlJxFfwWzLCGWvZFbd8ga"
	match := crypt.CheckPasswordHash("123456", hash)
	fmt.Println("Does 123456 match?", match)
}
