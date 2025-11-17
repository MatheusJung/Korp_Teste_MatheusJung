using BillingService.Domain.Entities;
using System.Collections.Generic;

namespace BillingService.Domain.Exceptions
{

    public class InvoiceException : DomainException
    {
        public InvoiceException(string message) : base(message) { }

    }
    // Nota já fechada
    public class InvoiceAlreadyClosedException : InvoiceException
    {
        public InvoiceAlreadyClosedException(int seqNumber) :
            base($"Invoice: {seqNumber} is already closed.") { }

    }
    // Nota já cancelada
    public class InvoiceAlreadyCanceledException : InvoiceException
    {
        public InvoiceAlreadyCanceledException(int seqNumber) :
            base($"Invoice:{seqNumber} is already closed.") { }

    }
    // Nota não encontrada
    public class InvoiceNotFoundException : InvoiceException
    {
        public InvoiceNotFoundException(int seqNumber) :
            base($"Invoice: {seqNumber} was not found.") { }

    }
    //  Tentativa de imprimir nota fechada
    public class CannotPrintClosedInvoiceException : InvoiceException
    {
        public CannotPrintClosedInvoiceException(int seqNumber) :
            base($"Cannot print invoice {seqNumber} because it is already closed.") { }

    }
    // Tentativa de adicionar item em nota fechada
    public class CannotAddItemToClosedInvoiceException : InvoiceException
    {
        public CannotAddItemToClosedInvoiceException(int seqNumber) : 
            base($"Cannot add items to closed invoice {seqNumber}.") { }

    }
    // Produto não encontrado
    public class ProductNotFoundException : InvoiceException
    {
        public ProductNotFoundException(string productCode) :
            base($"Product: {productCode} was not found in InventoryService."){ }

    }
    // Quantidade inválida
    public class InvalidQuantityException : InvoiceException
    {
        public InvalidQuantityException(int quantity)
            : base($"Quantity must be greater than zero. Received: {quantity}.") { }
    }

    // Estoque insuficiente
    public class InsufficientStockException : InvoiceException
    {
        public InsufficientStockException(string productCode, int requested, int available)
            : base($"Not enough stock for product '{productCode}'. Requested: {requested}, Available: {available}.") { }
    }

    // Produto inativo
    public class ProductInactiveException : InvoiceException
    {
        public ProductInactiveException(string productCode)
            : base($"Product:{productCode}: is inactive and cannot be added to an invoice.") { }
    }

    // Nota sem itens
    public class EmptyInvoiceException : InvoiceException
    {
        public EmptyInvoiceException()
            : base("An invoice must contain at least one item.") { }
    }

    // Tentativa de modificar nota fechada
    public class CannotModifyClosedInvoiceException : InvoiceException
    {
        public CannotModifyClosedInvoiceException(int seqNumber)
            : base($"Cannot modify a closed invoice:{seqNumber}.") { }
    }

    // Tentativa de cancelar nota fechada
    public class CannotCancelClosedInvoiceException : InvoiceException
    {
        public CannotCancelClosedInvoiceException(int seqNumber)
            : base($"Cannot cancel a closed invoice:{seqNumber}") { }
    }

    // Tentativa de remover item de nota fechada
    public class CannotRemoveItemsFromClosedInvoiceException : InvoiceException
    {
        public CannotRemoveItemsFromClosedInvoiceException(int seqNumber)
            : base($"Cannot remove items from a closed invoice: {seqNumber}") { }
    }

    // Itens duplicados na nota
    public class DuplicateItemsInInvoiceException : InvoiceException
    {
        public DuplicateItemsInInvoiceException(string duplicatedList)
            : base($"Existem produtos duplicados na nota: {duplicatedList}") { }
    }

    // Itens duplicados na nota
    public class FailedToCreateInvoiceException : InvoiceException
    {
        public FailedToCreateInvoiceException(int seqNumber)
            : base($"Failed to create invoice:{seqNumber}") { }
    }
}
