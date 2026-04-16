using Microsoft.OpenApi.Models;
using Swashbuckle.AspNetCore.SwaggerGen;

namespace Estoque.API.Infrastructure;

/// <summary>
/// Exclui endpoints de infraestrutura (health checks) da documentação Swagger.
/// </summary>
public class HealthCheckSchemaFilter : IDocumentFilter
{
    public void Apply(OpenApiDocument swaggerDoc, DocumentFilterContext context)
    {
        var healthPaths = swaggerDoc.Paths
            .Where(p => p.Key.StartsWith("/health"))
            .Select(p => p.Key)
            .ToList();

        foreach (var path in healthPaths)
            swaggerDoc.Paths.Remove(path);
    }
}
