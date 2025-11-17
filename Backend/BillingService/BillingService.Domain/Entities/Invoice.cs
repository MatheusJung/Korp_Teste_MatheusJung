using BillingService.Domain.Enums;
using BillingService.Domain.Exceptions;

namespace BillingService.Domain.Entities
{
    public sealed class Invoice
    {
        public int Id { get; private set; }
        public int SequentialNumber { get; private set; }
        public InvoiceStatus Status { get; private set; } = InvoiceStatus.Aberta;
        public DateTime CreatedAt { get; private set; } = DateTime.UtcNow;

        private readonly List<InvoiceItem> _items = new();
        public IReadOnlyCollection<InvoiceItem> Items => _items.AsReadOnly();

        private Invoice() { } // EF Core

        //Construtor — cria a nota sem itens ainda
        public Invoice(int number)
        {
            SequentialNumber = number;
        }

        public void AddItem(InvoiceItem item)
        {
            if (Status == InvoiceStatus.Fechada)
                throw new CannotAddItemToClosedInvoiceException(SequentialNumber);

            if (string.IsNullOrWhiteSpace(item.ProductCode))
                throw new ProductNotFoundException(item.ProductCode);

            if (item.Quantity <= 0)
                throw new InvalidQuantityException(item.Quantity);

            _items.Add(item);
            item.Invoice = this;
        }

        public void RemoveItem(string productCode)
        {
            if (Status != InvoiceStatus.Aberta)
                throw new CannotRemoveItemsFromClosedInvoiceException(SequentialNumber);

            var item = _items.FirstOrDefault(i => i.ProductCode == productCode) 
                ?? throw new ProductNotFoundException(productCode);
            _items.Remove(item);
        }

        public void Close()
        {
            if (Status == InvoiceStatus.Fechada)
                throw new InvoiceAlreadyClosedException(SequentialNumber);
            if (_items.Count == 0)
                throw new EmptyInvoiceException();

            Status = InvoiceStatus.Fechada;
        }

        public void Cancel()
        {
            if (Status == InvoiceStatus.Fechada)
                throw new CannotCancelClosedInvoiceException(SequentialNumber);
            if (Status == InvoiceStatus.Cancelada)
                throw new InvoiceAlreadyCanceledException(SequentialNumber);
            Status = InvoiceStatus.Cancelada;
        }
    }
}
