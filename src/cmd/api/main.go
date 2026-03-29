package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	brasilapi "github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/adapters/http/brasilapi"
	"github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/adapters/http/viacep"
	"github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/ports"
	"github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/usecase"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run ./src/cmd/api <cep>")
		os.Exit(1)
	}

	cep := os.Args[1]
	if !isValidCEP(cep) {
		fmt.Printf("Invalid CEP: %s\n", cep)
		os.Exit(1)
	}

	httpClient := &http.Client{
		Timeout: 1 * time.Second,
	}

	providers := []ports.CEPProvider{
		brasilapi.NewClient(httpClient),
		viacep.NewClient(httpClient),
	}

	findFastestAddress := usecase.NewFindFastestAddress(providers, 1*time.Second)

	address, source, err := findFastestAddress.Execute(context.Background(), cep)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	output := struct {
		API     string      `json:"api"`
		Address interface{} `json:"address"`
	}{
		API:     source,
		Address: address,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("error generating output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}

func isValidCEP(cep string) bool {
	ok, _ := regexp.MatchString("^[0-9]{8}$", cep)
	return ok
}
