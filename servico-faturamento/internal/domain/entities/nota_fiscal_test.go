package entities_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nf-system/servico-faturamento/internal/domain/entities"
	"github.com/nf-system/servico-faturamento/internal/domain/enums"
)

func itensValidos() []entities.ItemNota {
	return []entities.ItemNota{
		{ProdutoID: uuid.New(), Quantidade: 2},
		{ProdutoID: uuid.New(), Quantidade: 5},
	}
}

func TestNovaNotaFiscal_ComItensValidos_CriaComStatusAberta(t *testing.T) {
	nota, err := entities.NovaNotaFiscal(itensValidos())

	if err != nil {
		t.Fatalf("esperava sem erro, got: %v", err)
	}
	if nota.Status != enums.StatusAberta {
		t.Errorf("esperava Aberta, got: %s", nota.Status)
	}
	if len(nota.Itens) != 2 {
		t.Errorf("esperava 2 itens, got: %d", len(nota.Itens))
	}
	if nota.ID == uuid.Nil {
		t.Error("ID não deve ser zero")
	}
}

func TestNovaNotaFiscal_SemItens_RetornaErro(t *testing.T) {
	_, err := entities.NovaNotaFiscal([]entities.ItemNota{})
	if err == nil {
		t.Error("esperava erro para nota sem itens")
	}
}

func TestNovaNotaFiscal_QuantidadeZero_RetornaErro(t *testing.T) {
	itens := []entities.ItemNota{{ProdutoID: uuid.New(), Quantidade: 0}}
	_, err := entities.NovaNotaFiscal(itens)
	if err == nil {
		t.Error("esperava erro para quantidade zero")
	}
}

func TestIniciarProcessamento_DeAberta_Transita(t *testing.T) {
	nota, _ := entities.NovaNotaFiscal(itensValidos())

	err := nota.IniciarProcessamento()

	if err != nil {
		t.Fatalf("esperava sem erro, got: %v", err)
	}
	if nota.Status != enums.StatusProcessando {
		t.Errorf("esperava Processando, got: %s", nota.Status)
	}
}

func TestIniciarProcessamento_DeProcessando_RetornaErro(t *testing.T) {
	nota, _ := entities.NovaNotaFiscal(itensValidos())
	_ = nota.IniciarProcessamento()

	err := nota.IniciarProcessamento()
	if err == nil {
		t.Error("esperava erro ao tentar processar nota já em Processando")
	}
}

func TestFechar_DeProcessando_Transita(t *testing.T) {
	nota, _ := entities.NovaNotaFiscal(itensValidos())
	_ = nota.IniciarProcessamento()

	err := nota.Fechar()

	if err != nil {
		t.Fatalf("esperava sem erro, got: %v", err)
	}
	if nota.Status != enums.StatusFechada {
		t.Errorf("esperava Fechada, got: %s", nota.Status)
	}
}

func TestFechar_DeAberta_RetornaErro(t *testing.T) {
	nota, _ := entities.NovaNotaFiscal(itensValidos())

	err := nota.Fechar()
	if err == nil {
		t.Error("esperava erro ao fechar nota Aberta diretamente")
	}
}

func TestCancelar_DeProcessando_TransitaComMotivo(t *testing.T) {
	nota, _ := entities.NovaNotaFiscal(itensValidos())
	_ = nota.IniciarProcessamento()
	produtoID := uuid.New()

	err := nota.Cancelar("saldo insuficiente", &produtoID)

	if err != nil {
		t.Fatalf("esperava sem erro, got: %v", err)
	}
	if nota.Status != enums.StatusCancelada {
		t.Errorf("esperava Cancelada, got: %s", nota.Status)
	}
	if nota.MotivoCancelamento == nil || *nota.MotivoCancelamento != "saldo insuficiente" {
		t.Error("motivo de cancelamento não gravado")
	}
	if nota.ProdutoFalhouID == nil || *nota.ProdutoFalhouID != produtoID {
		t.Error("produto que falhou não gravado")
	}
}

func TestPodeImprimir_ApenasQuandoAberta(t *testing.T) {
	nota, _ := entities.NovaNotaFiscal(itensValidos())

	if !nota.PodeImprimir() {
		t.Error("nota Aberta deve poder imprimir")
	}

	_ = nota.IniciarProcessamento()
	if nota.PodeImprimir() {
		t.Error("nota Processando não deve poder imprimir")
	}
}
