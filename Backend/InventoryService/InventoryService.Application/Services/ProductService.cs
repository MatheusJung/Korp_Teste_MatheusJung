using InventoryService.Application.Exceptions;
using InventoryService.Domain.Entities;
using InventoryService.Domain.Repositories;
using InventoryService.Infrastructure.Repositories;

namespace InventoryService.Application.Services;

public class ProductService
{
    private readonly IProductRepository _productRepo;

    public ProductService(IProductRepository productRepo)
    {
        _productRepo = productRepo;
    }

    /// <summary>
    /// Creates a new product.
    /// </summary>
    public async Task<Product> CreateAsync(string code, string name, int initialQuantity = 0)
    {
        var existing = await _productRepo.GetByCodeAsync(code);
        if (existing is not null)
            throw new DuplicateProductCodeException(code);

        var product = new Product(code, name, initialQuantity);
        await _productRepo.AddAsync(product);
        return product;
    }

    /// <summary>
    /// Retrieves all Products
    /// </summary>
    public async Task<IEnumerable<Product>> GetAllAsync() =>
       (await _productRepo.GetAllAsync())
           .OrderByDescending(p => p.CreatedAt);

    /// <summary>
    /// Retrieves all Active Products
    /// </summary>
    public async Task<IEnumerable<Product>> GetAllActiveAsync() =>
       await _productRepo.GetAllActiveAsync();

    /// <summary>
    /// Retrieves all Inactive Products
    /// </summary>
    public async Task<IEnumerable<Product>> GetAllInactiveAsync() =>
       await _productRepo.GetAllInactiveAsync();

    /// <summary>
    /// Retrieves a product by its Id.
    /// </summary>
    public async Task<Product?> GetByIdAsync(int id) =>
        await _productRepo.GetByIdAsync(id);

    /// <summary>
    /// Retrieves a product by its Code.
    /// </summary>
    public async Task<Product?> GetByCodeAsync(string code) =>
        await _productRepo.GetByCodeAsync(code);

    /// <summary>
    /// Deactivates a product.
    /// </summary>
    public async Task DeactivateAsync(int id)
    {
        var product = await _productRepo.GetByIdAsync(id)
            ?? throw new KeyNotFoundException($"Product {id} not found.");

        if (!product.IsActive)
            throw new InvalidOperationException("Product already deactivated.");

        await _productRepo.DeactivateAsync(product);
    }

    /// <summary>
    /// Reactivates a product.
    /// </summary>
    public async Task ActivateAsync(int id)
    {
        var product = await _productRepo.GetByIdAsync(id)
              ?? throw new KeyNotFoundException($"Product {id} not found.");

        if (product.IsActive)
            throw new InvalidOperationException("Product already active.");

        await _productRepo.ActivateAsync(product);
    }
}
