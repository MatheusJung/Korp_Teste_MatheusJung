using BillingService.Domain.Entities;

namespace BillingService.Application.DTOs
{
    public sealed record InvoiceItemDto(
        string ProductCode,
        string Description,
        int Quantity
    );
}
