namespace InventoryService.Application.DTOs
{
    public sealed record StockMovementDto(
        int ProductId,
        int QuantityChange,
        string MovementType,
        string OperationKey,
        DateTime CreatedAt
    );
}
