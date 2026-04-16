package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nf-system/servico-faturamento/internal/application/ports"
	"github.com/nf-system/servico-faturamento/internal/application/usecases"
	"github.com/nf-system/servico-faturamento/internal/domain/entities"
	"go.uber.org/zap"
)

// -----------------------------------------------------------------------------
// Helper para adaptar os testes ao novo construtor do worker
// -----------------------------------------------------------------------------

func newWorkerForTest(
	outboxRepo ports.OutboxRepository,
	notaRepo ports.NotaFiscalRepository,
	estoqueClient ports.EstoqueClient,
	uow ports.UnitOfWork,
	logger *zap.Logger,
) *usecases.OutboxWorker {
	return usecases.NewOutboxWorker(
		outboxRepo,
		notaRepo,
		estoqueClient,
		uow,
		nil, // pdfService
		nil, // pdfStorage
		5*time.Second,
		5,
		logger,
	)
}

// -----------------------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------------------

type mockOutboxRepo struct {
	buscarPendentesFn func(ctx context.Context, limite int) ([]*entities.OutboxEvent, error)
	atualizarFn       func(ctx context.Context, event *entities.OutboxEvent) error
	salvarFn          func(ctx context.Context, event *entities.OutboxEvent) error
}

func (m *mockOutboxRepo) BuscarPendentes(ctx context.Context, limite int) ([]*entities.OutboxEvent, error) {
	if m.buscarPendentesFn != nil {
		return m.buscarPendentesFn(ctx, limite)
	}
	return nil, nil
}

func (m *mockOutboxRepo) Atualizar(ctx context.Context, event *entities.OutboxEvent) error {
	if m.atualizarFn != nil {
		return m.atualizarFn(ctx, event)
	}
	return nil
}

func (m *mockOutboxRepo) Salvar(ctx context.Context, event *entities.OutboxEvent) error {
	if m.salvarFn != nil {
		return m.salvarFn(ctx, event)
	}
	return nil
}

type mockNotaRepo struct {
	obterPorIDFn func(ctx context.Context, id uuid.UUID) (*entities.NotaFiscal, error)
	atualizarFn  func(ctx context.Context, nota *entities.NotaFiscal) error
	salvarFn     func(ctx context.Context, nota *entities.NotaFiscal) error
	listarFn     func(ctx context.Context) ([]*entities.NotaFiscal, error)
	proximoFn    func(ctx context.Context) (int64, error)
}

