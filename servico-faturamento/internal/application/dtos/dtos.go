package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/nf-system/servico-faturamento/internal/domain/entities"
	"github.com/nf-system/servico-faturamento/internal/domain/enums"
)

// Requests
type CriarNotaRequest struct {
	Itens []ItemNotaRequest `json:"itens" binding:"required,min=1,dive"`
}

type ItemNotaRequest struct {
	ProdutoID        uuid.UUID `json:"produtoId" binding:"required"`
	ProdutoCodigo    string    `json:"produtoCodigo" binding:"required"`
	ProdutoDescricao string    `json:"produtoDescricao" binding:"required"`
	Quantidade       float64   `json:"quantidade" binding:"required,gt=0"`
}

type ImprimirNotaRequest struct {
	// sem body — a idempotency key vem no header
}

// Responses
type NotaFiscalResponse struct {
	ID                 uuid.UUID          `json:"id"`
	Numero             int64              `json:"numero"`
	Status             enums.StatusNota   `json:"status"`
	Itens              []ItemNotaResponse `json:"itens"`
	CriadoEm           time.Time          `json:"criadoEm"`
	AtualizadoEm       time.Time          `json:"atualizadoEm"`
	MotivoCancelamento *string            `json:"motivoCancelamento,omitempty"`
	ProdutoFalhouID    *uuid.UUID         `json:"produtoFalhouId,omitempty"`
}

type ItemNotaResponse struct {
	ID               uuid.UUID `json:"id"`
	ProdutoID        uuid.UUID `json:"produtoId"`
	ProdutoCodigo    string    `json:"produtoCodigo"`
	ProdutoDescricao string    `json:"produtoDescricao"`
	Quantidade       float64   `json:"quantidade"`
}

type ImprimirResponse struct {
	NotaID   uuid.UUID        `json:"notaId"`
	Status   enums.StatusNota `json:"status"`
	Mensagem string           `json:"mensagem"`
}

type ListarNotasQuery struct {
	Page          int    `form:"page"`
	PageSize      int    `form:"pageSize"`
	Search        string `form:"search"`
	SortBy        string `form:"sortBy"`
	SortDirection string `form:"sortDirection"`
}

type NotaFiscalPagedResponse struct {
	Items      []NotaFiscalResponse `json:"items"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"pageSize"`
	TotalItems int64                `json:"totalItems"`
	TotalPages int                  `json:"totalPages"`
}

type PagedResponse[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

// Mappers
func ToNotaResponse(n *entities.NotaFiscal) NotaFiscalResponse {
	itens := make([]ItemNotaResponse, len(n.Itens))
	for i, item := range n.Itens {
		itens[i] = ItemNotaResponse{
			ID:               item.ID,
			ProdutoID:        item.ProdutoID,
			ProdutoCodigo:    item.ProdutoCodigo,
			ProdutoDescricao: item.ProdutoDescricao,
			Quantidade:       item.Quantidade,
		}
	}

	return NotaFiscalResponse{
		ID:                 n.ID,
		Numero:             n.Numero,
		Status:             n.Status,
		Itens:              itens,
		CriadoEm:           n.CriadoEm,
		AtualizadoEm:       n.AtualizadoEm,
		MotivoCancelamento: n.MotivoCancelamento,
		ProdutoFalhouID:    n.ProdutoFalhouID,
	}
}