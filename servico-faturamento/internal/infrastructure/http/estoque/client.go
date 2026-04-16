package estoque

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nf-system/servico-faturamento/internal/application/ports"
)

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

type deduzirLotePayload struct {
	NotaFiscalID uuid.UUID    `json:"notaFiscalId"`
	Itens        []itemPayload `json:"itens"`
}

type estornarLotePayload struct {
	NotaFiscalID uuid.UUID    `json:"notaFiscalId"`
	Itens        []itemPayload `json:"itens"`
}

type itemPayload struct {
	ProdutoID  uuid.UUID `json:"produtoId"`
	Quantidade float64   `json:"quantidade"`
}

type deduzirLoteResponse struct {
	Sucesso        bool          `json:"sucesso"`
	ItensDeduzidos []itemPayload `json:"itensDeduzidos"`
	Erro           string        `json:"erro"`
	ProdutoFalhou  *uuid.UUID    `json:"produtoFalhou"`
}

func (c *HTTPClient) DeduzirLote(ctx context.Context, req ports.DeduzirLoteRequest) (*ports.DeduzirLoteResponse, error) {
	payload := deduzirLotePayload{NotaFiscalID: req.NotaFiscalID}
	for _, item := range req.Itens {
		payload.Itens = append(payload.Itens, itemPayload{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
		})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/movimentacoes/deduzir-lote",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("estoque indisponível: %w", err)
	}
	defer resp.Body.Close()

	var result deduzirLoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta do estoque: %w", err)
	}

	out := &ports.DeduzirLoteResponse{
		Sucesso:       result.Sucesso,
		Erro:          result.Erro,
		ProdutoFalhou: result.ProdutoFalhou,
	}
	for _, item := range result.ItensDeduzidos {
		out.ItensDeduzidos = append(out.ItensDeduzidos, ports.ItemEstoque{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
		})
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return out, nil

	case http.StatusConflict:
		if out.Erro == "" {
			out.Erro = "conflito ao deduzir estoque"
		}
		out.Sucesso = false
		return out, nil

	default:
		return nil, fmt.Errorf("estoque retornou status inesperado: %d", resp.StatusCode)
	}
}

func (c *HTTPClient) EstornarLote(ctx context.Context, req ports.EstornarLoteRequest) error {
	payload := estornarLotePayload{NotaFiscalID: req.NotaFiscalID}
	for _, item := range req.Itens {
		payload.Itens = append(payload.Itens, itemPayload{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
		})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/api/movimentacoes/estornar-lote", bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("estoque indisponível: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("estoque retornou status inesperado no estorno: %d", resp.StatusCode)
	}
	return nil
}

func (c *HTTPClient) HealthCheck(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("estoque indisponível: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("estoque health check retornou: %d", resp.StatusCode)
	}
	return nil
}
