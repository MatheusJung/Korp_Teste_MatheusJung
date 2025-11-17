
namespace BillingService.Application.DTOs
{
    public sealed record InvoiceDto(
        int SequentialNumber,
        DateTime CreatedAt,
        string Status,
        IEnumerable<InvoiceItemDto> Items
    );
}
