# Frontend — Sistema de Notas Fiscais

Angular 20 · Standalone Components · RxJS · SCSS

## Pré-requisitos

- Node.js 20+
- npm 10+
- Serviços de estoque (porta 5001) e faturamento (porta 5002) rodando

## Instalação

```bash
npm install
```

> O `NotaFiscalService` usa `uuid` para gerar o `Idempotency-Key`.
> Instale a tipagem:
> ```bash
> npm install uuid
> npm install --save-dev @types/uuid
> ```

## Execução

```bash
npm start
# Acesse http://localhost:4200
```

O proxy redireciona automaticamente:
- `/api/estoque/*` → `http://localhost:5001/api/*`
- `/api/faturamento/*` → `http://localhost:5002/api/*`
- `/health/estoque` → `http://localhost:5001/health`
- `/health/faturamento` → `http://localhost:5002/health`

## Funcionalidades

- Cadastro de produtos com saldo inicial
- Listagem de produtos com saldo em tempo real
- Entrada manual de estoque
- Histórico de movimentações por produto
- Criação de notas fiscais com múltiplos itens
- Impressão de notas (Outbox Pattern — status Processando → Fechada/Cancelada)
- Polling automático a cada 3s enquanto nota estiver Processando
- Banner de saúde dos serviços (atualiza a cada 30s)
- Tratamento de erros com mensagens amigáveis

## Estrutura

```
src/app/
├── core/
│   ├── models/       → interfaces TypeScript
│   ├── services/     → ProdutoService, NotaFiscalService, HealthService
│   └── interceptors/ → ErrorInterceptor
├── features/
│   ├── produtos/     → ProdutosPageComponent
│   └── notas-fiscais/→ NotasPageComponent
└── app.component.ts  → shell com header, tabs e health banner
```
