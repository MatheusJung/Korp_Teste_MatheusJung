using InventoryService.Domain.Entities;

namespace InventoryService.Domain.Repositories
{
    public interface IStockMovementRepository
    {
        Task<IEnumerable<StockMovement>> GetAllAsync();
        Task AddAsync(StockMovement movement);
        Task<bool> OperationExistsAsync(string operationKey);
    }
}
