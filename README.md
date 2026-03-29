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

1. Inicialize o módulo (se ainda não existir):
   go mod init pos-goexpert-desafio-multithreading

2. Baixe dependências:
   go mod tidy

3. Execute:
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

## Observações

- Se nenhuma API responder em até 1 segundo, a aplicação retorna erro de timeout.
- Caso uma API falhe, a aplicação ainda pode retornar sucesso se a outra responder primeiro com sucesso.