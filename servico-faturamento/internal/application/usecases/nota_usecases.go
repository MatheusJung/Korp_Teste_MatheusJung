package usecases

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/nf-system/servico-faturamento/internal/application/dtos"
	"github.com/nf-system/servico-faturamento/internal/application/ports"
	"github.com/nf-system/servico-faturamento/internal/domain/entities"
	domainerrors "github.com/nf-system/servico-faturamento/internal/domain/errors"
	"go.uber.org/zap"
)

// CriarNotaUseCase cria uma nova nota fiscal com status Aberta.
type CriarNotaUseCase struct {
	uow    ports.UnitOfWork
	notas  ports.NotaFiscalRepository
	logger *zap.Logger
}

func NewCriarNotaUseCase(uow ports.UnitOfWork, notas ports.NotaFiscalRepository, logger *zap.Logger) *CriarNotaUseCase {
	return &CriarNotaUseCase{uow: uow, notas: notas, logger: logger}
}

func (uc *CriarNotaUseCase) Executar(ctx context.Context, req dtos.CriarNotaRequest) (*dtos.NotaFiscalResponse, error) {
	itens := make([]entities.ItemNota, len(req.Itens))
	for i, item := range req.Itens {
		itens[i] = entities.ItemNota{
			ID:               uuid.New(),
			ProdutoID:        item.ProdutoID,
			ProdutoCodigo:    item.ProdutoCodigo,
			ProdutoDescricao: item.ProdutoDescricao,
			Quantidade:       item.Quantidade,
		}
	}

	nota, err := entities.NovaNotaFiscal(itens)
	if err != nil {
		return nil, err
	}

	for i := range nota.Itens {
		nota.Itens[i].NotaID = nota.ID
	}

	numero, err := uc.notas.ProximoNumero(ctx)
	if err != nil {
		return nil, err
	}
	nota.Numero = numero

	if err := uc.notas.Salvar(ctx, nota); err != nil {
		return nil, err
	}

	uc.logger.Info("nota fiscal criada", zap.String("id", nota.ID.String()), zap.Int64("numero", nota.Numero))
	resp := dtos.ToNotaResponse(nota)
	return &resp, nil
}

// ObterNotaUseCase retorna uma nota por ID.
type ObterNotaUseCase struct {
	notas ports.NotaFiscalRepository
}

func NewObterNotaUseCase(notas ports.NotaFiscalRepository) *ObterNotaUseCase {
	return &ObterNotaUseCase{notas: notas}
}

func (uc *ObterNotaUseCase) Executar(ctx context.Context, id uuid.UUID) (*dtos.NotaFiscalResponse, error) {
	nota, err := uc.notas.ObterPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dtos.ToNotaResponse(nota)
	return &resp, nil
}

// ListarNotasUseCase lista todas as notas.
type ListarNotasUseCase struct {
	notas ports.NotaFiscalRepository
}

func NewListarNotasUseCase(notas ports.NotaFiscalRepository) *ListarNotasUseCase {
	return &ListarNotasUseCase{notas: notas}
}

func (uc *ListarNotasUseCase) Executar(ctx context.Context) ([]dtos.NotaFiscalResponse, error) {
	notas, err := uc.notas.Listar(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dtos.NotaFiscalResponse, len(notas))
	for i, n := range notas {
		result[i] = dtos.ToNotaResponse(n)
	}
	return result, nil
}

type ListarNotasPaginadoUseCase struct {
	notas ports.NotaFiscalRepository
}

func NewListarNotasPaginadoUseCase(notas ports.NotaFiscalRepository) *ListarNotasPaginadoUseCase {
	return &ListarNotasPaginadoUseCase{notas: notas}
}

func (uc *ListarNotasPaginadoUseCase) Executar(
	ctx context.Context,
	query dtos.ListarNotasQuery,
) (*dtos.NotaFiscalPagedResponse, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	if strings.TrimSpace(query.SortBy) == "" {
		query.SortBy = "criado_em"
	}
	if strings.TrimSpace(query.SortDirection) == "" {
		query.SortDirection = "desc"
	}

	notas, total, err := uc.notas.ListarPaginado(ctx, ports.ListarNotasFiltro{
		Page:          query.Page,
		PageSize:      query.PageSize,
		Search:        strings.TrimSpace(query.Search),
		SortBy:        query.SortBy,
		SortDirection: query.SortDirection,
	})
	if err != nil {
		return nil, err
	}

	items := make([]dtos.NotaFiscalResponse, len(notas))
	for i, n := range notas {
		items[i] = dtos.ToNotaResponse(n)
	}

	totalPages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))

	return &dtos.NotaFiscalPagedResponse{
		Items:      items,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

// ImprimirNotaUseCase inicia o processo de impressão da nota.
// Grava nota como Processando + outbox numa transação atômica.
// Retorna 202 Accepted imediatamente — o worker processa em background.
type ImprimirNotaUseCase struct {
	uow          ports.UnitOfWork
	notas        ports.NotaFiscalRepository
	idempotency  ports.IdempotencyRepository
	maxTentativas int
	logger       *zap.Logger
}

func NewImprimirNotaUseCase(
	uow ports.UnitOfWork,
	notas ports.NotaFiscalRepository,
	idempotency ports.IdempotencyRepository,
	maxTentativas int,
	logger *zap.Logger,
) *ImprimirNotaUseCase {
	return &ImprimirNotaUseCase{
		uow:           uow,
		notas:         notas,
		idempotency:   idempotency,
		maxTentativas: maxTentativas,
		logger:        logger,
	}
}

func (uc *ImprimirNotaUseCase) Executar(ctx context.Context, notaID uuid.UUID, idempotencyKey string) (*dtos.ImprimirResponse, error) {
	// Verificar idempotência — se já foi processada, retorna o resultado cached
	if existe, _ := uc.idempotency.Existe(ctx, idempotencyKey); existe {
		respBytes, err := uc.idempotency.ObterResposta(ctx, idempotencyKey)
		if err != nil {
			return nil, err
		}
		var resp dtos.ImprimirResponse
		if err := json.Unmarshal(respBytes, &resp); err != nil {
			return nil, err
		}
		uc.logger.Info("requisição idempotente retornada", zap.String("key", idempotencyKey))
		return &resp, domainerrors.ErrIdempotenciaRepetida
	}

	nota, err := uc.notas.ObterPorID(ctx, notaID)
	if err != nil {
		return nil, err
	}

	if !nota.PodeImprimir() {
		return nil, domainerrors.ErrNotaNaoPodeImprimir
	}

	// Transação atômica: muda status + grava outbox
	err = uc.uow.Execute(ctx, func(repos ports.TxRepositories) error {
		if err := nota.IniciarProcessamento(); err != nil {
			return err
		}
		if err := repos.Notas.Atualizar(ctx, nota); err != nil {
			return err
		}

		outboxEvent := entities.NovoOutboxEvent(nota.ID, nota.Itens, uc.maxTentativas)
		return repos.Outbox.Salvar(ctx, outboxEvent)
	})
	if err != nil {
		return nil, err
	}

	resp := &dtos.ImprimirResponse{
		NotaID:   nota.ID,
		Status:   nota.Status,
		Mensagem: "Impressão iniciada. A nota será fechada após dedução do estoque.",
	}

	// Persiste resposta para idempotência futura
	respBytes, _ := json.Marshal(resp)
	_ = uc.idempotency.Registrar(ctx, idempotencyKey, respBytes)

	uc.logger.Info("impressão iniciada", zap.String("notaId", nota.ID.String()))
	return resp, nil
}

