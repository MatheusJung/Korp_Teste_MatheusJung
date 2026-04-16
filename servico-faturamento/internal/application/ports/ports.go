package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/nf-system/servico-faturamento/internal/domain/entities"
)

type ListarNotasFiltro struct {
	Page          int
	PageSize      int
	Search        string
	SortBy        string
	SortDirection string
}

type NotaPDFService interface {
	Gerar(ctx context.Context, nota *entities.NotaFiscal) ([]byte, error)
}

type NotaPDFStorage interface {
	Salvar(ctx context.Context, notaID uuid.UUID, numero int64, conteudo []byte) error
	Obter(ctx context.Context, notaID uuid.UUID, numero int64) ([]byte, string, error)
}

type NotaFiscalRepository interface {
	ProximoNumero(ctx context.Context) (int64, error)
	Salvar(ctx context.Context, nota *entities.NotaFiscal) error
	ObterPorID(ctx context.Context, id uuid.UUID) (*entities.NotaFiscal, error)
	Listar(ctx context.Context) ([]*entities.NotaFiscal, error)
	ListarPaginado(ctx context.Context, filtro ListarNotasFiltro) ([]*entities.NotaFiscal, int64, error)
	Atualizar(ctx context.Context, nota *entities.NotaFiscal) error
}

type OutboxRepository interface {
	Salvar(ctx context.Context, event *entities.OutboxEvent) error
	// FOR UPDATE SKIP LOCKED — evita que dois workers peguem o mesmo evento
	BuscarPendentes(ctx context.Context, limite int) ([]*entities.OutboxEvent, error)
	Atualizar(ctx context.Context, event *entities.OutboxEvent) error
}

type IdempotencyRepository interface {
	Existe(ctx context.Context, chave string) (bool, error)
	Registrar(ctx context.Context, chave string, resposta []byte) error
	ObterResposta(ctx context.Context, chave string) ([]byte, error)
}

// UnitOfWork agrupa nota + outbox numa mesma transação
type UnitOfWork interface {
	Execute(ctx context.Context, fn func(repos TxRepositories) error) error
}

type TxRepositories struct {
	Notas  NotaFiscalRepository
	Outbox OutboxRepository
}

// EstoqueClient — porta para o serviço externo de estoque
type EstoqueClient interface {
	DeduzirLote(ctx context.Context, req DeduzirLoteRequest) (*DeduzirLoteResponse, error)
	EstornarLote(ctx context.Context, req EstornarLoteRequest) error
	HealthCheck(ctx context.Context) error
}

type DeduzirLoteRequest struct {
	NotaFiscalID uuid.UUID
	Itens        []ItemEstoque
}

type EstornarLoteRequest struct {
	NotaFiscalID uuid.UUID
	Itens        []ItemEstoque
}

type ItemEstoque struct {
	ProdutoID  uuid.UUID
	Quantidade float64
}

type DeduzirLoteResponse struct {
	Sucesso        bool
	ItensDeduzidos []ItemEstoque
	Erro           string
	ProdutoFalhou  *uuid.UUID
}
