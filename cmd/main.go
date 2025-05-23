package main

import (
	"fmt"
	"idm/inner"
)

func main() {
	fmt.Print("Hello Go")
	fmt.Print(inner.RandomInt(5))
}
