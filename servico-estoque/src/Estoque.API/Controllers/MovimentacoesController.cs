using Estoque.Application.DTOs;
using Estoque.Application.Interfaces;
using Estoque.Application.UseCases.Movimentacoes;
using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;

namespace Estoque.API.Controllers;

[ApiController]
[Route("api/movimentacoes")]
[Produces("application/json")]
public class MovimentacoesController : ControllerBase
{
    private readonly DeduzirLoteUseCase _deduzirLoteUseCase;
    private readonly EstornarLoteUseCase _estornarLoteUseCase;
    private readonly AdicionarEntradaUseCase _adicionarEntradaUseCase;
    private readonly IMovimentacaoRepository _movimentacaoRepo;

    public MovimentacoesController(
        DeduzirLoteUseCase deduzirLoteUseCase,
        EstornarLoteUseCase estornarLoteUseCase,
        AdicionarEntradaUseCase adicionarEntradaUseCase,
        IMovimentacaoRepository movimentacaoRepo)
    {
        _deduzirLoteUseCase = deduzirLoteUseCase;
        _estornarLoteUseCase = estornarLoteUseCase;
        _adicionarEntradaUseCase = adicionarEntradaUseCase;
        _movimentacaoRepo = movimentacaoRepo;
    }

    /// <summary>
    /// Deduz estoque de múltiplos produtos
    /// </summary>
    /// <remarks>
    /// Sempre retorna HTTP 200.
    /// </remarks>
    [HttpPost("deduzir-lote")]
    [SwaggerOperation(Summary = "Deduz estoque em lote")]
    [ProducesResponseType(typeof(DeduzirLoteResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<IActionResult> DeduzirLote(
        [FromBody] DeduzirLoteRequest request,
        CancellationToken ct)
    {
        var response = await _deduzirLoteUseCase.ExecutarAsync(request, ct);
        return Ok(response);
    }

    /// <summary>
    /// Estorna estoque de produtos
    /// </summary>
    /// <remarks>
    /// Usado quando uma nota fiscal é cancelada.
    /// </remarks>
    [HttpPost("estornar-lote")]
    [SwaggerOperation(Summary = "Estornar estoque")]
    [ProducesResponseType(StatusCodes.Status204NoContent)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<IActionResult> EstornarLote(
        [FromBody] EstornarLoteRequest request,
        CancellationToken ct)
    {
        await _estornarLoteUseCase.ExecutarAsync(request, ct);
        return NoContent();
    }

    /// <summary>
    /// Adiciona entrada de estoque
    /// </summary>
    /// <remarks>
    /// Entrada manual de estoque.
    /// </remarks>
    [HttpPost("entrada")]
    [SwaggerOperation(Summary = "Adicionar entrada de estoque")]
    [ProducesResponseType(typeof(ProdutoResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<IActionResult> AdicionarEntrada(
        [FromBody] AdicionarEntradaRequest request,
        CancellationToken ct)
    {
        var produto = await _adicionarEntradaUseCase.ExecutarAsync(request, ct);
        return Ok(produto);
    }

    /// <summary>
    /// Lista movimentações de um produto
    /// </summary>
    /// <remarks>
    /// Retorna o histórico completo de movimentações.
    /// </remarks>
    [HttpGet("produto/{produtoId:guid}")]
    [SwaggerOperation(Summary = "Listar movimentações por produto")]
    [ProducesResponseType(typeof(IEnumerable<MovimentacaoResponse>), StatusCodes.Status200OK)]
    public async Task<IActionResult> ListarPorProduto(
        Guid produtoId,
        CancellationToken ct)
    {
        var movimentacoes = await _movimentacaoRepo.ListarPorProdutoAsync(produtoId, ct);

        var response = movimentacoes.Select(m => new MovimentacaoResponse(
            m.Id,
            m.ProdutoId,
            m.Tipo,
            m.Quantidade,
            m.SaldoAnterior,
            m.SaldoResultante,
            m.NotaFiscalId,
            m.IsEstorno,
            m.OcorridoEm
        ));

        return Ok(response);
    }
}