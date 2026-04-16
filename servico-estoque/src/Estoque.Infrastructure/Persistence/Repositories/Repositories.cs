using Estoque.Application.Interfaces;
using Estoque.Domain.Entities;
using Microsoft.Data.SqlClient;
using Microsoft.EntityFrameworkCore;

namespace Estoque.Infrastructure.Persistence.Repositories;

public class ProdutoRepository(EstoqueDbContext context) : IProdutoRepository
{
    public async Task<Produto?> ObterPorIdAsync(Guid id, CancellationToken ct = default) =>
        await context.Produtos
            .Include(p => p.Movimentacoes)
            .FirstOrDefaultAsync(p => p.Id == id, ct);

    public async Task<Produto?> ObterPorCodigoAsync(string codigo, CancellationToken ct = default) =>
        await context.Produtos
            .Include(p => p.Movimentacoes)
            .FirstOrDefaultAsync(p => p.Codigo == codigo.ToUpper(), ct);

    public async Task<Produto?> ObterSemMovimentacoesAsync(Guid id, CancellationToken ct = default) =>
    await context.Produtos
        .AsNoTracking()
        .FirstOrDefaultAsync(p => p.Id == id, ct);

    public async Task<IEnumerable<Produto>> ListarAsync(CancellationToken ct = default) =>
        await context.Produtos
            .OrderBy(p => p.Codigo)
            .ToListAsync(ct);
    public async Task<(IEnumerable<Produto> Items, int TotalItems)> ListarPaginadoAsync(
        int page,
        int pageSize,
        string? search,
        string? sortBy,
        string? sortDirection,
        CancellationToken ct = default)
    {
        var query = context.Produtos.AsNoTracking();

        if (!string.IsNullOrWhiteSpace(search))
        {
            query = query.Where(p =>
                p.Codigo.Contains(search) ||
                p.Descricao.Contains(search));
        }

        query = (sortBy, sortDirection) switch
        {
            ("codigo", "desc") => query.OrderByDescending(p => p.Codigo),
            ("descricao", "desc") => query.OrderByDescending(p => p.Descricao),
            ("saldo", "desc") => query.OrderByDescending(p => p.Saldo),
            ("criadoem", "desc") => query.OrderByDescending(p => p.CriadoEm),
            ("atualizadoem", "desc") => query.OrderByDescending(p => p.AtualizadoEm),

            ("descricao", _) => query.OrderBy(p => p.Descricao),
            ("saldo", _) => query.OrderBy(p => p.Saldo),
            ("criadoem", _) => query.OrderBy(p => p.CriadoEm),
            ("atualizadoem", _) => query.OrderBy(p => p.AtualizadoEm),
            _ => query.OrderBy(p => p.Codigo)
        };

        var totalItems = await query.CountAsync(ct);

        var items = await query
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(ct);

        return (items, totalItems);
    }

    public async Task AdicionarAsync(Produto produto, CancellationToken ct = default) =>
        await context.Produtos.AddAsync(produto, ct);

    public async Task<bool> ExisteCodigoAsync(string codigo, Guid? ignorarId = null, CancellationToken ct = default) =>
        await context.Produtos.AnyAsync(p =>
            p.Codigo == codigo.ToUpper() && (ignorarId == null || p.Id != ignorarId), ct);

    public async Task AtualizarAsync(Produto produto, CancellationToken ct = default)
    {
        // foreach (var mov in produto.Movimentacoes)
        //context.Movimentacoes.Add(mov);
        //return Task.CompletedTask;

        var trackedProduto = await context.Produtos
            .Include(p => p.Movimentacoes)
            .FirstOrDefaultAsync(p => p.Id == produto.Id, ct);

        if (trackedProduto is null)
            throw new InvalidOperationException($"Produto '{produto.Id}' não encontrado para atualização.");

        context.Entry(trackedProduto).CurrentValues.SetValues(produto);

        var idsExistentes = trackedProduto.Movimentacoes
            .Select(m => m.Id)
            .ToHashSet();

        var movimentacoesNovas = produto.Movimentacoes
            .Where(m => !idsExistentes.Contains(m.Id))
            .ToList();

        foreach (var mov in movimentacoesNovas)
            context.Movimentacoes.Add(mov);
    }
}

public class MovimentacaoRepository(EstoqueDbContext context) : IMovimentacaoRepository
{
    public async Task<IEnumerable<Movimentacao>> ListarPorProdutoAsync(Guid produtoId, CancellationToken ct = default) =>
        await context.Movimentacoes
            .Where(m => m.ProdutoId == produtoId)
            .OrderByDescending(m => m.OcorridoEm)
            .ToListAsync(ct);

}
public class UnitOfWork(EstoqueDbContext context) : IUnitOfWork
{
    public async Task<int> CommitAsync(CancellationToken ct = default)
    {
        try
        {
            return await context.SaveChangesAsync(ct);
        }
        catch (DbUpdateConcurrencyException ex)
        {
            // Concorrência detectada — o EF lança isso quando RowVersion diverge
            throw new Application.Exceptions.ConcurrencyException(
                "Conflito de concorrência detectado. O produto foi modificado por outra operação simultânea.", ex);
        }
    }
}
