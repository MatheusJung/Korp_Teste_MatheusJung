package usecases

import (
	"context"
	"math"
	"time"

	"github.com/nf-system/servico-faturamento/internal/application/ports"
	"github.com/nf-system/servico-faturamento/internal/domain/entities"
	"github.com/nf-system/servico-faturamento/internal/domain/enums"
	"go.uber.org/zap"
)

type OutboxWorker struct {
	outbox  ports.OutboxRepository
	notas   ports.NotaFiscalRepository
	estoque ports.EstoqueClient
	uow     ports.UnitOfWork
	pdfService ports.NotaPDFService
	pdfStorage ports.NotaPDFStorage
	logger  *zap.Logger

	pollInterval  time.Duration
	maxTentativas int
}

func NewOutboxWorker(
	outbox ports.OutboxRepository,
	notas ports.NotaFiscalRepository,
	estoque ports.EstoqueClient,
	uow ports.UnitOfWork,
	pdfService ports.NotaPDFService,
	pdfStorage ports.NotaPDFStorage,
	pollInterval time.Duration,
	maxTentativas int,
	logger *zap.Logger,
) *OutboxWorker {
	return &OutboxWorker{
		outbox:        outbox,
		notas:         notas,
		estoque:       estoque,
		uow:           uow,
		pdfService:    pdfService,
		pdfStorage:    pdfStorage,
		pollInterval:  pollInterval,
		maxTentativas: maxTentativas,
		logger:        logger,
	}
}

// Start inicia o worker em background. Respeita context cancellation para shutdown gracioso.
func (w *OutboxWorker) Start(ctx context.Context) {
	w.logger.Info("outbox worker iniciado", zap.Duration("pollInterval", w.pollInterval))

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("outbox worker encerrado")
			return
		case <-ticker.C:
			w.processar(ctx)
		}
	}
}

// ProcessarParaTeste expõe o ciclo de processamento para testes unitários sem ticker.
func (w *OutboxWorker) ProcessarParaTeste(ctx context.Context) {
	w.processar(ctx)
}

func (w *OutboxWorker) processar(ctx context.Context) {
	eventos, err := w.outbox.BuscarPendentes(ctx, 10)
	if err != nil {
		w.logger.Error("erro ao buscar eventos pendentes", zap.Error(err))
		return
	}

	for _, evento := range eventos {
		w.processarEvento(ctx, evento)
	}
}

func (w *OutboxWorker) processarEvento(ctx context.Context, evento *entities.OutboxEvent) {
	log := w.logger.With(
		zap.String("eventoId", evento.ID.String()),
		zap.String("notaId", evento.NotaFiscalID.String()),
		zap.Int("tentativa", evento.Tentativas+1),
	)

	nota, err := w.notas.ObterPorID(ctx, evento.NotaFiscalID)
	if err != nil {
		log.Error("nota não encontrada para evento outbox", zap.Error(err))
		return
	}

	req := ports.DeduzirLoteRequest{
		NotaFiscalID: nota.ID,
	}

	for _, item := range evento.Payload.Itens {
		req.Itens = append(req.Itens, ports.ItemEstoque{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
		})
	}

	resultado, err := w.estoque.DeduzirLote(ctx, req)
	if err != nil {
		log.Warn("erro técnico ao chamar estoque, reagendando evento", zap.Error(err))
		w.reagendar(ctx, evento, log)
		return
	}

	if resultado == nil {
		log.Error("resultado do estoque veio nil sem erro, reagendando por segurança")
		w.reagendar(ctx, evento, log)
		return
	}

	if resultado.Sucesso {
		w.fecharNota(ctx, nota, evento, log)
		return
	}

	log.Warn("erro de negócio ao deduzir estoque, cancelando nota",
		zap.String("erro", resultado.Erro),
		zap.Any("produtoFalhou", resultado.ProdutoFalhou),
	)
	w.cancelarNota(ctx, nota, evento, resultado, log)
}

func (w *OutboxWorker) fecharNota(
	ctx context.Context,
	nota *entities.NotaFiscal,
	evento *entities.OutboxEvent,
	log *zap.Logger,
) {
	err := w.uow.Execute(ctx, func(repos ports.TxRepositories) error {
		if err := nota.Fechar(); err != nil {
			return err
		}
		if err := repos.Notas.Atualizar(ctx, nota); err != nil {
			return err
		}

		evento.Marcar(enums.StatusOutboxProcessado)
		return repos.Outbox.Atualizar(ctx, evento)
	})
	if err != nil {
		log.Error("erro ao fechar nota", zap.Error(err))
		return
	}

	if w.pdfService != nil && w.pdfStorage != nil {
		pdfBytes, err := w.pdfService.Gerar(ctx, nota)
		if err != nil {
			log.Error("erro ao gerar pdf da nota", zap.Error(err))
			return
		}

		if err := w.pdfStorage.Salvar(ctx, nota.ID, nota.Numero, pdfBytes); err != nil {
			log.Error("erro ao salvar pdf da nota", zap.Error(err))
			return
		}

		log.Info("pdf da nota gerado e salvo com sucesso")
	}

	log.Info("nota fechada com sucesso")
}

func (w *OutboxWorker) cancelarNota(
	ctx context.Context,
	nota *entities.NotaFiscal,
	evento *entities.OutboxEvent,
	resultado *ports.DeduzirLoteResponse,
	log *zap.Logger,
) {
	err := w.uow.Execute(ctx, func(repos ports.TxRepositories) error {
		if err := nota.Cancelar(resultado.Erro, resultado.ProdutoFalhou); err != nil {
			return err
		}
		if err := repos.Notas.Atualizar(ctx, nota); err != nil {
			return err
		}

		evento.Marcar(enums.StatusOutboxFalhou)
		return repos.Outbox.Atualizar(ctx, evento)
	})
	if err != nil {
		log.Error("erro ao cancelar nota", zap.Error(err))
		return
	}

	log.Info("nota cancelada por erro de negócio no estoque")
}

func (w *OutboxWorker) reagendar(ctx context.Context, evento *entities.OutboxEvent, log *zap.Logger) {
	if evento.EsgotouTentativas() {
		log.Error("evento esgotou tentativas, marcando como falhou")

		nota, err := w.notas.ObterPorID(ctx, evento.NotaFiscalID)
		if err != nil {
			log.Error("erro ao buscar nota para cancelamento após esgotar tentativas", zap.Error(err))
			return
		}

		motivo := "Serviço de estoque indisponível após múltiplas tentativas"

		err = w.uow.Execute(ctx, func(repos ports.TxRepositories) error {
			if err := nota.Cancelar(motivo, nil); err != nil {
				return err
			}
			if err := repos.Notas.Atualizar(ctx, nota); err != nil {
				return err
			}

			evento.Marcar(enums.StatusOutboxFalhou)
			return repos.Outbox.Atualizar(ctx, evento)
		})
		if err != nil {
			log.Error("erro ao marcar evento como falho após esgotar tentativas", zap.Error(err))
			return
		}

		log.Info("evento marcado como falho após esgotar tentativas")
		return
	}

	backoff := time.Duration(math.Pow(2, float64(evento.Tentativas))) * 5 * time.Second
	proxima := time.Now().UTC().Add(backoff)

	evento.RegistrarTentativa(proxima)

	if err := w.outbox.Atualizar(ctx, evento); err != nil {
		log.Error("erro ao reagendar evento", zap.Error(err))
		return
	}

	log.Info("evento reagendado",
		zap.Time("proximaTentativa", proxima),
		zap.Duration("backoff", backoff),
	)
}