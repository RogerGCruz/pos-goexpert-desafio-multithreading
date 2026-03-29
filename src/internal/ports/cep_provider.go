package ports

import (
	"context"

	"github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/domain"
)

type CEPProvider interface {
	Name() string
	Lookup(ctx context.Context, cep string) (domain.Address, error)
}
