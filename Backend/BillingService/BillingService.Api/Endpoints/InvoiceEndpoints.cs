using BillingService.Application.DTOs;
using BillingService.Application.Interfaces;
using BillingService.Application.Mappers;
using BillingService.Application.Services;
using BillingService.Application.Validators;
using BillingService.Domain.Exceptions;
using BillingService.Infrastructure.Repositories;
using Microsoft.AspNetCore.Http;

namespace BillingService.Api.Endpoints
{
    public static class InvoiceEndpoints
    {
        public static void MapInvoiceEndpoints(this WebApplication app)
        {
            var group = app.MapGroup("/invoices").WithTags("Invoices");

            // Listar notas
            group.MapGet("/all", async (InvoiceService service) =>
            {
                try
                {
                    var invoices = (await service.ListAsync())
                        .OrderByDescending(i => i.CreatedAt)
                        .ToInvoiceDtoList();

                    return Results.Ok(invoices);
                }
                catch (Exception ex)
                {
                    return Results.Problem("Erro interno ao listar as notas fiscais.", ex.Message);
                }
            })
            .WithName("GetAllInvoices")
            .WithDisplayName("Get Invoices")
            .WithSummary("Lista todas as notas fiscais")
            .WithDisplayName("Listar Notas Fiscais")
            .Produces<List<InvoiceDto>>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status500InternalServerError);

            // Listar notas paginadas
            group.MapGet("/", async (InvoiceService service, int pageNumber = 1, int pageSize = 10) =>
            {
                try
                {
                    if (pageNumber <= 0 || pageSize <= 0)
                        return Results.BadRequest("pageNumber e pageSize devem ser maiores que zero.");

                    var allInvoices = await service.ListAsync();
                    var totalItems = allInvoices.Count();

                    var pagedInvoices = allInvoices
                        .OrderByDescending(i => i.CreatedAt)
                        .Skip((pageNumber - 1) * pageSize)
                        .Take(pageSize)
                        .ToInvoiceDtoList();

                    var result = new
                    {
                        PageNumber = pageNumber,
                        PageSize = pageSize,
                        TotalItems = totalItems,
                        TotalPages = (int)Math.Ceiling(totalItems / (double)pageSize),
                        Items = pagedInvoices
                    };

                    return Results.Ok(result);
                }
                catch (Exception ex)
                {
                    return Results.Problem("Erro interno ao listar as notas fiscais.", ex.Message);
                }
            })
            .WithName("GetInvoicesPaged")
            .WithDisplayName("Get Invoices Paginated")
            .WithSummary("Lista todas as notas fiscais com paginação")
            .Produces(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status500InternalServerError);

            // Consultar por sequencial
            group.MapGet("/{seqNumber:int}", async (int seqNumber, InvoiceService service) =>
            {
                try
                {
                    var invoice = await service.GetBySequentialNumberAsync(seqNumber);
                    InvoiceValidator.EnsureInvoiceExists(invoice, seqNumber);

                    return Results.Ok(invoice!.ToInvoiceDto());
                }
                catch (InvoiceNotFoundException ex)
                {
                    return Results.NotFound(ex.Message);
                }
                catch (Exception ex)
                {
                    return Results.Problem("Erro interno ao buscar a nota fiscal.", ex.Message);
                }
            })
            .WithName("GetInvoiceBySequentialNumber")
            .WithSummary("Consulta uma nota fiscal pelo número sequencial")
            .WithDisplayName("Consultar Nota Fiscal por Sequencial")
            .WithDescription("Consulta uma nota fiscal pelo seu número sequencial.")
            .Produces<InvoiceDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound)
            .Produces(StatusCodes.Status500InternalServerError);

