using BillingService.Application.ContractsHttp;
using BillingService.Application.DTOs;
using BillingService.Application.Validators;
using BillingService.Domain.Entities;
using BillingService.Domain.Enums;
using BillingService.Domain.Exceptions;
using BillingService.Domain.Interfaces;
using static System.Runtime.InteropServices.JavaScript.JSType;

namespace BillingService.Application.Services
{

    public class InvoiceService
    {
        private readonly IInvoiceRepository _invoiceRepo;
        private readonly IHttpInventoryClient _inventoryClient;

        public InvoiceService(IInvoiceRepository invoiceRepo, IHttpInventoryClient inventoryClient)
        {
            _invoiceRepo = invoiceRepo;
            _inventoryClient = inventoryClient;
        }

        // Criar nova nota
        public async Task<Invoice> CreateAsync(List<InvoiceItemDto> items)
        {
            InvoiceValidator.ValidateNoDuplicateItems(items);

            foreach (var itemDto in items)
            {
                // Busca o produto no InventoryService
                var product = await _inventoryClient.GetProductByCodeAsync(itemDto.ProductCode)
                    ?? throw new ProductNotFoundException(itemDto.ProductCode);
               
                // Verifica se o produto está ativo
                if (product.Status == "False")
                    throw new ProductInactiveException(product.ProductCode);

                // Verifica se há estoque suficiente
                if (itemDto.Quantity > product.Quantity)
                    throw new InsufficientStockException(itemDto.ProductCode, itemDto.Quantity, product.Quantity);
            }

            // Gera o próximo número sequencial da nota
            var sequentialNumber = await _invoiceRepo.GetNextSequentialNumberAsync();
            // Cria a nota sem itens inicialmente
            var invoice = new Invoice(sequentialNumber);

            try
            {
                foreach (var itemDto in items)
                {
                    // Gera operation key determinística para o estoque
                    string operationKey = $"INV-{sequentialNumber}-ITEM-{itemDto.ProductCode}-QTY-{itemDto.Quantity}";

                    // Ajusta o estoque (débito)
                    await _inventoryClient.AdjustStockAsync(itemDto.ProductCode, -itemDto.Quantity, operationKey);

                    // Cria o item já associado à nota
                    var invoiceItem = new InvoiceItem(itemDto.ProductCode, itemDto.Description, itemDto.Quantity, invoice);

                    // Adiciona o item à nota
                    invoice.AddItem(invoiceItem);
                }

                // Salva a nota e todos os itens no banco
                await _invoiceRepo.AddAsync(invoice);

                return invoice;

            }
            catch
            {
                await CancelAsync(sequentialNumber);
                throw new FailedToCreateInvoiceException(sequentialNumber);
            }
        }

        // Adicionar itens a nota aberta
        public async Task AddItemAsync(int sequentialNumber, string productCode, int quantity)
        {
            // Busca a nota
            var invoice = await _invoiceRepo.GetBySequentialNumberAsync(sequentialNumber)
                          ?? throw new KeyNotFoundException($"Invoice {sequentialNumber} not found.");

            if (invoice.Status != InvoiceStatus.Aberta)
                throw new CannotAddItemToClosedInvoiceException(sequentialNumber);

            // Busca produto no InventoryService
            var product = await _inventoryClient.GetProductByCodeAsync(productCode)
                          ?? throw new KeyNotFoundException($"Product {productCode} not found.");

            if (product.Status == "False")
                throw new ProductInactiveException(product.ProductCode);

            if (product.Quantity < quantity)
                throw new InsufficientStockException(product.ProductCode, quantity, product.Quantity);

            // Gera operation key determinística
            string operationKey = $"INV-{sequentialNumber}-ITEM-{product.ProductCode}-QTY-{quantity}";

            // Ajusta o estoque (débito da quantidade do pedido)
            await _inventoryClient.AdjustStockAsync(product.ProductCode, -quantity, operationKey);

            // Cria o item com a quantidade do pedido e associa à nota
            var item = new InvoiceItem(product.ProductCode, product.Description, quantity, invoice);
            invoice.AddItem(item);

            // Atualiza a nota no repositório
            await _invoiceRepo.UpdateAsync(invoice);
        }

        // Cancelar nota aberta e devolver itens
        public async Task CancelAsync(int sequentialNumber)
        {
            var invoice = await _invoiceRepo.GetBySequentialNumberAsync(sequentialNumber)
                           ?? throw new InvoiceNotFoundException(sequentialNumber);

            if (invoice.Status == InvoiceStatus.Fechada)
                throw new InvoiceAlreadyClosedException(sequentialNumber);

            if (invoice.Status == InvoiceStatus.Cancelada)
                throw new InvoiceAlreadyCanceledException(sequentialNumber);

            foreach (var item in invoice.Items)
            {
                var description = $"INV-CANCEL-{sequentialNumber}-ITEM-{item.ProductCode}-QTY-{item.Quantity}";

                await _inventoryClient.AdjustStockAsync(
                    item.ProductCode,
                    item.Quantity,
                    description
                );
            }

            invoice.Cancel();
            await _invoiceRepo.UpdateAsync(invoice);
        }


        // Remover item de nota aberta
        public async Task RemoveItemAsync(int sequentialNumber, string productCode)
        {
            var invoice = await _invoiceRepo.GetBySequentialNumberAsync(sequentialNumber)
                          ?? throw new KeyNotFoundException($"Invoice {sequentialNumber} not found.");
            if (invoice.Status != InvoiceStatus.Aberta)
                throw new InvoiceAlreadyClosedException(sequentialNumber);

            var item = invoice.Items.FirstOrDefault(i => i.ProductCode == productCode)
                       ?? throw new ProductNotFoundException(productCode);
            await _inventoryClient.AdjustStockAsync(item.ProductCode, item.Quantity, $"INV-{invoice.SequentialNumber}-REMOVE-{item.ProductCode}");
            invoice.RemoveItem(productCode);
            await _invoiceRepo.UpdateAsync(invoice);
        }

        // Fechar nota
        public async Task CloseAsync(int sequentialNumber)
        {
            var invoice = await _invoiceRepo.GetBySequentialNumberAsync(sequentialNumber)
                          ?? throw new KeyNotFoundException($"Invoice {sequentialNumber} not found.");
            if (invoice.Status == InvoiceStatus.Fechada)
                throw new InvoiceAlreadyClosedException(sequentialNumber);
            if (invoice.Status == InvoiceStatus.Cancelada)
                throw new CannotCancelClosedInvoiceException(sequentialNumber);

            invoice.Close();
            await _invoiceRepo.UpdateAsync(invoice);

        }

        // Listar todas as notas
        public Task<IEnumerable<Invoice>> ListAsync() => _invoiceRepo.GetAllAsync();

        // Consultar nota por numero sequencial
        public Task<Invoice?> GetBySequentialNumberAsync(int sequentialNumber) =>
            _invoiceRepo.GetBySequentialNumberAsync(sequentialNumber);
    }
}