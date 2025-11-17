namespace BillingService.Application.DTOs
{
    public sealed record CreateInvoiceDto(
        List<InvoiceItemDto> Items
    );
}
