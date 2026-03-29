package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/domain"
	"github.com/rogergcruz/pos-goexpert-desafio-multithreading/src/internal/ports"
)

// fakeProvider implementa a interface CEPProvider para testes
type fakeProvider struct {
	name    string
	delay   time.Duration
	address domain.Address
	err     error
}

func (f fakeProvider) Name() string {
	return f.name
}

func (f fakeProvider) Lookup(ctx context.Context, cep string) (domain.Address, error) {
	select {
	case <-ctx.Done():
		return domain.Address{}, ctx.Err()
	case <-time.After(f.delay):
		return f.address, f.err
	}
}

// TestExecuteReturnsFastestProvider: REQUISITO 2 (Race Condition)
// Verifica que o sistema retorna apenas a resposta da API mais rápida
// e descarta a resposta da mais lenta
func TestExecuteReturnsFastestProvider(t *testing.T) {
	// Provider mais lento (200ms)
	slow := fakeProvider{
		name:  "slow",
		delay: 200 * time.Millisecond,
		address: domain.Address{
			CEP:          "40279680",
			City:         "Salvador",
			State:        "BA",
			Street:       "Rua Assaré",
			Neighborhood: "Brotas",
		},
	}

	// Provider mais rápido (50ms)
	fast := fakeProvider{
		name:  "fast",
		delay: 50 * time.Millisecond,
		address: domain.Address{
			CEP:          "40279680",
			City:         "Salvador",
			State:        "BA",
			Street:       "Rua Assaré",
			Neighborhood: "Brotas",
		},
	}

	// REQUISITO 1: Requisições simultâneas em duas APIs
	providers := []ports.CEPProvider{slow, fast}
	uc := NewFindFastestAddress(providers, 1*time.Second)

	// Execute a busca concorrente
	address, source, err := uc.Execute(context.Background(), "40279680")

	// Validações
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// REQUISITO 2: A API mais rápida deve ser o resultado
	if source != "fast" {
		t.Fatalf("expected fast provider to win the race, got %s", source)
	}

	if address.CEP != "40279680" {
		t.Fatalf("unexpected cep: %s", address.CEP)
	}
}

// TestExecuteTimeout: REQUISITO 4 (Timeout)
// Verifica que o sistema retorna erro de timeout quando nenhuma
// API responde dentro de 1 segundo
func TestExecuteTimeout(t *testing.T) {
	// Ambos os providers são lentos (2 segundos cada)
	slow1 := fakeProvider{
		name:  "slow1",
		delay: 2 * time.Second,
	}
	slow2 := fakeProvider{
		name:  "slow2",
		delay: 2 * time.Second,
	}

	providers := []ports.CEPProvider{slow1, slow2}
	uc := NewFindFastestAddress(providers, 1*time.Second)

	// Tentar executar deve resultar em timeout
	_, _, err := uc.Execute(context.Background(), "40279680")

	// REQUISITO 4: Erro de timeout deve ser retornado
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("expected timeout error, got %v", err)
	}
}

// TestExecuteOneProviderFailsOtherSucceeds: Caso de falha parcial
// Verifica que o sistema ainda retorna sucesso se uma API falha
// mas a outra responde com sucesso (mesmo que mais lentamente)
func TestExecuteOneProviderFailsOtherSucceeds(t *testing.T) {
	// Provider que falha rapidamente (100ms com erro)
	failing := fakeProvider{
		name:  "failing",
		delay: 100 * time.Millisecond,
		err:   errors.New("connection refused"),
	}

	// Provider que sucede mais lentamente (300ms, mas dentro do timeout)
	succeeding := fakeProvider{
		name:  "succeeding",
		delay: 300 * time.Millisecond,
		address: domain.Address{
			CEP:          "40279680",
			City:         "Salvador",
			State:        "BA",
			Street:       "Rua Assaré",
			Neighborhood: "Brotas",
		},
	}

	providers := []ports.CEPProvider{failing, succeeding}
	uc := NewFindFastestAddress(providers, 1*time.Second)

	// Execute a busca
	address, source, err := uc.Execute(context.Background(), "40279680")

	// Validações
	if err != nil {
		t.Fatalf("expected no error (one provider succeeded), got %v", err)
	}

	// O provider que sucedeu deve ser retornado
	if source != "succeeding" {
		t.Fatalf("expected succeeding provider, got %s", source)
	}

	if address.CEP != "40279680" {
		t.Fatalf("unexpected cep: %s", address.CEP)
	}
}

// TestExecuteBothProvidersFail: Caso de falha total
// Verifica que o sistema retorna erro quando ambas as APIs falham
func TestExecuteBothProvidersFail(t *testing.T) {
	// Ambos os providers falham
	failing1 := fakeProvider{
		name:  "failing1",
		delay: 100 * time.Millisecond,
		err:   errors.New("API 1 error"),
	}

	failing2 := fakeProvider{
		name:  "failing2",
		delay: 150 * time.Millisecond,
		err:   errors.New("API 2 error"),
	}

	providers := []ports.CEPProvider{failing1, failing2}
	uc := NewFindFastestAddress(providers, 1*time.Second)

	// Execute a busca
	_, _, err := uc.Execute(context.Background(), "40279680")

	// Validações
	if err == nil {
		t.Fatalf("expected error when both providers fail, got nil")
	}

	// O erro não deve ser timeout (ambos responderam, mas com erro)
	if errors.Is(err, ErrTimeout) {
		t.Fatalf("expected combined errors, not timeout")
	}
}

// TestExecuteRespeitoConcorrencia: REQUISITO 1 (Requisições Simultâneas)
// Verifica que ambas as APIs são chamadas simultaneamente,
// não sequencialmente
func TestExecuteRespeitoConcorrencia(t *testing.T) {
	startTime := time.Now()

	// Criar providers que registram quando foram chamados
	fast := fakeProvider{
		name:    "fast",
		delay:   50 * time.Millisecond,
		address: domain.Address{CEP: "40279680"},
	}

	slow := fakeProvider{
		name:    "slow",
		delay:   150 * time.Millisecond,
		address: domain.Address{CEP: "40279680"},
	}

	providers := []ports.CEPProvider{fast, slow}
	uc := NewFindFastestAddress(providers, 1*time.Second)

	address, source, err := uc.Execute(context.Background(), "40279680")

	// Validações
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if source != "fast" {
		t.Fatalf("expected fast provider, got %s", source)
	}

	// O tempo total deve ser próximo ao tempo do mais rápido (~50ms)
	// Se fossem sequenciais, levaria ~200ms (50 + 150)
	elapsed := time.Since(startTime)
	if elapsed > 500*time.Millisecond {
		t.Fatalf("execution took %v, indicando que não foi paralelo (deveria ser ~50ms)", elapsed)
	}

	if address.CEP != "40279680" {
		t.Fatalf("unexpected cep: %s", address.CEP)
	}
}
