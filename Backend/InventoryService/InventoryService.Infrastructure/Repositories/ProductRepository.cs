using InventoryService.Domain.Entities;
using InventoryService.Domain.Repositories;
using InventoryService.Infrastructure.Data;
using Microsoft.EntityFrameworkCore;

namespace InventoryService.Infrastructure.Repositories;

public class ProductRepository : IProductRepository
{
    private readonly InventoryDbContext _context;
    public ProductRepository(InventoryDbContext context) => _context = context;

    public async Task<IEnumerable<Product>> GetAllAsync() =>
        await _context.Products.ToListAsync();

    public async Task<IEnumerable<Product>> GetAllActiveAsync() =>
        await _context.Products.Where(p => p.IsActive).ToListAsync();

    public async Task<IEnumerable<Product>> GetAllInactiveAsync() =>
        await _context.Products.Where(p => !p.IsActive).ToListAsync();

    public async Task<Product?> GetByCodeAsync(string code)
    {
        return await _context.Products.FirstOrDefaultAsync(p => p.Code == code);

    }

    public async Task<Product?> GetByIdAsync(int id) =>
        await _context.Products.FindAsync(id);

    public async Task AddAsync(Product product)
    {
        _context.Products.Add(product);
        await _context.SaveChangesAsync();
    }

    public async Task UpdateAsync(Product product)
    {
        _context.Products.Update(product);
        await _context.SaveChangesAsync();
    }
    public async Task DeactivateAsync(Product product)
    {
        product.Deactivate();
        _context.Products.Update(product);
        await _context.SaveChangesAsync();
    }
    public async Task ActivateAsync(Product product)
    {
        product.Activate();
        _context.Products.Update(product);
        await _context.SaveChangesAsync();
    }
}
