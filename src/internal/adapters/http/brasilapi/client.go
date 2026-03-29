package brasilapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
}

func (c *Client) Lookup(ctx context.Context, cep string) (domain.Address, error) {
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

	if resp.StatusCode != http.StatusOK {
		return domain.Address{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var payload response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return domain.Address{}, err
	}

	return domain.Address{
		CEP:          payload.CEP,
		State:        payload.State,
		City:         payload.City,
		Street:       payload.Street,
		Neighborhood: payload.Neighborhood,
	}, nil
}
