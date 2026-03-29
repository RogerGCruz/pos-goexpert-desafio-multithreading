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

	// REQUISITO 4: Cliente HTTP com timeout de 1 segundo
	httpClient := &http.Client{
		Timeout: 1 * time.Second,
	}

	// REQUISITO 1: Instanciando dois providers (BrasilAPI e ViaCEP) que farão requisições simultâneas
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

	// REQUISITO 3: Output: Qual API entregou a resposta: BrasilAPI ou ViaCEP
	output := struct {
		API     string      `json:"api"`
		Address interface{} `json:"address"`
	}{
		API:     source,
		Address: address,
	}

	// REQUISITO 3: Exibindo o resultado no terminal (stdout) em formato JSON
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
