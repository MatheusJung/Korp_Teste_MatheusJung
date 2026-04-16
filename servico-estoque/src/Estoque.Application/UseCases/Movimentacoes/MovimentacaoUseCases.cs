using Estoque.Application.DTOs;
using Estoque.Application.Interfaces;
using Estoque.Domain.Exceptions;

namespace Estoque.Application.UseCases.Movimentacoes;

/// <summary>
/// Deduz estoque de múltiplos produtos atomicamente.
/// Se qualquer item falhar por saldo insuficiente, estorna os já deduzidos e retorna erro.
/// Chamado pelo Serviço de Faturamento via outbox worker.
/// </summary>
public class DeduzirLoteUseCase(IProdutoRepository repo, IUnitOfWork uow)
{
    public async Task<DeduzirLoteResponse> ExecutarAsync(DeduzirLoteRequest request, CancellationToken ct = default)
    {
        var itensDeduzidos = new List<ItemDeduzidoResponse>();
        var produtosModificados = new List<(Estoque.Domain.Entities.Produto produto, decimal quantidade)>();

        foreach (var item in request.Itens)
        {
            var produto = await repo.ObterSemMovimentacoesAsync(item.ProdutoId, ct);

            if (produto is null)
            {
                await EstornarParcialAsync(produtosModificados, request.NotaFiscalId, ct);

                return new DeduzirLoteResponse(
                    Sucesso: false,
                    ItensDeduzidos: itensDeduzidos,
                    Erro: $"Produto '{item.ProdutoId}' não encontrado.",
                    ProdutoFalhou: item.ProdutoId);
            }

            try
            {
                produto.Deduzir(item.Quantidade, request.NotaFiscalId);
                await repo.AtualizarAsync(produto, ct);

                produtosModificados.Add((produto, item.Quantidade));
                itensDeduzidos.Add(new ItemDeduzidoResponse(item.ProdutoId, item.Quantidade));
            }
            catch (SaldoInsuficienteException ex)
            {
                await EstornarParcialAsync(produtosModificados, request.NotaFiscalId, ct);

                return new DeduzirLoteResponse(
                    Sucesso: false,
                    ItensDeduzidos: itensDeduzidos,
                    Erro: ex.Message,
                    ProdutoFalhou: ex.ProdutoId);
            }
        }

        await uow.CommitAsync(ct);

        return new DeduzirLoteResponse(
            Sucesso: true,
            ItensDeduzidos: itensDeduzidos);
    }

    private async Task EstornarParcialAsync(
        List<(Estoque.Domain.Entities.Produto produto, decimal quantidade)> itens,
        Guid notaFiscalId,
        CancellationToken ct)
    {
        foreach (var (produto, quantidade) in itens)
        {
            produto.Estornar(quantidade, notaFiscalId);
            await repo.AtualizarAsync(produto, ct);
        }

        if (itens.Count > 0)
            await uow.CommitAsync(ct);
    }
}

/// <summary>
/// Estorna estoque de múltiplos produtos (usado quando nota é cancelada).
/// </summary>
public class EstornarLoteUseCase(IProdutoRepository repo, IUnitOfWork uow)
{
    public async Task ExecutarAsync(EstornarLoteRequest request, CancellationToken ct = default)
    {
        foreach (var item in request.Itens)
        {
            var produto = await repo.ObterSemMovimentacoesAsync(item.ProdutoId, ct)
                ?? throw new ProdutoNaoEncontradoException(item.ProdutoId);

            produto.Estornar(item.Quantidade, request.NotaFiscalId);
            await repo.AtualizarAsync(produto, ct);
        }

        await uow.CommitAsync(ct);
    }
}

/// <summary>
/// Adiciona entrada manual de estoque (reposição).
/// </summary>
public class AdicionarEntradaUseCase(IProdutoRepository repo, IUnitOfWork uow)
{
    public async Task<ProdutoResponse> ExecutarAsync(AdicionarEntradaRequest request, CancellationToken ct = default)
    {
        var produto = await repo.ObterSemMovimentacoesAsync(request.ProdutoId, ct)
            ?? throw new ProdutoNaoEncontradoException(request.ProdutoId);

        produto.AdicionarEntrada(request.Quantidade);

        await repo.AtualizarAsync(produto, ct);
        await uow.CommitAsync(ct);

        return new ProdutoResponse(
            produto.Id,
            produto.Codigo,
            produto.Descricao,
            produto.Saldo,
            produto.CriadoEm,
            produto.AtualizadoEm
        );
    }
}