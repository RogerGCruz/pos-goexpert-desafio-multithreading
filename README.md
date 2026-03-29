# Desafio Multithreading - Busca de CEP

Aplicação em Go que consulta BrasilAPI e ViaCEP em paralelo e retorna apenas a resposta mais rápida.

## Requisitos atendidos

- Requisições simultâneas para duas APIs
- Retorno da API mais rápida
- Timeout global de 1 segundo
- Exibição no terminal do endereço e da API vencedora

## Tecnologias

- Go
- Goroutines
- Channels
- Select
- net/http
- context

## Como executar

1. Baixe as dependências:
   go mod tidy

2. Execute:
   go run ./src/cmd/api 40279680

## Exemplo de saída

{
  "api": "ViaCEP",
  "address": {
    "cep": "40279-680",
    "street": "Rua Assaré",
    "neighborhood": "Brotas",
    "city": "Salvador",
    "state": "BA"
  }
}

## Exemplo de erro de timeout

Se nenhuma API responder em até 1 segundo, a aplicação retorna erro de timeout.

## Observações

- Se nenhuma API responder em até 1 segundo, a aplicação retorna erro de timeout.
- Caso uma API falhe, a aplicação ainda pode retornar sucesso se a outra responder primeiro com sucesso.