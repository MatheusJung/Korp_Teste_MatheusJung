// Genericos
export interface PagedResult<T> {
  items: T[];
  page: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
}

export interface ListarParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortDirection?: "asc" | "desc";
}

// ── Estoque ────────────────────────────────────────────────────────────────────

export interface Produto {
  id: string;
  codigo: string;
  descricao: string;
  saldo: number;
  criadoEm: string;
  atualizadoEm: string;
}

export interface CriarProdutoRequest {
  codigo: string;
  descricao: string;
  saldoInicial: number;
}

export interface Movimentacao {
  id: string;
  produtoId: string;
  tipo: "0" | "1" | "2";
  quantidade: number;
  saldoAnterior: number;
  saldoResultante: number;
  notaFiscalId?: string;
  isEstorno: boolean;
  ocorridoEm: string;
}

export interface AdicionarEntradaRequest {
  produtoId: string;
  quantidade: number;
}

// ── Faturamento ────────────────────────────────────────────────────────────────

export type StatusNota = "Aberta" | "Processando" | "Fechada" | "Cancelada";

export interface NotaFiscal {
  id: string;
  numero: number;
  status: StatusNota;
  itens: ItemNota[];
  criadoEm: string;
  atualizadoEm: string;
  motivoCancelamento?: string;
  produtoFalhouId?: string;
}

export interface ItemNota {
  id: string;
  produtoId: string;
  quantidade: number;
}

export interface CriarNotaRequest {
  itens: {
    produtoId: string;
    produtoCodigo: string;
    produtoDescricao: string;
    quantidade: number;
  }[];
}

export interface ImprimirResponse {
  notaId: string;
  status: StatusNota;
  mensagem: string;
}

// ── Health ─────────────────────────────────────────────────────────────────────

export type ServiceStatus = "ok" | "degraded" | "down";

export interface HealthResponse {
  status: ServiceStatus;
  checks: Record<string, string>;
  timestamp: string;
}
