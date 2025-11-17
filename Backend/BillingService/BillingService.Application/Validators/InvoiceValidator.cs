using BillingService.Application.DTOs;
using BillingService.Domain.Entities;
using BillingService.Domain.Exceptions;

namespace BillingService.Application.Validators
{
    public static class InvoiceValidator
    {
        public static void ValidateNoDuplicateItems(IEnumerable<InvoiceItemDto> items)
        {
            var duplicatedProducts = items
                .GroupBy(i => i.ProductCode)
                .Where(g => g.Count() > 1)
                .Select(g => g.Key)
                .ToList();

            if (duplicatedProducts.Any())
            {
                var duplicatedList = string.Join(", ", duplicatedProducts);
                throw new DuplicateItemsInInvoiceException(duplicatedList);
            }
        }

        public static void EnsureInvoiceExists(Invoice? invoice, int sequentialNumber)
        {
            if (invoice is null)
                throw new DirectoryNotFoundException($"Invoice with Sequential Number {sequentialNumber} was not found.");
        }

        public static void ValidateQuantity(int quantity)
        {
            if (quantity <= 0)
                throw new ArgumentException("Quantidade inválida. Deve ser maior que zero.");
        }
    }
}
