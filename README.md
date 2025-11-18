# Sistema de Emissão de Notas Fiscais

Projeto técnico: Sistema de emissão de notas fiscais com arquitetura de microsserviços e Clean Architecture.

## Visão Geral
O sistema possui os seguintes componentes:

- **InventoryService** (Postgres): cadastro de produtos e controle de estoque.
- **BillingService** (SQL Server): criação, atualização e impressão de notas fiscais.
- **Frontend**: Angular 20 + RxJS, consumindo os serviços via REST.
- **Orquestração**: dois Docker Compose separados para os backends; frontend pode rodar fora do container ou em container próprio.

Cada serviço expõe um endpoint de **health check** (`/health`) para monitoramento.

---

## Arquitetura Limpa (Clean Architecture)

Cada backend segue **arquitetura limpa**, separando responsabilidades em camadas:

```
Services/
├─ InventoryService/
│  ├─ InventoryService.Api          # Controllers / Endpoints
│  ├─ InventoryService.Application  # Casos de uso / Services
│  ├─ InventoryService.Domain       # Entidades / Regras de negócio
│  └─ InventoryService.Infrastructure # EF Core, Repositórios, DB
│
├─ BillingService/
│  ├─ BillingService.Api
│  ├─ BillingService.Application
│  ├─ BillingService.Domain
│  └─ BillingService.Infrastructure
```

**Camadas principais:**

1. **Domain**: entidades, agregados, regras de negócio puras, enums.
2. **Application**: casos de uso (services), DTOs, interfaces de repositórios.
3. **Infrastructure**: persistência (EF Core), acesso HTTP a outros serviços, mapeamento para o banco.
4. **API**: controllers, endpoints, validação de requests, tratamento de erros HTTP.

> Essa separação garante **testabilidade**, **manutenção mais fácil** e **independência de frameworks/bancos**.

---

## Como Rodar Localmente

### 1. Criar network Docker compartilhada
```bash
docker network create invoice_network
```

### 2. Rodar InventoryService
```bash
cd services/inventoryservice
docker compose up --build -d
```
> Configure no `docker-compose.yml` do InventoryService para usar a network criada:
```yaml
networks:
  default:
    external:
      name: invoice_network
```

### 3. Rodar BillingService
```bash
cd services/billingservice
docker compose up --build -d
```
> Também configure a mesma network (`invoice_network`) para que o BillingService consiga acessar o InventoryService.

### 4. Frontend
- Fora do container:
```bash
cd frontend
npm ci
npm run start
```
- Em container próprio, conecte-o também à network `invoice_network`.

5. Acesse:
```
http://localhost:4200
```

> **Observação**: BillingService depende do InventoryService; ambos devem estar ativos e na mesma network.

---

## Endpoints Principais

### InventoryService
**Health**
- `GET /health`

**Products**
- `POST /products` — cria um produto
- `GET /products` — lista produtos com paginação
- `GET /products/Products` — lista todos os produtos
- `GET /products/ActiveProducts` — lista produtos ativos
- `GET /products/{productCode}` — consulta produto por código
- `DELETE /products/{productCode}` — desativa produto

**Stock**
- `GET /stock` — lista movimentações do estoque
- `POST /stock/add` — adiciona quantidade ao estoque
- `POST /stock/remove` — remove quantidade do estoque

### BillingService
**Health**
- `GET /health`

**Invoices**
- `GET /invoices/all` — lista todas as notas fiscais
- `GET /invoices` — lista notas fiscais com paginação
- `POST /invoices` — cria nova nota fiscal com lista de itens
- `GET /invoices/{seqNumber}` — consulta nota fiscal pelo número sequencial
- `PUT /invoices/{seqNumber}/items` — adiciona item a nota aberta
- `POST /invoices/{seqNumber}/cancel` — cancela nota
- `POST /invoices/{seqNumber}/close` — fecha nota aberta
- `POST /invoices/{seqNumber}/close-and-print` — imprime e fecha nota aberta

**Exemplo de body para criar nota:**
```json
{
  "sequentialNumber": 12345,
  "items": [
    { "productCode": "P001", "quantity": 2 },
    { "productCode": "P002", "quantity": 1 }
  ]
}
```

---

## Geração de PDF de Notas Fiscais

O sistema permite gerar e baixar PDFs das notas fiscais diretamente do backend.  

### Endpoint para gerar PDF
- **POST** `/invoices/{seqNumber}/close-and-print`
  - Fecha a nota fiscal (se estiver aberta) e gera o PDF.
  - Retorna o PDF como `application/pdf`.

### Como o frontend lida com o PDF
- Ao clicar no botão “Imprimir” na interface Angular, a aplicação chama o endpoint acima.
- O PDF pode ser:
  1. **Visualizado em nova aba**:
  ```ts
  this.http.post(`/api/billing/invoices/${seqNumber}/close-and-print`, {}, { responseType: 'blob' })
    .subscribe(blob => {
      const url = window.URL.createObjectURL(blob);
      window.open(url);
    });
  ```
  2. **Baixado automaticamente**:
  ```ts
  this.http.post(`/api/billing/invoices/${seqNumber}/close-and-print`, {}, { responseType: 'blob' })
    .subscribe(blob => {
      const a = document.createElement('a');
      a.href = window.URL.createObjectURL(blob);
      a.download = `Invoice_${seqNumber}.pdf`;
      a.click();
    });
  ```

### Biblioteca utilizada no backend
- **QuestPDF** (.NET 8) — gera PDFs de forma programática, integrando dados da nota fiscal (número sequencial, status e itens).
- Garante que o PDF seja **consistente** com os dados do banco e **profissionalmente formatado**.

> Observação: PDF é gerado apenas no backend para garantir a consistência dos dados oficiais da nota fiscal.

---

## Tratamento de Falhas
- Health check em `/health`.
- Retry + Circuit Breaker usando **Polly** no BillingService ao chamar InventoryService.
- Transações distribuídas via **saga / compensação** para manter consistência.
- Erros HTTP claros:
  - `400` — validação / regras de negócio
  - `404` — recurso não encontrado
  - `409` — conflito de saldo
  - `503` — serviço externo indisponível

---

## Tecnologias e Bibliotecas
- **Backend (.NET 8)**:
  - Entity Framework Core (ORM)
  - Polly (Retry / Circuit Breaker)
  - Swashbuckle/Swagger (Documentação)
  - QuestPDF (Geração de PDFs)
  - Clean Architecture (Domain, Application, Infrastructure, API)
- **Frontend (Angular 20)**:
  - RxJS (reatividade)
  - Angular Material / Bootstrap (UI)
  - ngx-toastr (notificações)
- **Orquestração / Containers**:
  - Docker / Docker Compose
  - Network Docker compartilhada (`invoice_network`)

---

## Observações
- Cada backend possui seu `docker-compose.yml` separado; ambos devem estar na **mesma network Docker**.
- Frontend consome os endpoints via REST e deve apontar para os containers correspondentes.

