package brasilapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/domain"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

func (c *Client) Name() string {
	return "BrasilAPI"
}

type response struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	Error        bool   `json:"error"`
}

// REQUISITO 1: Método que participa das requisições simultâneas
// Será chamado em uma goroutine separada no usecase
func (c *Client) Lookup(ctx context.Context, cep string) (addr domain.Address, err error) {
	startedAt := time.Now()
	defer func() {
		log.Printf("[api=%s] cep=%s duration=%s err=%v", c.Name(), cep, time.Since(startedAt), err)
	}()

	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Address{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Address{}, err
	}
	defer resp.Body.Close()

	// REQUISITO 4: Utilizando context com timeout para respeitar o limite de 1 segundo
	if resp.StatusCode != http.StatusOK {
		return domain.Address{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var payload response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return domain.Address{}, err
	}

	if payload.Error {
		return domain.Address{}, fmt.Errorf("CEP not found")
	}

	// REQUISITO 3: Retornando os dados do endereço e o nome da API para identificar qual foi a mais rápida
	return domain.Address{
		CEP:          payload.CEP,
		State:        payload.State,
		City:         payload.City,
		Street:       payload.Street,
		Neighborhood: payload.Neighborhood,
	}, nil
}
