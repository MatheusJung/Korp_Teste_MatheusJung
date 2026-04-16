package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/nf-system/servico-faturamento/internal/domain/enums"
	"github.com/nf-system/servico-faturamento/internal/domain/errors"
)

type NotaFiscal struct {
	ID             uuid.UUID
	Numero         int64
	Status         enums.StatusNota
	Itens          []ItemNota
	CriadoEm      time.Time
	AtualizadoEm  time.Time
	MotivoCancelamento *string
	ProdutoFalhouID    *uuid.UUID
}

type ItemNota struct {
	ID               uuid.UUID
	NotaID           uuid.UUID
	ProdutoID        uuid.UUID
	ProdutoCodigo    string
	ProdutoDescricao string
	Quantidade       float64
}

func NovaNotaFiscal(itens []ItemNota) (*NotaFiscal, error) {
	if len(itens) == 0 {
		return nil, errors.ErrNotaSemItens
	}

	for _, item := range itens {
		if item.Quantidade <= 0 {
			return nil, errors.ErrQuantidadeInvalida
		}
	}

	now := time.Now().UTC()
	nota := &NotaFiscal{
		ID:            uuid.New(),
		Status:        enums.StatusAberta,
		Itens:         itens,
		CriadoEm:     now,
		AtualizadoEm: now,
	}

	// atribui IDs aos itens
	for i := range nota.Itens {
		nota.Itens[i].ID = uuid.New()
		nota.Itens[i].NotaID = nota.ID
	}

	return nota, nil
}

func (n *NotaFiscal) IniciarProcessamento() error {
	if n.Status != enums.StatusAberta {
		return errors.ErrTransicaoInvalida(n.Status, enums.StatusProcessando)
	}
	n.Status = enums.StatusProcessando
	n.AtualizadoEm = time.Now().UTC()
	return nil
}

func (n *NotaFiscal) Fechar() error {
	if n.Status != enums.StatusProcessando {
		return errors.ErrTransicaoInvalida(n.Status, enums.StatusFechada)
	}
	n.Status = enums.StatusFechada
	n.AtualizadoEm = time.Now().UTC()
	return nil
}

func (n *NotaFiscal) Cancelar(motivo string, produtoFalhouID *uuid.UUID) error {
	if n.Status != enums.StatusProcessando {
		return errors.ErrTransicaoInvalida(n.Status, enums.StatusCancelada)
	}
	n.Status = enums.StatusCancelada
	n.MotivoCancelamento = &motivo
	n.ProdutoFalhouID = produtoFalhouID
	n.AtualizadoEm = time.Now().UTC()
	return nil
}

func (n *NotaFiscal) PodeImprimir() bool {
	return n.Status == enums.StatusAberta
}
