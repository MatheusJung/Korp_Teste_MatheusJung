package pdf

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/nf-system/servico-faturamento/internal/domain/entities"
)

type FPDFService struct{}

func NewFPDFService() *FPDFService {
	return &FPDFService{}
}

func (s *FPDFService) Gerar(ctx context.Context, nota *entities.NotaFiscal) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(fmt.Sprintf("DANFE %d", nota.Numero), false)
	pdf.SetAuthor("Servico Faturamento", false)
	pdf.SetMargins(8, 8, 8)
	pdf.SetAutoPageBreak(true, 12)
	pdf.AddPage()

	pageW, _ := pdf.GetPageSize()
	left, _, right, _ := pdf.GetMargins()
	contentW := pageW - left - right

	drawSectionTitle := func(title string) {
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(0, 6, title, "1", 1, "L", false, 0, "")
	}

	drawField := func(x, y, w float64, label, value string) {
		pdf.SetXY(x, y)
		pdf.SetFont("Arial", "B", 8)
		pdf.CellFormat(w, 4, label, "LTR", 1, "L", false, 0, "")

		pdf.SetX(x)
		pdf.SetFont("Arial", "", 9)
		pdf.CellFormat(w, 6, value, "LBR", 0, "L", false, 0, "")
	}

	// Cabeçalho
	leftCol := contentW * 0.62
	rightCol := contentW - leftCol
	startY := pdf.GetY()

	pdf.SetXY(left, startY)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(leftCol, 10, "DANFE SIMPLIFICADA", "1", 1, "C", false, 0, "")

	pdf.SetX(left)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(leftCol, 5,
		"Documento Auxiliar da Nota Fiscal\nGerado pelo sistema de faturamento",
		"1", "C", false,
	)

	pdf.SetXY(left+leftCol, startY)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(rightCol, 8, "CONTROLE DO DOCUMENTO", "1", 1, "C", false, 0, "")

	pdf.SetX(left + leftCol)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(rightCol, 6, fmt.Sprintf("Numero da nota: %d", nota.Numero), "1", 1, "L", false, 0, "")

	pdf.SetX(left + leftCol)
	pdf.CellFormat(rightCol, 6, fmt.Sprintf("Status: %s", nota.Status), "1", 1, "L", false, 0, "")

	pdf.SetX(left + leftCol)
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(rightCol, 6, fmt.Sprintf("Codigo da nota: NF-%06d", nota.Numero), "1", 1, "L", false, 0, "")

	if pdf.GetY() < startY+22 {
		pdf.SetY(startY + 22)
	}
	pdf.Ln(2)

	// Emitente
	drawSectionTitle("EMITENTE")
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 5,
		"NF System LTDA\n"+
			"Rua Exemplo, 123 - Centro\n"+
			"Cidade/UF - CEP 00000-000\n"+
			"CNPJ: 00.000.000/0001-00",
		"1", "L", false,
	)

	pdf.Ln(1)

	// Dados gerais - corrigido para não quebrar
	y := pdf.GetY()
	col1 := contentW * 0.34
	col2 := contentW * 0.33
	col3 := contentW - col1 - col2

	drawField(left, y, col1, "DATA DE EMISSAO", time.Now().Format("02/01/2006 15:04"))
	drawField(left+col1, y, col2, "TIPO DE OPERACAO", "SAIDA")
	drawField(left+col1+col2, y, col3, "QTD ITENS", fmt.Sprintf("%d", len(nota.Itens)))

	pdf.SetY(y + 10)
	pdf.Ln(2)

	// Itens
	drawSectionTitle("DADOS DOS PRODUTOS / SERVICOS")

	codeW := contentW * 0.20
	descW := contentW * 0.58
	qtdW := contentW * 0.22

	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(codeW, 7, "CODIGO", "1", 0, "C", false, 0, "")
	pdf.CellFormat(descW, 7, "DESCRICAO", "1", 0, "C", false, 0, "")
	pdf.CellFormat(qtdW, 7, "QUANTIDADE", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 8)

	var totalQtd float64

	for _, item := range nota.Itens {
		codigo := item.ProdutoCodigo
		descricao := item.ProdutoDescricao
		if codigo == "" {
			codigo = "SEM-CODIGO"
		}
		if descricao == "" {
			descricao = "Produto sem descricao"
		}

		lineH := 6.0
		x := pdf.GetX()
		y := pdf.GetY()

		// calcula altura da descrição
		pdf.SetXY(x+codeW, y)
		linesDesc := pdf.SplitLines([]byte(descricao), descW)
		rowH := float64(len(linesDesc)) * lineH
		if rowH < lineH {
			rowH = lineH
		}

		// código
		pdf.SetXY(x, y)
		pdf.CellFormat(codeW, rowH, codigo, "1", 0, "L", false, 0, "")

		// descrição
		pdf.SetXY(x+codeW, y)
		pdf.MultiCell(descW, lineH, descricao, "1", "L", false)

		// quantidade
		pdf.SetXY(x+codeW+descW, y)
		pdf.CellFormat(qtdW, rowH, fmt.Sprintf("%.2f", item.Quantidade), "1", 1, "R", false, 0, "")

		totalQtd += item.Quantidade
	}

	pdf.Ln(2)

	// Total
	totalBoxW := 72.0
	pdf.SetX(pageW - right - totalBoxW)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(totalBoxW, 7, "TOTAL GERAL", "1", 1, "C", false, 0, "")

	pdf.SetX(pageW - right - totalBoxW)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(totalBoxW, 8, fmt.Sprintf("Quantidade total: %.2f", totalQtd), "1", 1, "R", false, 0, "")

	pdf.Ln(3)

	// Observações
	drawSectionTitle("INFORMACOES COMPLEMENTARES")
	pdf.SetFont("Arial", "", 8)
	pdf.MultiCell(0, 5,
		"Documento gerado eletronicamente pelo servico de faturamento.\n"+
			fmt.Sprintf("Codigo da nota: NF-%06d", nota.Numero),
		"1", "L", false,
	)

	// Rodapé
	pdf.SetY(280)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 5,
		fmt.Sprintf("Gerado em %s", time.Now().Format("02/01/2006 15:04:05")),
		"T", 0, "R", false, 0, "",
	)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}