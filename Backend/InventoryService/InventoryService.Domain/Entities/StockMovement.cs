using InventoryService.Domain.Enums;

namespace InventoryService.Domain.Entities
{
    public sealed class StockMovement
    {
        public int Id { get; private set; }
        public int ProductId { get; private set; }
        public Product Product { get; private set; } = null!;
        public int QuantityChange { get; private set; }
        public MovementType MovementType { get; private set; }
        public string OperationKey { get; private set; } = null!;
        public DateTime CreatedAt { get; private set; } = DateTime.UtcNow;

        private StockMovement() { } // EF Core

        public StockMovement(int productId, int quantityChange, MovementType movementType, string operationKey)
        {
            if (quantityChange == 0)
                throw new ArgumentOutOfRangeException(nameof(quantityChange), "QuantityChange cannot be zero.");
            if (string.IsNullOrWhiteSpace(operationKey))
                throw new ArgumentException("OperationKey is required.", nameof(operationKey));

            ProductId = productId;
            QuantityChange = quantityChange;
            MovementType = movementType;
            OperationKey = operationKey;
            CreatedAt = DateTime.UtcNow;
        }
    }
}
