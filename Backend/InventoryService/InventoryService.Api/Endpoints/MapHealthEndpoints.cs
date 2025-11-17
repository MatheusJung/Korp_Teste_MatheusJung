namespace InventoryService.Api.Endpoints
{
    public static class HealthEndpoints
    {
        public static void MapHealthEndpoints(this WebApplication app)
        {
            var group = app.MapGroup("/health").WithTags("Health");

            // Health endpoint simples
            group.MapGet("/", () => new { status = "UP", service = "InventoryService" });
        }
    }
}
