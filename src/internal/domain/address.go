package domain

type Address struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
}
