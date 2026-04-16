package enums

type StatusNota string

const (
	StatusAberta      StatusNota = "Aberta"
	StatusProcessando StatusNota = "Processando"
	StatusFechada     StatusNota = "Fechada"
	StatusCancelada   StatusNota = "Cancelada"
)

type StatusOutbox string

const (
	StatusOutboxPendente   StatusOutbox = "pendente"
	StatusOutboxProcessado StatusOutbox = "processado"
	StatusOutboxFalhou     StatusOutbox = "falhou"
)
