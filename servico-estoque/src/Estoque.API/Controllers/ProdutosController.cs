using Estoque.Application.DTOs;
using Estoque.Application.UseCases.Produtos;
using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;

namespace Estoque.API.Controllers;

[ApiController]
[Route("api/produtos")]
[Produces("application/json")]
public class ProdutosController : ControllerBase
{
    private readonly CriarProdutoUseCase _criarUseCase;
    private readonly ObterProdutoUseCase _obterUseCase;
    private readonly ListarProdutosUseCase _listarUseCase;
    private readonly ListarProdutosPaginadoUseCase _listarPaginadoUseCase;

    public ProdutosController(
        CriarProdutoUseCase criarUseCase,
        ObterProdutoUseCase obterUseCase,
        ListarProdutosUseCase listarUseCase,
        ListarProdutosPaginadoUseCase listarPaginadoUseCase)
    {
        _criarUseCase = criarUseCase;
        _obterUseCase = obterUseCase;
        _listarUseCase = listarUseCase;
        _listarPaginadoUseCase = listarPaginadoUseCase;
    }

    /// <summary>
    /// Lista todos os produtos
    /// </summary>
    [HttpGet]
    [SwaggerOperation(Summary = "Listar produtos")]
    [ProducesResponseType(typeof(IEnumerable<ProdutoResponse>), StatusCodes.Status200OK)]
    public async Task<IActionResult> Listar(CancellationToken ct)
    {
        var produtos = await _listarUseCase.ExecutarAsync(ct);
        return Ok(produtos);
    }

    /// <summary>
    /// Lista produtos com paginação
    /// </summary>
    [HttpGet("paginado")]
    [SwaggerOperation(Summary = "Listar produtos paginado")]
    [ProducesResponseType(typeof(PagedResult<ProdutoResponse>), StatusCodes.Status200OK)]
    public async Task<IActionResult> ListarPaginado(
    [FromQuery] int page = 1,
    [FromQuery] int pageSize = 10,
    [FromQuery] string? search = null,
    [FromQuery] string? sortBy = "codigo",
    [FromQuery] string? sortDirection = "asc",
    CancellationToken ct = default)
    {
        var resultado = await _listarPaginadoUseCase.ExecutarAsync(
            page,
            pageSize,
            search,
            sortBy,
            sortDirection,
            ct);

        return Ok(resultado);
    }

    /// <summary>
    /// Obtém um produto pelo ID
    /// </summary>
    [HttpGet("{id:guid}")]
    [SwaggerOperation(Summary = "Obter produto por ID")]
    [ProducesResponseType(typeof(ProdutoResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(typeof(ProblemDetails), StatusCodes.Status404NotFound)]
    public async Task<IActionResult> Obter(Guid id, CancellationToken ct)
    {
        var produto = await _obterUseCase.ExecutarAsync(id, ct);

        if (produto is null)
            return NotFound(new ProblemDetails
            {
                Title = "Produto não encontrado",
                Status = StatusCodes.Status404NotFound
            });

        return Ok(produto);
    }

    /// <summary>
    /// Cria um novo produto
    /// </summary>
    [HttpPost]
    [SwaggerOperation(Summary = "Criar produto")]
    [ProducesResponseType(typeof(ProdutoResponse), StatusCodes.Status201Created)]
    [ProducesResponseType(typeof(ValidationProblemDetails), StatusCodes.Status400BadRequest)]
    public async Task<IActionResult> Criar(
        [FromBody] CriarProdutoRequest request,
        CancellationToken ct)
    {
        var resultado = await _criarUseCase.ExecutarAsync(request, ct);

        return CreatedAtAction(
            nameof(Obter),
            new { id = resultado.Id },
            resultado);
    }
}