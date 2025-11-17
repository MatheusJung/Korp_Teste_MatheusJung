using BillingService.Application.DTOs;

namespace BillingService.Application.Interfaces
{
    public interface IInvoicePrinter
    {
        byte[] Print(InvoiceDto invoice);
    }
}
