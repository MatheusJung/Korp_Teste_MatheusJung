using InventoryService.Application.DTOs;
using InventoryService.Application.Services;
using InventoryService.Domain.Entities;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore.Storage.ValueConversion;

namespace InventoryService.Api.Endpoints
{
    public static class ProductEndpoints
    {
        public static void MapProductEndpoints(this WebApplication app)
        {
            var group = app.MapGroup("/products").WithTags("Products");

            //Cadastrar Produto
            group.MapPost("/", async (CreateProductDto dto, [FromServices] ProductService service) =>
            {
                var product = await service.CreateAsync(dto.Code, dto.Name, dto.InitialStock);
                return Results.Created($"/products/{product.Id}", product);
            })
            .WithName("CreateProduct")
            .WithDisplayName("Criar novo produto")
            .WithSummary("Cria um novo produto.")
            .WithDescription("Cadastra um produto com código, nome e estoque inicial.");

            //Listar todos produtos
            group.MapGet("/Products", async ([FromServices] ProductService service) =>
            {
                var products = (await service.GetAllAsync());
                var dtos = ProductMapper.ToDtoList(products);
                return Results.Ok(dtos);
            })
            .WithName("GetAllProducts")
            .WithDisplayName("Listar produtos")
            .WithSummary("Lista todos os produtos.")
            .WithDescription("Retorna todos os produtos cadastrados, ativos ou não.");

            group.MapGet("/", async (ProductService service, int pageNumber = 1, int pageSize = 10) =>
            {
                try
                {
                    if (pageNumber <= 0 || pageSize <= 0)
                        return Results.BadRequest("pageNumber e pageSize devem ser maiores que zero.");

                    var allProducts = await service.GetAllAsync(); // retorna IEnumerable<Product>
                    var totalItems = allProducts.Count();

                    var pagedProducts = allProducts
                        .OrderBy(p => p.Code) // ou .OrderByDescending(p => p.CreatedAt) se tiver
                        .Skip((pageNumber - 1) * pageSize)
                        .Take(pageSize)
                        .ToList();

                    var pagedDtos = ProductMapper.ToDtoList(pagedProducts);

                    var result = new
                    {
                        PageNumber = pageNumber,
                        PageSize = pageSize,
                        TotalItems = totalItems,
                        TotalPages = (int)Math.Ceiling(totalItems / (double)pageSize),
                        Items = pagedDtos
                    };

                    return Results.Ok(result);
                }
                catch (Exception ex)
                {
                    return Results.Problem("Erro interno ao listar os produtos.", ex.Message);
                }
            })
            .WithName("GetProductsPaged")
            .WithDisplayName("Get Products Paginated")
            .WithSummary("Lista todos os produtos com paginação")
            .WithDescription("Retorna produtos paginados com base nos parâmetros fornecidos.")
            .Produces(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status500InternalServerError);
            

            //Listar todos produtos ativos
            group.MapGet("/ActiveProducts", async ([FromServices] ProductService service) =>
            {
                try
                {
                    var products = await service.GetAllActiveAsync();
                    var dtos = ProductMapper.ToDtoList(products);
                    return Results.Ok(dtos);
                }
                catch (Exception ex)
                {
                    return Results.BadRequest(new { error = ex.Message });
                }
            })
            .WithName("GetActiveProducts")
            .WithDisplayName("Listar produtos ativos")
            .WithSummary("Lista produtos ativos.")
            .WithDescription("Retorna apenas produtos marcados como ativos.");

            //Consulta Produto por Código do produto
            group.MapGet("/{productCode}", async (string productCode, [FromServices] ProductService service) =>
            {
                var product = await service.GetByCodeAsync(productCode);
                if(product is null)
                    return Results.NotFound();
                var dto = ProductMapper.ToDto(product);
                return Results.Ok(dto);
            })
            .WithName("GetProductByCode")
            .WithDisplayName("Consultar produto por código")
            .WithSummary("Consulta produto por código.")
            .WithDescription("Busca um produto usando o código único.");

            //Deleta Produto
            group.MapDelete("/{productCode}", async (string productCode, [FromServices] ProductService service) =>
            {
                var product = await service.GetByCodeAsync(productCode);

                if (product is null)
                    return Results.NotFound();

                await service.DeactivateAsync(product.Id);
                return Results.NoContent();
            })
            .WithName("DeactivateProduct")
            .WithDisplayName("Desativar produto")
            .WithSummary("Desativa um produto.")
            .WithDescription("Marca o produto como inativo, mas mantém no banco.");
        }
    }
}
