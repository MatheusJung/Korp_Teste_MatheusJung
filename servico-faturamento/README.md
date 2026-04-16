# Serviço de Faturamento

Microsserviço responsável pela gestão de notas fiscais com Outbox Pattern e idempotência.

## Stack
- **Go 1.22** — net/http + Gin
- **PostgreSQL 16** — banco de dados
- **pgx/v5** — driver PostgreSQL
- **goose** — migrations
- **zap** — logging estruturado

## Arquitetura

Clean Architecture com separação por ports & adapters:

```
internal/
├── domain/
│   ├── entities/     → NotaFiscal, OutboxEvent (regras de negócio puras)
│   ├── enums/        → StatusNota, StatusOutbox
│   └── errors/       → erros de domínio tipados
├── application/
│   ├── ports/        → interfaces (repositórios, cliente estoque, UoW)
│   ├── dtos/         → request/response + mappers
│   └── usecases/     → CriarNota, Imprimir, OutboxWorker
└── infrastructure/
    ├── persistence/  → repositórios pgx + migrations goose
    └── http/estoque/ → cliente HTTP para o serviço de estoque
```

## Endpoints

| Método | Rota                        | Descrição                                      |
|--------|-----------------------------|------------------------------------------------|
| GET    | /api/notas                  | Lista todas as notas                           |
| GET    | /api/notas/:id              | Obtém nota por id                              |
| POST   | /api/notas                  | Cria nova nota (status: Aberta)                |
| POST   | /api/notas/:id/imprimir     | Inicia impressão (202 Accepted)                |
| GET    | /health                     | Health check com status do banco e do estoque  |

## Idempotência

O endpoint `POST /api/notas/:id/imprimir` exige o header `Idempotency-Key: <uuid>`.
Chamadas repetidas com a mesma chave retornam o resultado original sem efeitos colaterais.

## Outbox Pattern

```
Imprimir → [transação atômica] → nota.status = Processando
                               → INSERT outbox_events

OutboxWorker (goroutine, a cada 5s):
  → SELECT ... FOR UPDATE SKIP LOCKED
  → POST estoque/deduzir-lote
    ↳ sucesso  → nota = Fechada, outbox = processado
    ↳ saldo insuficiente → nota = Cancelada, outbox = falhou
    ↳ estoque fora → reagenda com backoff exponencial (5s, 10s, 20s...)
    ↳ esgotou tentativas → nota = Cancelada (timeout)
```

## Execução local

```bash
# Certifique-se que o serviço de estoque está rodando primeiro
cd ../servico-estoque && docker compose up -d

# Subir PostgreSQL + API de faturamento
docker compose up -d

# Acessar API
http://localhost:5002/api/notas

# Health check
http://localhost:5002/health
```

## Testes

```bash
go test ./...
```

## Variáveis de ambiente

| Variável      | Padrão                                                        |
|---------------|---------------------------------------------------------------|
| DATABASE_URL  | postgres://nf:nf@localhost:5432/faturamento?sslmode=disable   |
| ESTOQUE_URL   | http://localhost:5001                                         |
| PORT          | 8080                                                          |
