# Desafio Multithreading - Busca de CEP

Aplicacao em Go que consulta duas APIs de CEP em paralelo e retorna apenas a resposta mais rapida, com timeout global de 1 segundo.

## Objetivo

Implementar um fluxo concorrente para consulta de CEP usando:
- Goroutines
- Channels
- Select
- Context com timeout
- net/http

## APIs utilizadas

- BrasilAPI: https://brasilapi.com.br/api/cep/v1/{cep}
- ViaCEP: https://viacep.com.br/ws/{cep}/json/

## Requisitos Tecnicos (enunciado)

1. Requisicoes simultaneas
   - O sistema deve fazer a requisicao para BrasilAPI e ViaCEP ao mesmo tempo (paralelamente).

2. Race condition (corrida)
   - O sistema deve aceitar apenas a resposta da API que responder mais rapido e descartar a resposta mais lenta.

3. Output (saida)
   - O resultado deve ser exibido na linha de comando (terminal), contendo:
     - Os dados do endereco recebido.
     - Qual API entregou a resposta (BrasilAPI ou ViaCEP).

4. Timeout
   - O tempo limite de resposta deve ser de 1 segundo.
   - Caso nenhuma API responda dentro desse tempo, o sistema deve exibir erro de timeout.

## Requisitos atendidos

- Requisicoes simultaneas para duas APIs
- Race condition controlada: apenas a resposta mais rapida com sucesso e retornada
- Saida em terminal contendo:
   - API vencedora
   - Dados do endereco
- Timeout global de 1 segundo

## Arquitetura do projeto

Estrutura organizada por camadas com separacao de responsabilidades:

- cmd: ponto de entrada da aplicacao
- internal/domain: entidades de dominio
- internal/ports: contratos/interfaces
- internal/adapters/http: integracoes com APIs externas
- internal/usecase: regra de negocio de concorrencia (busca mais rapida)

## Estrutura de pastas

```text
.
├── src
│   ├── cmd
│   │   └── api
│   │       └── main.go
│   └── internal
│       ├── adapters
│       │   └── http
│       │       ├── brasilapi
│       │       │   └── client.go
│       │       └── viacep
│       │           └── client.go
│       ├── domain
│       │   └── address.go
│       ├── ports
│       │   └── cep_provider.go
│       └── usecase
│           ├── find_fasted_address.go
│           └── find_fasted_address_test.go
└── README.md
```

## Pre-requisitos

- Go instalado (recomendado Go 1.22+)
- Conexao com internet para consumir as APIs
- Terminal na raiz do projeto

## Como executar

1. Clone o repositorio:
```bash
git clone https://github.com/rogergcruz/pos-goexpert-desafio-multithreading.git
```

2. Entre na pasta do projeto:
```bash
cd pos-goexpert-desafio-multithreading
```

3. Baixe e organize as dependencias:
```bash
go mod tidy
```

4. Execute a aplicacao com um CEP valido (8 digitos numericos):
```bash
go run ./src/cmd/api 40279680
```

5. Resultado esperado:
- A aplicacao dispara chamadas paralelas para BrasilAPI e ViaCEP
- Retorna apenas a primeira resposta valida
- Imprime no terminal um JSON com API vencedora e endereco

## Exemplo de saida (sucesso)

```json
{
   "api": "ViaCEP",
   "address": {
      "cep": "40279-680",
      "street": "Rua Assare",
      "neighborhood": "Brotas",
      "city": "Salvador",
      "state": "BA"
   }
}
```

## Exemplo de saida (erro de timeout)

Se nenhuma API responder em ate 1 segundo, a aplicacao retorna erro de timeout, por exemplo:

```text
error: timeout for all providers in the given time one second
```

## Logs de tempo de resposta por API

Cada client registra a duracao de sua requisicao. Exemplo de logs no terminal:

```text
[api=BrasilAPI] cep=40279680 duration=85.2ms err=<nil>
[api=ViaCEP] cep=40279680 duration=112.7ms err=<nil>
```

## Como testar

Os testes unitarios cobrem o caso de uso principal de concorrencia e timeout.

1. Execute os testes da camada de use case:
```bash
go test ./src/internal/usecase -v
```

2. Rodar um teste especifico:
```bash
go test ./src/internal/usecase -run TestExecuteTimeout -v
```

## O que os testes validam

- A API mais rapida vence a corrida
- Timeout e respeitado em 1 segundo
- Cenarios de sucesso e falha dos providers no fluxo concorrente

## Tratamento de erros

A aplicacao trata:
- CEP invalido (formato incorreto)
- Timeout global
- Falhas de rede/API
- Respostas invalidas das APIs

## Observacoes

- Se uma API falhar e a outra responder com sucesso dentro do timeout, a aplicacao retorna sucesso
- Se ambas falharem, retorna erro consolidado
- O foco esta em concorrencia, previsibilidade do timeout e clareza de saida em CLI