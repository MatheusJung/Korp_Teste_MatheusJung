package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/nf-system/servico-faturamento/internal/domain/enums"
)

type OutboxEvent struct {
	ID              uuid.UUID
	NotaFiscalID    uuid.UUID
	Payload         OutboxPayload
	Status          enums.StatusOutbox
	Tentativas      int
	MaxTentativas   int
	CriadoEm       time.Time
	ProcessadoEm   *time.Time
	ProximaTentativa time.Time
}

type OutboxPayload struct {
	NotaFiscalID   uuid.UUID        `json:"notaFiscalId"`
	Itens          []ItemOutbox     `json:"itens"`
	ItensDeduzidos []ItemOutbox     `json:"itensDeduzidos"`
}

type ItemOutbox struct {
	ProdutoID  uuid.UUID `json:"produtoId"`
	Quantidade float64   `json:"quantidade"`
}

func NovoOutboxEvent(notaID uuid.UUID, itens []ItemNota, maxTentativas int) *OutboxEvent {
	payload := OutboxPayload{
		NotaFiscalID:   notaID,
		ItensDeduzidos: []ItemOutbox{},
	}
	for _, item := range itens {
		payload.Itens = append(payload.Itens, ItemOutbox{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
		})
	}

	return &OutboxEvent{
		ID:               uuid.New(),
		NotaFiscalID:     notaID,
		Payload:          payload,
		Status:           enums.StatusOutboxPendente,
		Tentativas:       0,
		MaxTentativas:    maxTentativas,
		CriadoEm:        time.Now().UTC(),
		ProximaTentativa: time.Now().UTC(),
	}
}

func (e *OutboxEvent) RegistrarTentativa(proximaEm time.Time) {
	e.Tentativas++
	e.ProximaTentativa = proximaEm
}

func (e *OutboxEvent) Marcar(status enums.StatusOutbox) {
	e.Status = status
	if status == enums.StatusOutboxProcessado {
		now := time.Now().UTC()
		e.ProcessadoEm = &now
	}
}

func (e *OutboxEvent) EsgotouTentativas() bool {
	return e.Tentativas >= e.MaxTentativas
}
