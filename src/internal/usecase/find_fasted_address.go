package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/domain"
	"github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/ports"
)

// REQUISITO 4: Timeout em 1 segundo
var ErrTimeout = errors.New("timeout for all providers in the given time one second")

type FindFastestAddressUseCase struct {
	providers []ports.CEPProvider
	timeout   time.Duration
}

func NewFindFastestAddress(providers []ports.CEPProvider, timeout time.Duration) *FindFastestAddressUseCase {
	return &FindFastestAddressUseCase{
		providers: providers,
		timeout:   timeout,
	}
}

type result struct {
	providerName domain.Address
	source       string
	err          error
}

func (uc *FindFastestAddressUseCase) Execute(ctx context.Context, cep string) (domain.Address, string, error) {
	// REQUISITO 4: Criando contexto com timeout de 1 segundo
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	ch := make(chan result, len(uc.providers))

	// REQUISITO 1: Requisições simultâneas cada provider executado em uma goroutine separada
	for _, provider := range uc.providers {
		p := provider
		go func() {
			addr, err := p.Lookup(ctx, cep)
			ch <- result{
				providerName: addr,
				source:       p.Name(),
				err:          err,
			}
		}()
	}

	var errs []error
	// REQUISITO 2: Race Condition (Corrida)
	// O select aguarda o primeiro resultado bem-sucedido que chegar no canal
	// Assim que uma API responde com sucesso, a função retorna imediatamente,
	// descartando a resposta da outra API (mais lenta)
	for i := 0; i < len(uc.providers); i++ {
		select {
		case res := <-ch:
			if res.err == nil {
				return res.providerName, res.source, nil
			}
			errs = append(errs, fmt.Errorf("%s: %w", res.source, res.err))
		// REQUISITO 4: Se nenhuma API responder em até 1 segundo, retorna erro de timeout
		case <-ctx.Done():
			return domain.Address{}, "", ErrTimeout
		}
	}
	return domain.Address{}, "", errors.Join(errs...)

}
