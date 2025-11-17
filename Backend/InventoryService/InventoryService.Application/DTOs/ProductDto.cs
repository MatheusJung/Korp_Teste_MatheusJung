namespace InventoryService.Application.DTOs
{
    public sealed record ProductDto(
        string ProductCode,
        string Description,
        int Quantity,
        string Status
    );
}