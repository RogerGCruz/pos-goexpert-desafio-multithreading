package viacep

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
	return "ViaCEP"
}

type response struct {
	CEP          string `json:"cep"`
	State        string `json:"uf"`
	City         string `json:"localidade"`
	Street       string `json:"logradouro"`
	Neighborhood string `json:"bairro"`
	Error        bool   `json:"erro"`
}

func (c *Client) Lookup(ctx context.Context, cep string) (addr domain.Address, err error) {
	startedAt := time.Now()
	defer func() {
		log.Printf("[api=%s] cep=%s duration=%s err=%v", c.Name(), cep, time.Since(startedAt), err)
	}()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Address{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Address{}, err
	}
	defer resp.Body.Close()

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

	return domain.Address{
		CEP:          payload.CEP,
		State:        payload.State,
		City:         payload.City,
		Street:       payload.Street,
		Neighborhood: payload.Neighborhood,
	}, nil
}