            // Criar nota
            group.MapPost("/", async (CreateInvoiceDto createDto, InvoiceService service) =>
            {
                try
                {
                    var invoice = await service.CreateAsync(createDto.Items);
                    return Results.Created($"/invoices/{invoice.SequentialNumber}", invoice.ToInvoiceDto());
                }
                catch (ArgumentException ex)
                {
                    return Results.BadRequest(ex.Message);
                }
                catch (Exception ex)
                {
                    return Results.Problem("Erro interno ao criar a nota fiscal.", ex.Message);
                }
            })
            .WithName("CreateInvoice")
            .WithSummary("Cria uma nova nota fiscal a partir de uma lista de itens")
            .WithDisplayName("Criar Nota Fiscal")
            .WithDescription("Cria uma nova nota fiscal com os itens especificados.")
            .Produces<InvoiceDto>(StatusCodes.Status201Created)
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status500InternalServerError);

            // Adicionar item a nota aberta
            group.MapPut("/{seqNumber:int}/items", async (int seqNumber, InvoiceItemDto dto, InvoiceService service) =>
            {
                try
                {
                    InvoiceValidator.ValidateQuantity(dto.Quantity);

                    await service.AddItemAsync(seqNumber, dto.ProductCode, dto.Quantity);

                    var updatedInvoice = await service.GetBySequentialNumberAsync(seqNumber);
                    return Results.Ok(updatedInvoice!.ToInvoiceDto());
                }
                catch (InvoiceNotFoundException ex)
                {
                    return Results.NotFound(ex.Message);
                }
                catch (ArgumentException ex)
                {
                    return Results.BadRequest(ex.Message);
                }
                catch (Exception ex)
                {
                    return Results.Problem("Erro ao adicionar item à nota fiscal.", ex.Message);
                }
            })
            .WithName("AddItemToInvoice")
            .WithSummary("Adiciona um item a uma nota fiscal aberta")
            .WithDisplayName("Adicionar Item à Nota Fiscal")
            .WithDescription("Adiciona um item a uma nota fiscal que esteja em aberto.")
            .Produces<InvoiceDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status404NotFound)
            .Produces(StatusCodes.Status500InternalServerError);

            // Cancelar nota
            group.MapPost("/{seqNumber:int}/cancel", async (int seqNumber, InvoiceService service) =>
            {
                try
                {
                    await service.CancelAsync(seqNumber);

                    var updatedInvoice = await service.GetBySequentialNumberAsync(seqNumber);
                    return Results.Ok(updatedInvoice!.ToInvoiceDto());
                }
                catch (InvoiceNotFoundException ex)
                {
                    return Results.NotFound(ex.Message);
                }
                catch (Exception ex)
                {
                    return Results.Problem("Erro ao cancelar a nota fiscal.", ex.Message);
                }
            })
            .WithName("CancelInvoice")
            .WithSummary("Cancela uma nota fiscal existente")
            .WithDisplayName("Cancelar Nota Fiscal")
            .WithDescription("Cancela uma nota fiscal existente.")
            .Produces<InvoiceDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound)
            .Produces(StatusCodes.Status500InternalServerError);

            //fechar nota
            group.MapPost("/{seqNumber:int}/close", async (int seqNumber, InvoiceService service) =>
            {
                try
                {
                    await service.CloseAsync(seqNumber);

                    var closedInvoice = await service.GetBySequentialNumberAsync(seqNumber);
                    return Results.Ok(closedInvoice!.ToInvoiceDto());
                }
                catch (InvoiceNotFoundException ex)
                {
                    return Results.NotFound(ex.Message);
                }
                catch (InvalidOperationException ex)
                {
                    return Results.BadRequest(ex.Message);
                }
                catch (Exception ex)
                {
                    return Results.Problem("Erro ao fechar a nota fiscal.", ex.Message);
                }
            })
            .WithName("CloseInvoice")
            .WithDescription("fecha uma nota fiscal existente.")
            .WithSummary("fecha uma nota fiscal aberta")
            .WithDisplayName("Fechar Nota Fiscal")
            .Produces<InvoiceDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status404NotFound)
            .Produces(StatusCodes.Status500InternalServerError);

            // Imprimir e fechar nota
            group.MapPost("/{seqNumber:int}/close-and-print", async (int seqNumber, InvoiceService service, IInvoicePrinter printer) =>
            {
                try
                {
                    // Busca a nota fiscal pelo numero sequencial
                    var invoice = await service.GetBySequentialNumberAsync(seqNumber);
                    if (invoice == null)
                        return Results.NotFound("NF não encontrada");

                    // Fecha a nota fiscal
                    await service.CloseAsync(seqNumber);

                    // Recarrega a nota fiscal atualizada
                    var updatedInvoice = await service.GetBySequentialNumberAsync(seqNumber);

                    // Conveter para DTO
                    var invoiceDto = updatedInvoice?.ToInvoiceDto();
                    if (invoiceDto == null)
                        return Results.NoContent();

                    // Gerar PDF da nota fiscal
                    var pdfBytes = printer.Print(invoiceDto);

                    // Retorna o PDF como resposta
                    return Results.File(pdfBytes,"application/pdf", $"NF-{invoice.SequentialNumber}.pdf");

                }
                catch (InvoiceNotFoundException ex)
                {
                    return Results.NotFound(ex.Message);
                }
                catch (InvalidOperationException ex)
                {
                    return Results.BadRequest(ex.Message);
                }
                catch (Exception ex)
                {
                    return Results.Problem("Erro ao fechar a nota fiscal.", ex.Message);
                }
            })
            .WithName("CloseAndPrintInvoice")
            .WithDescription("Imprime e fecha uma nota fiscal existente.")
            .WithSummary("Imprime e fecha uma nota fiscal aberta")
            .WithDisplayName("Fechar Nota Fiscal")
            .Produces(StatusCodes.Status200OK, contentType: "application/pdf")
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status404NotFound)
            .Produces(StatusCodes.Status500InternalServerError);
        }
    }
}