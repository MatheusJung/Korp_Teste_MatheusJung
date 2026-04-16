# Sistema de Emissão de Notas Fiscais

Sistema de emissão de notas fiscais com arquitetura de microsserviços, desenvolvido como projeto de avaliação técnica.

## Visão geral

```
├── servico-estoque/        → .NET 8 + SQL Server
├── servico-faturamento/    → Go + PostgreSQL
└── frontend-angular/       → Angular 20
```

## Stack

| Camada | Tecnologia |
|---|---|
| Frontend | Angular 20 · Standalone Components · RxJS · SCSS |
| Serviço de Estoque | .NET 8 · ASP.NET Core · Entity Framework Core · SQL Server |
| Serviço de Faturamento | Go 1.23 · Gin · pgx · PostgreSQL |
| Infraestrutura | Docker · Docker Compose |

## Arquitetura

### Microsserviços

**Serviço de Estoque** (porta 5001)
Responsável pelo controle de produtos, saldos e movimentações. Expõe endpoints REST consumidos diretamente pelo frontend e pelo Serviço de Faturamento.

**Serviço de Faturamento** (porta 5002)
Responsável pela gestão de notas fiscais. Implementa o **Outbox Pattern** para garantir consistência na dedução de estoque sem acoplamento síncrono.

**Frontend** (porta 4200)
Single Page Application com proxy reverso para os dois backends.

### Fluxo de impressão (Outbox Pattern)

```
Usuário clica "Imprimir"
  → Faturamento: nota = Processando + evento gravado na tabela outbox (transação atômica)
  → Resposta 202 Accepted imediata ao usuário
  → Outbox Worker (goroutine, polling a cada 5s):
      → Chama Estoque: POST /api/movimentacoes/deduzir-lote
          → Sucesso: nota = Fechada, saldos deduzidos
          → Saldo insuficiente: estorno parcial automático, nota = Cancelada
          → Estoque indisponível: retry com backoff exponencial (5s → 10s → 20s...)
          → Esgotou tentativas: nota = Cancelada (timeout)
  → Frontend: polling a cada 3s enquanto status = Processando
```

### Requisitos implementados

- [x] Cadastro de produtos com saldo inicial
- [x] Cadastro de notas fiscais com múltiplos itens
- [x] Impressão com atualização de status e dedução de estoque
- [x] Arquitetura de microsserviços (Estoque + Faturamento)
- [x] Tratamento de falhas com feedback ao usuário
- [x] Conexão real com banco de dados (SQL Server + PostgreSQL)
- [x] Controle de concorrência (RowVersion / Optimistic Locking)
- [x] Idempotência (Idempotency-Key no header HTTP)
- [x] Histórico de movimentações por produto (Criação, Entrada, Saída, Estorno)

## Pré-requisitos

- Docker Desktop
- Node.js 20+ e npm (para o frontend)
- Go 1.23+ (opcional, para rodar fora do Docker)

## Como executar

### 1. Rede compartilhada

Os dois serviços precisam se enxergar via Docker network. Suba o estoque primeiro — ele cria a rede `nf-network`:

```bash
cd servico-estoque
docker compose up -d
```

### 2. Serviço de Faturamento

```bash
cd servico-faturamento
docker compose up -d
```

### 3. Frontend

```bash
cd frontend-angular
npm install
npm start
```

Acesse: http://localhost:4200

### URLs disponíveis

| Serviço | URL |
|---|---|
| Frontend | http://localhost:4200 |
| Estoque API | http://localhost:5001 |
| Estoque Swagger | http://localhost:5001/swagger |
| Estoque Health | http://localhost:5001/health |
| Faturamento API | http://localhost:5002 |
| Faturamento Swagger | http://localhost:5002/swagger/index.html |
| Faturamento Health | http://localhost:5002/health |

## Estrutura dos projetos

### Serviço de Estoque (.NET 8)

Clean Architecture com 4 camadas:

```
Estoque.Domain          → Entidades, regras de negócio, exceções
Estoque.Application     → Use cases, DTOs, interfaces de repositório
Estoque.Infrastructure  → EF Core, repositórios, migrations
Estoque.API             → Controllers, middleware, health check
```

### Serviço de Faturamento (Go)

Ports & Adapters:

```
internal/domain/        → Entidades, enums, erros de domínio
internal/application/   → Use cases, DTOs, ports (interfaces)
internal/infrastructure/→ Repositórios pgx, cliente HTTP do estoque, migrations
internal/handler/       → Handlers HTTP (Gin), middleware
cmd/api/                → Bootstrap (main.go)
```

### Frontend (Angular 20)

Feature-based com standalone components:

```
core/models/            → Interfaces TypeScript
core/services/          → ProdutoService, NotaFiscalService, HealthService
core/interceptors/      → ErrorInterceptor
features/produtos/      → ProdutosPageComponent, modais
features/notas-fiscais/ → NotasPageComponent, ModalCriarNota, ModalDetalhesNota
```

## Parar os serviços

```bash
# Para cada pasta (servico-estoque, servico-faturamento)
docker compose down

# Para remover volumes (apaga dados do banco)
docker compose down -v
```
