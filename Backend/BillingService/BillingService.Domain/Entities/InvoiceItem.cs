using BillingService.Domain.Exceptions;
using System.Text.Json.Serialization;

namespace BillingService.Domain.Entities{ 

public sealed class InvoiceItem
    {
        public int Id { get; private set; }
        public string ProductCode { get; private set; } = string.Empty;
        public string Description { get; private set; } = string.Empty;
        public int Quantity { get; private set; }

        public int InvoiceId { get; set; }
        [JsonIgnore]
        public Invoice? Invoice { get; set; }

        private InvoiceItem() { } // EF Core

        public InvoiceItem(string productCode, string description, int quantity, Invoice invoice)
        {
            if (string.IsNullOrWhiteSpace(productCode))
                throw new ProductNotFoundException(productCode);

            if (quantity <= 0)
                throw new InvalidQuantityException(quantity);

            ProductCode = productCode;
            Description = description;
            Quantity = quantity;

            Invoice = invoice;
        }
    }
}