package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nf-system/servico-faturamento/internal/application/dtos"
	"github.com/nf-system/servico-faturamento/internal/application/ports"
	"github.com/nf-system/servico-faturamento/internal/application/usecases"
	domainerrors "github.com/nf-system/servico-faturamento/internal/domain/errors"
)

type NotaFiscalHandler struct {
	criar      *usecases.CriarNotaUseCase
	obter      *usecases.ObterNotaUseCase
	listar     *usecases.ListarNotasUseCase
	listarPaginado     *usecases.ListarNotasPaginadoUseCase
	imprimir   *usecases.ImprimirNotaUseCase
	pdfStorage ports.NotaPDFStorage
}

func NewNotaFiscalHandler(
	criar *usecases.CriarNotaUseCase,
	obter *usecases.ObterNotaUseCase,
	listar *usecases.ListarNotasUseCase,
	listarPaginado *usecases.ListarNotasPaginadoUseCase,
	imprimir *usecases.ImprimirNotaUseCase,
	pdfStorage ports.NotaPDFStorage,
) *NotaFiscalHandler {
	return &NotaFiscalHandler{
		criar:      criar,
		obter:      obter,
		listar:     listar,
		listarPaginado:		listarPaginado,
		imprimir:   imprimir,
		pdfStorage: pdfStorage,
	}
}

func (h *NotaFiscalHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/notas")
	{
		api.GET("", h.Listar)
		api.GET("/paginado", h.ListarPaginado)
		api.GET("/:id", h.Obter)
		api.GET("/:id/pdf", h.ObterPDF)
		api.POST("", h.Criar)
		api.POST("/:id/imprimir", h.Imprimir)
	}
}

// Listar godoc
// @Summary     Lista notas fiscais
// @Description Retorna todas as notas fiscais cadastradas
// @Tags        notas
// @Produce     json
// @Success     200 {array}  dtos.NotaFiscalResponse
// @Failure     500 {object} map[string]string
// @Router      /notas [get]
func (h *NotaFiscalHandler) Listar(c *gin.Context) {
	notas, err := h.listar.Executar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notas)
}

// ListarPaginado godoc
// @Summary     Lista notas fiscais com paginação
// @Description Retorna notas fiscais com paginação, pesquisa e ordenação
// @Tags        notas
// @Produce     json
// @Param       page     query int    false "Página" default(1)
// @Param       pageSize query int    false "Tamanho da página" default(10)
// @Param       search   query string false "Pesquisa por número, status ou produto"
// @Param       sortBy   query string false "Campo de ordenação" Enums(numero,status,criado_em,atualizado_em) default(criado_em)
// @Param       sortDir  query string false "Direção da ordenação" Enums(asc,desc) default(desc)
// @Success     200 {object} dtos.NotaFiscalPagedResponse
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /notas/paginado [get]
func (h *NotaFiscalHandler) ListarPaginado(c *gin.Context) {
	var query dtos.ListarNotasQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "parâmetros inválidos"})
		return
	}

	result, err := h.listarPaginado.Executar(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Obter godoc
// @Summary     Obtém nota fiscal
// @Description Retorna uma nota fiscal pelo ID
// @Tags        notas
// @Produce     json
// @Param       id  path     string true "ID da nota fiscal"
// @Success     200 {object} dtos.NotaFiscalResponse
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /notas/{id} [get]
func (h *NotaFiscalHandler) Obter(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id inválido"})
		return
	}

	nota, err := h.obter.Executar(c.Request.Context(), id)
	if err != nil {
		if domainerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nota)
}

// Criar godoc
// @Summary     Cria nota fiscal
// @Description Cria uma nova nota fiscal com status Aberta
// @Tags        notas
// @Accept      json
// @Produce     json
// @Param       nota body     dtos.CriarNotaRequest true "Dados da nota"
// @Success     201  {object} dtos.NotaFiscalResponse
// @Failure     400  {object} map[string]string
// @Failure     422  {object} map[string]string
// @Router      /notas [post]
func (h *NotaFiscalHandler) Criar(c *gin.Context) {
	var req dtos.CriarNotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	nota, err := h.criar.Executar(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, nota)
}

// Imprimir godoc
// @Summary     Imprime nota fiscal
// @Description Inicia o processo de impressão. Retorna 202 enquanto o estoque é deduzido em background.
// @Tags        notas
// @Produce     json
// @Param       id              path     string true  "ID da nota fiscal"
// @Param       Idempotency-Key header   string true  "Chave única para idempotência (UUID)"
// @Success     202             {object} dtos.ImprimirResponse
// @Success     200             {object} dtos.ImprimirResponse "Retornado quando a chave já foi processada"
// @Failure     400             {object} map[string]string
// @Failure     404             {object} map[string]string
// @Failure     422             {object} map[string]string
// @Router      /notas/{id}/imprimir [post]
func (h *NotaFiscalHandler) Imprimir(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id inválido"})
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "header Idempotency-Key é obrigatório"})
		return
	}

	resp, err := h.imprimir.Executar(c.Request.Context(), id, idempotencyKey)
	if err != nil {
		if domainerrors.IsIdempotencia(err) {
			c.JSON(http.StatusOK, resp)
			return
		}
		if domainerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, resp)
}

// ObterPDF godoc
// @Summary     Obtém o PDF da nota fiscal
// @Description Retorna o arquivo PDF da nota fiscal gerado após o processamento
// @Tags        notas
// @Produce     application/pdf
// @Param       id   path      string true "ID da nota fiscal"
// @Success     200  {file}    file
// @Failure     400  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Router      /notas/{id}/pdf [get]
func (h *NotaFiscalHandler) ObterPDF(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id inválido"})
		return
	}

	//BUSCA A NOTA PRIMEIRO
	nota, err := h.obter.Executar(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "nota não encontrada"})
		return
	}

	//AGORA USA O NUMERO
	pdfBytes, nomeArquivo, err := h.pdfStorage.Obter(c.Request.Context(), nota.ID, nota.Numero)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", `inline; filename="`+nomeArquivo+`"`)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}