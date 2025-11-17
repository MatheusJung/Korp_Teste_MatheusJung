using BillingService.Domain.Entities;
using BillingService.Domain.Interfaces;
using BillingService.Infrastructure.Data;
using Microsoft.EntityFrameworkCore;

namespace BillingService.Infrastructure.Repositories
{
    public class InvoiceRepository : IInvoiceRepository
    {
        private readonly BillingDbContext _context;
        public InvoiceRepository(BillingDbContext context)
        {
            _context = context;
        }
        public async Task AddAsync(Invoice invoice)
        {
            await _context.Invoices.AddAsync(invoice);
            await _context.SaveChangesAsync();
        }
        public async Task UpdateAsync(Invoice invoice)
        {
            _context.Invoices.Update(invoice);
            await _context.SaveChangesAsync();
        }
        public async Task<IEnumerable<Invoice>> GetAllAsync()
        {
            return await _context.Invoices
                .Include(i => i.Items)
                .ToListAsync();
        }
        public async Task<int> GetNextSequentialNumberAsync()
        {
            // Busca o maior número de nota existente, se não houver, começa em 1
            var maxNumber = await _context.Invoices.MaxAsync(i => (int?)i.SequentialNumber) ?? 0;
            return maxNumber + 1;
        }
        public async Task<Invoice?> GetBySequentialNumberAsync(int sequentialNumber)
        {
            return await _context.Invoices
                .Include(i => i.Items)
                .AsNoTracking()
                .FirstOrDefaultAsync(i => i.SequentialNumber == sequentialNumber);
        }
    }
}
