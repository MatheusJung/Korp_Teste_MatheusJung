using InventoryService.Application.DTOs;
using InventoryService.Domain.Entities;
using InventoryService.Domain.Enums;
using InventoryService.Domain.Exceptions;
using InventoryService.Domain.Repositories;

public class StockService
{
    private readonly IProductRepository _productRepo;
    private readonly IStockMovementRepository _movementRepo;

    public StockService(IProductRepository productRepo, IStockMovementRepository movementRepo)
    {
        _productRepo = productRepo;
        _movementRepo = movementRepo;
    }

    public async Task<IEnumerable<StockMovement>> GetAllAsync() =>
        await _movementRepo.GetAllAsync();

    public Task AddStockAsync(AdjustStockDto dto)
     => AdjustStockByCodeAsync(dto, MovementType.Entrada);

    public Task RemoveStockAsync(AdjustStockDto dto)
        => AdjustStockByCodeAsync(dto, MovementType.Saida);

    public Task AdjustStockAsync(AdjustStockDto dto)
        => AdjustStockByCodeAsync(dto, MovementType.Ajuste);

    public async Task AdjustStockByCodeAsync(AdjustStockDto dto, MovementType movementType)
    {
        // Verifica idempotência
        if (!string.IsNullOrEmpty(dto.OperationKey) && await _movementRepo.OperationExistsAsync(dto.OperationKey))
            throw new OperationAlreadyAppliedException(dto.OperationKey);

        // Busca produto pelo código
        var product = await _productRepo.GetByCodeAsync(dto.ProductCode)
                      ?? throw new KeyNotFoundException($"Produto '{dto.ProductCode}' não encontrado");

        // Verifica se o produto está ativo
        if (!product.IsActive)
            throw new InvalidOperationException($"Cannot adjust stock. Product '{dto.ProductCode}' is inactive.");

        if (dto.Quantity <= 0)
            throw new ArgumentOutOfRangeException(nameof(dto.Quantity), "Quantity must be non-zero.");
        // Ajusta estoque
        if (movementType == MovementType.Entrada)
            product.IncreaseStock(dto.Quantity);
        else if (movementType == MovementType.Saida)
            product.DecreaseStock(dto.Quantity);

        await _productRepo.UpdateAsync(product);

        // Cria registro de movimento
        var movement = new StockMovement(
            product.Id,
            dto.Quantity,
            movementType,
            dto.OperationKey
        );

        await _movementRepo.AddAsync(movement);
    }
}