using Estoque.Domain.Entities;

namespace Estoque.Application.Interfaces;

public interface IProdutoRepository
{
    Task<Produto?> ObterPorIdAsync(Guid id, CancellationToken ct = default);
    Task<Produto?> ObterPorCodigoAsync(string codigo, CancellationToken ct = default);
    Task<IEnumerable<Produto>> ListarAsync(CancellationToken ct = default);
    Task<(IEnumerable<Produto> Items, int TotalItems)> ListarPaginadoAsync(
        int page,
        int pageSize,
        string? search,
        string? sortBy,
        string? sortDirection,
        CancellationToken ct = default);
    Task AdicionarAsync(Produto produto, CancellationToken ct = default);
    Task<bool> ExisteCodigoAsync(string codigo, Guid? ignorarId = null, CancellationToken ct = default);
    Task<Produto?> ObterSemMovimentacoesAsync(Guid id, CancellationToken ct = default);
    Task AtualizarAsync(Produto produto, CancellationToken ct = default);
}

public interface IMovimentacaoRepository
{
    Task<IEnumerable<Movimentacao>> ListarPorProdutoAsync(Guid produtoId, CancellationToken ct = default);
}

public interface IUnitOfWork
{
    Task<int> CommitAsync(CancellationToken ct = default);
}