func (m *mockNotaRepo) ObterPorID(ctx context.Context, id uuid.UUID) (*entities.NotaFiscal, error) {
	if m.obterPorIDFn != nil {
		return m.obterPorIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockNotaRepo) Atualizar(ctx context.Context, nota *entities.NotaFiscal) error {
	if m.atualizarFn != nil {
		return m.atualizarFn(ctx, nota)
	}
	return nil
}

func (m *mockNotaRepo) Salvar(ctx context.Context, nota *entities.NotaFiscal) error {
	if m.salvarFn != nil {
		return m.salvarFn(ctx, nota)
	}
	return nil
}

func (m *mockNotaRepo) Listar(ctx context.Context) ([]*entities.NotaFiscal, error) {
	if m.listarFn != nil {
		return m.listarFn(ctx)
	}
	return nil, nil
}

func (m *mockNotaRepo) ListarPaginado(
	ctx context.Context,
	filtro ports.ListarNotasFiltro,
) ([]*entities.NotaFiscal, int64, error) {
	return []*entities.NotaFiscal{}, 0, nil
}

func (m *mockNotaRepo) ProximoNumero(ctx context.Context) (int64, error) {
	if m.proximoFn != nil {
		return m.proximoFn(ctx)
	}
	return m.proximoFn(ctx)
}

type mockEstoqueClient struct {
	deduzirLoteFn func(ctx context.Context, req ports.DeduzirLoteRequest) (*ports.DeduzirLoteResponse, error)
	estornarFn    func(ctx context.Context, req ports.EstornarLoteRequest) error
	healthFn      func(ctx context.Context) error
}

func (m *mockEstoqueClient) DeduzirLote(ctx context.Context, req ports.DeduzirLoteRequest) (*ports.DeduzirLoteResponse, error) {
	if m.deduzirLoteFn != nil {
		return m.deduzirLoteFn(ctx, req)
	}
	return nil, nil
}

func (m *mockEstoqueClient) EstornarLote(ctx context.Context, req ports.EstornarLoteRequest) error {
	if m.estornarFn != nil {
		return m.estornarFn(ctx, req)
	}
	return nil
}

func (m *mockEstoqueClient) HealthCheck(ctx context.Context) error {
	if m.healthFn != nil {
		return m.healthFn(ctx)
	}
	return nil
}

type mockUoW struct {
	executeFn func(ctx context.Context, fn func(repos ports.TxRepositories) error) error
}

func (m *mockUoW) Execute(ctx context.Context, fn func(repos ports.TxRepositories) error) error {
	if m.executeFn != nil {
		return m.executeFn(ctx, fn)
	}
	return fn(ports.TxRepositories{})
}

// -----------------------------------------------------------------------------
// Helpers de teste
// -----------------------------------------------------------------------------

func mustNovaNotaFiscal(t *testing.T) *entities.NotaFiscal {
	t.Helper()

	nota, err := entities.NovaNotaFiscal([]entities.ItemNota{
		{
			ProdutoID:  uuid.New(),
			Quantidade: 1,
		},
	})
	if err != nil {
		t.Fatalf("erro ao criar nota de teste: %v", err)
	}
	return nota
}

// -----------------------------------------------------------------------------
// Testes
// -----------------------------------------------------------------------------

func TestOutboxWorker_Processar_NaoFalhaAoBuscarPendentes(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	outboxRepo := &mockOutboxRepo{
		buscarPendentesFn: func(ctx context.Context, limite int) ([]*entities.OutboxEvent, error) {
			return nil, errors.New("erro ao buscar")
		},
	}
	notaRepo := &mockNotaRepo{}
	estoqueClient := &mockEstoqueClient{}
	uow := &mockUoW{}

	worker := newWorkerForTest(outboxRepo, notaRepo, estoqueClient, uow, logger)

	worker.ProcessarParaTeste(ctx)
}

func TestOutboxWorker_ProcessarEvento_ReagendaQuandoEstoqueFalhaTecnica(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	nota := mustNovaNotaFiscal(t)
	_ = nota.IniciarProcessamento()

	evento := entities.NovoOutboxEvent(nota.ID, nota.Itens, 5)

	atualizouOutbox := false

	outboxRepo := &mockOutboxRepo{
		buscarPendentesFn: func(ctx context.Context, limite int) ([]*entities.OutboxEvent, error) {
			return []*entities.OutboxEvent{evento}, nil
		},
		atualizarFn: func(ctx context.Context, event *entities.OutboxEvent) error {
			atualizouOutbox = true
			return nil
		},
	}

	notaRepo := &mockNotaRepo{
		obterPorIDFn: func(ctx context.Context, id uuid.UUID) (*entities.NotaFiscal, error) {
			return nota, nil
		},
	}

	estoqueClient := &mockEstoqueClient{
		deduzirLoteFn: func(ctx context.Context, req ports.DeduzirLoteRequest) (*ports.DeduzirLoteResponse, error) {
			return nil, errors.New("timeout")
		},
	}

	uow := &mockUoW{}

	worker := newWorkerForTest(outboxRepo, notaRepo, estoqueClient, uow, logger)

	worker.ProcessarParaTeste(ctx)

	if !atualizouOutbox {
		t.Fatal("esperava atualização do outbox ao reagendar evento")
	}
}

func TestOutboxWorker_ProcessarEvento_CancelaNotaQuandoErroDeNegocio(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	nota := mustNovaNotaFiscal(t)
	_ = nota.IniciarProcessamento()

	evento := entities.NovoOutboxEvent(nota.ID, nota.Itens, 5)

	atualizouNota := false
	atualizouOutbox := false

	outboxRepo := &mockOutboxRepo{
		buscarPendentesFn: func(ctx context.Context, limite int) ([]*entities.OutboxEvent, error) {
			return []*entities.OutboxEvent{evento}, nil
		},
		atualizarFn: func(ctx context.Context, event *entities.OutboxEvent) error {
			atualizouOutbox = true
			return nil
		},
	}

	notaRepo := &mockNotaRepo{
		obterPorIDFn: func(ctx context.Context, id uuid.UUID) (*entities.NotaFiscal, error) {
			return nota, nil
		},
		atualizarFn: func(ctx context.Context, nota *entities.NotaFiscal) error {
			atualizouNota = true
			return nil
		},
	}

	estoqueClient := &mockEstoqueClient{
		deduzirLoteFn: func(ctx context.Context, req ports.DeduzirLoteRequest) (*ports.DeduzirLoteResponse, error) {
			return &ports.DeduzirLoteResponse{
				Sucesso: false,
				Erro:    "conflito de negócio",
			}, nil
		},
	}

	uow := &mockUoW{
		executeFn: func(ctx context.Context, fn func(repos ports.TxRepositories) error) error {
			repos := ports.TxRepositories{
				Notas:  notaRepo,
				Outbox: outboxRepo,
			}
			return fn(repos)
		},
	}

	worker := newWorkerForTest(outboxRepo, notaRepo, estoqueClient, uow, logger)

	worker.ProcessarParaTeste(ctx)

	if !atualizouNota {
		t.Fatal("esperava atualização da nota ao cancelar")
	}
	if !atualizouOutbox {
		t.Fatal("esperava atualização do outbox ao marcar falha")
	}
}

func TestOutboxWorker_ProcessarEvento_FechaNotaQuandoSucesso(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	nota := mustNovaNotaFiscal(t)
	_ = nota.IniciarProcessamento()

	evento := entities.NovoOutboxEvent(nota.ID, nota.Itens, 5)

	atualizouNota := false
	atualizouOutbox := false

	outboxRepo := &mockOutboxRepo{
		buscarPendentesFn: func(ctx context.Context, limite int) ([]*entities.OutboxEvent, error) {
			return []*entities.OutboxEvent{evento}, nil
		},
		atualizarFn: func(ctx context.Context, event *entities.OutboxEvent) error {
			atualizouOutbox = true
			return nil
		},
	}

	notaRepo := &mockNotaRepo{
		obterPorIDFn: func(ctx context.Context, id uuid.UUID) (*entities.NotaFiscal, error) {
			return nota, nil
		},
		atualizarFn: func(ctx context.Context, nota *entities.NotaFiscal) error {
			atualizouNota = true
			return nil
		},
	}

	estoqueClient := &mockEstoqueClient{
		deduzirLoteFn: func(ctx context.Context, req ports.DeduzirLoteRequest) (*ports.DeduzirLoteResponse, error) {
			return &ports.DeduzirLoteResponse{
				Sucesso: true,
			}, nil
		},
	}

	uow := &mockUoW{
		executeFn: func(ctx context.Context, fn func(repos ports.TxRepositories) error) error {
			repos := ports.TxRepositories{
				Notas:  notaRepo,
				Outbox: outboxRepo,
			}
			return fn(repos)
		},
	}

	worker := newWorkerForTest(outboxRepo, notaRepo, estoqueClient, uow, logger)

	worker.ProcessarParaTeste(ctx)

	if !atualizouNota {
		t.Fatal("esperava atualização da nota ao fechar")
	}
	if !atualizouOutbox {
		t.Fatal("esperava atualização do outbox ao marcar processado")
	}
}