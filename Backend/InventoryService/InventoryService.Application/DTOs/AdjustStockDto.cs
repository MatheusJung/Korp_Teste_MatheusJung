namespace InventoryService.Application.DTOs
{
    public class AdjustStockDto
    {
        public string ProductCode { get; set; } = string.Empty;
        public int Quantity { get; set; }
        public string OperationKey { get; set; } = string.Empty;
    }
}
