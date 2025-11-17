using BillingService.Application.DTOs;

namespace BillingService.Application.ContractsHttp
{
    public interface IHttpInventoryClient
    {
        Task<bool> IsInventoryServiceAvailableAsync();
        Task<ItemDto?> GetProductByCodeAsync(string productCode);
        Task AdjustStockAsync(string productCode, int quantity, string operationKey);
    }
}
