using BillingService.Domain.Entities;

namespace BillingService.Domain.Interfaces
{
    public interface IInvoiceRepository
    {
        Task<int> GetNextSequentialNumberAsync();
        Task<Invoice?> GetBySequentialNumberAsync(int sequentialNumber);
        Task<IEnumerable<Invoice>> GetAllAsync();
        Task AddAsync(Invoice invoice);
        Task UpdateAsync(Invoice invoice);
    }
}
