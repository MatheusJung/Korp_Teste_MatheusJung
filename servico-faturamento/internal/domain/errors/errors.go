package errors

import (
	"errors"
	"fmt"

	"github.com/nf-system/servico-faturamento/internal/domain/enums"
)

var (
	ErrNotaSemItens        = errors.New("a nota fiscal deve ter ao menos um item")
	ErrQuantidadeInvalida  = errors.New("quantidade dos itens deve ser positiva")
	ErrNotaNaoEncontrada   = errors.New("nota fiscal não encontrada")
	ErrNotaNaoPodeImprimir = errors.New("apenas notas com status 'Aberta' podem ser impressas")
	ErrIdempotenciaRepetida = errors.New("operação já processada anteriormente")
)

type TransicaoInvalidaError struct {
	De  enums.StatusNota
	Para enums.StatusNota
}

func (e *TransicaoInvalidaError) Error() string {
	return fmt.Sprintf("transição de status inválida: %s → %s", e.De, e.Para)
}

func ErrTransicaoInvalida(de, para enums.StatusNota) error {
	return &TransicaoInvalidaError{De: de, Para: para}
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotaNaoEncontrada)
}

func IsIdempotencia(err error) bool {
	return errors.Is(err, ErrIdempotenciaRepetida)
}
