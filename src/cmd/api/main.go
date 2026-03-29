package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run .src/cmd/api/main.go <CEP>")
		os.Exit(1)
	}

	cep := os.Args[1]
	if !isValidCEP(cep) {
		fmt.Printf("Invalid CEP format: %s\n", cep)
		os.Exit(1)
	}
}

func isValidCEP(cep string) bool {
	ok, _ := regexp.MatchString("^[0-9]{8}$", cep)
	return ok
}
