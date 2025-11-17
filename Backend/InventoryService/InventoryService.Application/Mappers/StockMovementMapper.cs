using InventoryService.Application.DTOs;
using InventoryService.Domain.Entities;

namespace InventoryService.Application.Mappers
{
    public static class StockMovementMapper
    {
        public static StockMovementDto ToDto(StockMovement movement) =>
            new StockMovementDto(
                movement.ProductId,
                movement.QuantityChange,
                movement.MovementType.ToString(),
                movement.OperationKey,
                movement.CreatedAt
            );

        public static List<StockMovementDto> ToDtoList(IEnumerable<StockMovement> movements) =>
            movements.Select(ToDto).ToList();
    }
}