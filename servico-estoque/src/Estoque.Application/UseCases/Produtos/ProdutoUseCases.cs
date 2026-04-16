using Estoque.Application.DTOs;
using Estoque.Application.Interfaces;
using Estoque.Domain.Entities;
using Estoque.Domain.Exceptions;

namespace Estoque.Application.UseCases.Produtos;

public class CriarProdutoUseCase(IProdutoRepository repo, IUnitOfWork uow)
{
    public async Task<ProdutoResponse> ExecutarAsync(
        CriarProdutoRequest request, CancellationToken ct = default)
    {
        if (await repo.ExisteCodigoAsync(request.Codigo, ct: ct))
            throw new DomainException($"Já existe um produto com o código '{request.Codigo}'.");

        var produto = Produto.Criar(request.Codigo, request.Descricao, request.SaldoInicial);

        await repo.AdicionarAsync(produto, ct);
        await uow.CommitAsync(ct);

        return produto.ToResponse();
    }
}

public class ObterProdutoUseCase(IProdutoRepository repo)
{
    public async Task<ProdutoResponse> ExecutarAsync(Guid id, CancellationToken ct = default)
    {
        var produto = await repo.ObterPorIdAsync(id, ct)
            ?? throw new ProdutoNaoEncontradoException(id);

        return produto.ToResponse();
    }
}

public class ListarProdutosUseCase(IProdutoRepository repo)
{
    public async Task<IEnumerable<ProdutoResponse>> ExecutarAsync(CancellationToken ct = default)
    {
        var produtos = await repo.ListarAsync(ct);
        return produtos.Select(p => p.ToResponse());
    }
}

public class ListarProdutosPaginadoUseCase(IProdutoRepository repo)
{
    public async Task<PagedResult<ProdutoResponse>> ExecutarAsync(
        int page,
        int pageSize,
        string? search,
        string? sortBy,
        string? sortDirection,
        CancellationToken ct = default)
    {
        if (page < 1)
            page = 1;

        if (pageSize < 1)
            pageSize = 10;

        if (pageSize > 100)
            pageSize = 100;

        search = string.IsNullOrWhiteSpace(search)
            ? null
            : search.Trim();

        sortBy = string.IsNullOrWhiteSpace(sortBy)
            ? "codigo"
            : sortBy.Trim().ToLower();

        sortDirection = string.IsNullOrWhiteSpace(sortDirection)
            ? "asc"
            : sortDirection.Trim().ToLower();

        var (produtos, totalItems) = await repo.ListarPaginadoAsync(
            page,
            pageSize,
            search,
            sortBy,
            sortDirection,
            ct);

        var items = produtos.Select(p => p.ToResponse());

        var totalPages = totalItems == 0
            ? 0
            : (int)Math.Ceiling(totalItems / (double)pageSize);

        return new PagedResult<ProdutoResponse>(
            items,
            page,
            pageSize,
            totalItems,
            totalPages
        );
    }
}

// Extension para mapeamento sem AutoMapper
internal static class ProdutoMappingExtensions
{
    public static ProdutoResponse ToResponse(this Produto p) =>
        new(p.Id, p.Codigo, p.Descricao, p.Saldo, p.CriadoEm, p.AtualizadoEm);
}
