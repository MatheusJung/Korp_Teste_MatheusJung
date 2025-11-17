using InventoryService.Domain.Entities;

namespace InventoryService.Domain.Repositories
{
    public interface IProductRepository
    {
        Task<IEnumerable<Product>> GetAllAsync();
        Task<IEnumerable<Product>> GetAllInactiveAsync();
        Task<IEnumerable<Product>> GetAllActiveAsync();
        Task<Product?> GetByIdAsync(int id);
        Task<Product?> GetByCodeAsync(string code);
        Task AddAsync(Product product);
        Task UpdateAsync(Product product);
        Task DeactivateAsync(Product product);
        Task ActivateAsync(Product product);
    }
}
