using InventoryService.Application.DTOs;
using InventoryService.Domain.Entities;

public static class ProductMapper
{
    public static ProductDto ToDto(Product product)
    {
        return new ProductDto(
            product.Code,
            product.Name,
            product.Quantity,
            product.IsActive.ToString()
        );
    }

    public static List<ProductDto> ToDtoList(IEnumerable<Product> products)
    {
        return products.Select(ToDto).ToList();
    }
}