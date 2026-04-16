using Estoque.Domain.Enums;

namespace Estoque.Application.DTOs;

public record CriarProdutoRequest(
    string Codigo,
    string Descricao,
    decimal SaldoInicial = 0);

public record AtualizarProdutoRequest(
    string Descricao);

public record DeduzirEstoqueRequest(
    Guid ProdutoId,
    decimal Quantidade,
    Guid? NotaFiscalId = null);

public record EstornarEstoqueRequest(
    Guid ProdutoId,
    decimal Quantidade,
    Guid? NotaFiscalId = null);

public record AdicionarEntradaRequest(
    Guid ProdutoId,
    decimal Quantidade,
    Guid? NotaFiscalId = null);

// Payload para dedução em lote (vindo do serviço de faturamento)
public record DeduzirLoteRequest(
    Guid NotaFiscalId,
    IEnumerable<ItemDeduzirRequest> Itens);

public record ItemDeduzirRequest(
    Guid ProdutoId,
    decimal Quantidade);

public record EstornarLoteRequest(
    Guid NotaFiscalId,
    IEnumerable<ItemDeduzirRequest> Itens);

// Responses
public record ProdutoResponse(
    Guid Id,
    string Codigo,
    string Descricao,
    decimal Saldo,
    DateTime CriadoEm,
    DateTime AtualizadoEm);
public record PagedResult<T>(
    IEnumerable<T> Items,
    int Page,
    int PageSize,
    int TotalItems,
    int TotalPages);

public record MovimentacaoResponse(
    Guid Id,
    Guid ProdutoId,
    TipoMovimentacao Tipo,
    decimal Quantidade,
    decimal SaldoAnterior,
    decimal SaldoResultante,
    Guid? NotaFiscalId,
    bool IsEstorno,
    DateTime OcorridoEm);

public record DeduzirLoteResponse(
    bool Sucesso,
    IEnumerable<ItemDeduzidoResponse> ItensDeduzidos,
    string? Erro = null,
    Guid? ProdutoFalhou = null);

public record ItemDeduzidoResponse(
    Guid ProdutoId,
    decimal Quantidade);
