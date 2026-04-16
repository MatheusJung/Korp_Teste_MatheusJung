# Serviço de Estoque

Microsserviço responsável pelo controle de produtos e saldos em estoque.

## Stack
- **.NET 8** — ASP.NET Core Web API
- **SQL Server 2022** — banco de dados
- **Entity Framework Core 8** — ORM com LINQ
- **Serilog** — logging estruturado

## Arquitetura

Clean Architecture com 4 camadas:

```
Estoque.Domain          → Entidades, regras de negócio, exceções (sem dependências)
Estoque.Application     → Use cases, DTOs, interfaces de repositório
Estoque.Infrastructure  → EF Core, repositórios concretos, migrations
Estoque.API             → Controllers, middleware, health check, DI
```

## Endpoints

| Método | Rota                                  | Descrição                              |
|--------|---------------------------------------|----------------------------------------|
| GET    | /api/produtos                         | Lista todos os produtos                |
| GET    | /api/produtos/{id}                    | Obtém produto por id                   |
| POST   | /api/produtos                         | Cadastra novo produto                  |
| POST   | /api/movimentacoes/deduzir-lote       | Deduz estoque (chamado pelo faturamento)|
| POST   | /api/movimentacoes/estornar-lote      | Estorna estoque (nota cancelada)       |
| POST   | /api/movimentacoes/entrada            | Entrada manual (reposição)             |
| GET    | /api/movimentacoes/produto/{id}       | Histórico de movimentações             |
| GET    | /health                               | Health check com status do banco       |

## Controle de concorrência

Utiliza **Optimistic Concurrency** via `RowVersion` (timestamp) no SQL Server.
Quando duas requisições tentam modificar o mesmo produto simultaneamente,
o EF Core lança `DbUpdateConcurrencyException`, que é capturada e retornada
como HTTP 409 Conflict.

## Execução local

```bash
# Subir SQL Server + API
docker compose up -d

# Acessar Swagger
http://localhost:5001/swagger

# Health check
http://localhost:5001/health
```

## Testes

```bash
dotnet test tests/Estoque.UnitTests
```

## Variáveis de ambiente

| Variável                          | Padrão                          |
|-----------------------------------|---------------------------------|
| ConnectionStrings__SqlServer      | ver appsettings.json            |
| ASPNETCORE_ENVIRONMENT            | Production                      |
