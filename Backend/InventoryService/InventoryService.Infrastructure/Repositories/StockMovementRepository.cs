
using InventoryService.Domain.Entities;
using InventoryService.Domain.Repositories;
using InventoryService.Infrastructure.Data;
using Microsoft.EntityFrameworkCore;

namespace InventoryService.Infrastructure.Repositories
{
    public class StockMovementRepository : IStockMovementRepository
    {
        private readonly InventoryDbContext _context;
        public StockMovementRepository(InventoryDbContext context) => _context = context;

        public async Task AddAsync(StockMovement movement)
        {
            _context.StockMovements.Add(movement);
            await _context.SaveChangesAsync();
        }

        public async Task<bool> OperationExistsAsync(string operationKey)
        {
            return await _context.StockMovements
                .AnyAsync(m => m.OperationKey == operationKey);
        }

        public async Task<IEnumerable<StockMovement>> GetAllAsync() =>
            await _context.StockMovements.ToListAsync();
    }
}