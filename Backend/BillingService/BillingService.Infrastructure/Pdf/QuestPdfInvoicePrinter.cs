using BillingService.Application.DTOs;
using BillingService.Application.Interfaces;
using QuestPDF.Fluent;

namespace BillingService.Infrastructure.Pdf
{

    public class QuestPdfInvoicePrinter : IInvoicePrinter
    {
        public byte[] Print(InvoiceDto invoice)
        {
            var doc = Document.Create(container =>
            {
                container.Page(page =>
                {
                    page.Margin(40);
                    page.Header().Text($"Nota Fiscal Nº {invoice.SequentialNumber}").FontSize(20).Bold();

                    page.Content().Column(col =>
                    {
                        col.Item().Text($"Data: {invoice.CreatedAt:dd/MM/yyyy}");
                        col.Item().Text($"Cliente: {"John Doe"}");
                        col.Item().Text($"Total: R$ 0,00");
                        col.Item().LineHorizontal(1);

                        col.Item().Table(table =>
                        {
                            table.ColumnsDefinition(columns =>
                            {
                                columns.ConstantColumn(80);
                                columns.RelativeColumn(80);
                                columns.ConstantColumn(80);
                                columns.ConstantColumn(80);
                            });

                            table.Header(header =>
                            {
                                header.Cell().Text("Código").Bold();
                                header.Cell().Text("Descrição").Bold();
                                header.Cell().Text("Quantidade").Bold();
                                header.Cell().Text("Valor").Bold();
                            });

                            foreach (var item in invoice.Items ?? Enumerable.Empty<InvoiceItemDto>())
                            {
                                table.Cell().Text(item.ProductCode ?? "");
                                table.Cell().Text(item.Description ?? "");
                                table.Cell().Text(item.Quantity.ToString());
                                table.Cell().Text($"R$ 0,00");
                            }
                        });
                    });
                });
            });

            return doc.GeneratePdf();
        }
    }
}