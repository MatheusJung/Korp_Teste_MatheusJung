using Estoque.Application.Exceptions;
using Estoque.Domain.Exceptions;
using System.Net;
using System.Text.Json;

namespace Estoque.API.Middleware;

public class ExceptionMiddleware(RequestDelegate next, ILogger<ExceptionMiddleware> logger)
{
    private static readonly JsonSerializerOptions JsonOptions = new() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase };

    public async Task InvokeAsync(HttpContext context)
    {
        try
        {
            await next(context);
        }
        catch (Exception ex)
        {
            logger.LogError(ex, "Erro não tratado: {Message}", ex.Message);
            await HandleExceptionAsync(context, ex);
        }
    }

    private static async Task HandleExceptionAsync(HttpContext context, Exception ex)
    {
        var (status, titulo, detalhe) = ex switch
        {
            SaldoInsuficienteException e => (
                HttpStatusCode.UnprocessableEntity,
                "Saldo insuficiente",
                e.Message),

            ProdutoNaoEncontradoException e => (
                HttpStatusCode.NotFound,
                "Produto não encontrado",
                e.Message),

            DomainException e => (
                HttpStatusCode.BadRequest,
                "Regra de negócio violada",
                e.Message),

            ConcurrencyException e => (
                HttpStatusCode.Conflict,
                "Conflito de concorrência",
                e.Message),

            OperationCanceledException => (
                HttpStatusCode.ServiceUnavailable,
                "Operação cancelada",
                "A requisição foi cancelada."),

            _ => (
                HttpStatusCode.InternalServerError,
                "Erro interno",
                "Ocorreu um erro inesperado. Tente novamente.")
        };

        context.Response.ContentType = "application/json";
        context.Response.StatusCode = (int)status;

        var body = JsonSerializer.Serialize(new
        {
            status = (int)status,
            titulo,
            detalhe,
            timestamp = DateTime.UtcNow
        }, JsonOptions);

        await context.Response.WriteAsync(body);
    }
}
