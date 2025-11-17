using BillingService.Application.DTOs;
using BillingService.Domain.Entities;

namespace BillingService.Application.Mappers
{
    public static class InvoiceMapper
    {
        public static InvoiceItemDto ToInvoiceItemDto(InvoiceItem item)
        {
            if (item == null)
                throw new InvalidOperationException("Item da NF não pode ser nulo.");
            if (string.IsNullOrWhiteSpace(item.ProductCode))
                throw new InvalidOperationException("Código do produto inválido.");

            return new InvoiceItemDto(
                item.ProductCode,
                item.Description,
                item.Quantity
            );
        }

        public static InvoiceDto ToInvoiceDto(this Invoice invoice)
        {
            if (invoice == null)
                throw new InvalidOperationException("A nota fiscal não pode ser nula.");
            if (invoice.SequentialNumber <= 0)
                throw new InvalidOperationException("Número sequencial inválido.");

            var itemDtos = invoice.Items.Select(i => new InvoiceItemDto(
                i.ProductCode,
                i.Description,
                i.Quantity
            ));

            return new InvoiceDto(
                invoice.SequentialNumber,
                invoice.CreatedAt,
                invoice.Status.ToString(),
                itemDtos
            );
        }
        public static IEnumerable<InvoiceDto> ToInvoiceDtoList(this IEnumerable<Invoice> invoices)
        {
            return invoices.Select(invoice => invoice.ToInvoiceDto());
        }
    }
}
