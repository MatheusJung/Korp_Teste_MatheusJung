using BillingService.Domain.Entities;

namespace BillingService.Application.DTOs
{
    public sealed record ItemDto(
        string ProductCode,
        string Description,
        int Quantity,
        string Status
    );
}
