namespace InventoryService.Application.DTOs
{
    public class CreateProductDto
    {
        public string Code { get; set; } = null!;
        public string Name { get; set; } = null!;
        public int InitialStock { get; set; }
    }
}
