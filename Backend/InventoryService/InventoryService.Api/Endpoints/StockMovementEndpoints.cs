using InventoryService.Application.DTOs;
using InventoryService.Application.Mappers;
using InventoryService.Application.Services;
using InventoryService.Domain.Enums;
using InventoryService.Domain.Exceptions;
using Microsoft.AspNetCore.Mvc;

public static class StockEndpoints
{
    public static void MapStockEndpoints(this WebApplication app)
    {
        var group = app.MapGroup("/stock").WithTags("Stock");

        //Listar todos produtos
        group.MapGet("/", async ([FromServices] StockService service) =>
        {
            var movements = (await service.GetAllAsync());
            var dtos = StockMovementMapper.ToDtoList(movements);
            return Results.Ok(dtos);
        })
        .WithName("GetAllStockMovements")
        .WithDisplayName("Listar movimentações")
        .WithSummary("Lista de movimentações do estoque")
        .WithDescription("Retorna todos os movimentos do estoque");


        group.MapPost("/add", async ([FromBody] AdjustStockDto dto, [FromServices] StockService service) =>
        {
            try
            {
                await service.AddStockAsync(dto);
                return Results.Ok(new { message = "Stock added successfully" });
            }
            catch (Exception ex)
            {
                return Results.BadRequest(new { error = ex.Message });
            }
        })
        .WithName("AddStock")
        .WithDisplayName("Adicionar estoque de produto")
        .WithSummary("Adiciona uma quantidade ao estoque de um produto")
        .WithDescription("Movimenta quantia de produto para dentro do estoque")
        .Produces(StatusCodes.Status200OK)
        .Produces(StatusCodes.Status400BadRequest);

        group.MapPost("/remove", async ([FromBody] AdjustStockDto dto, [FromServices] StockService service) =>
        {
            try
            {
                await service.RemoveStockAsync(dto);
                return Results.Ok(new { message = "Stock added successfully" });
            }
            catch (Exception ex)
            {
                return Results.BadRequest(new { error = ex.Message });
            }

        })
        .WithName("RemoveStock")
        .WithDisplayName("Remove estoque de produto")
        .WithSummary("Remove uma quantidade do estoque de um produto")
        .WithDescription("Movimenta quantia de produto para fora do estoque")
        .Produces(StatusCodes.Status200OK)
        .Produces(StatusCodes.Status400BadRequest);
    }
}
