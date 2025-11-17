import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { InvoiceDto,InvoiceItemDto } from '../../../../core/Invoice/invoiceModel';
import { InvoiceService } from '../../../../core/Invoice/invoiceService';
import { ProductDto } from '../../../../core/product/productModel';
import { ProductService } from '../../../../core/product/productService';

@Component({
  selector: 'app-invoice-add-Invoice-modal',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './invoice-add-Invoice-modal.html',
  styleUrl: './invoice-add-Invoice-modal.scss',
})

export class InvoiceAddInvoiceModal{
  @Input() visible = false;
  @Output() close = new EventEmitter<void>();

  // item completo
  newItem = {
    productCode: '',
    description: '',
    quantityAvailable: 0,
    quantityRequested: 1
  };

  // Notas fiscal e lista de itens da nota
  invoice:InvoiceDto | null = null
  items: InvoiceItemDto[] = [];

  // exibir mensagens de erro
  errorMessage: string | null = null;

  invoiceCreated: InvoiceDto | null = null;

  constructor(
    private invoiceService: InvoiceService,
    private productService: ProductService
  ) {}

  /** Reset modal */
  onOpenModal() {
    this.visible = true;
    this.items = [];
    this.invoiceCreated = null;
    this.errorMessage = null;

    this.newItem = {
      productCode: '',
      description: '',
      quantityAvailable: 0,
      quantityRequested: 1
    };
  }

  /** Fecha modal */
  onCloseModal() {
    this.visible = false;
    this.items = []; // limpa a lista de itens
    this.invoice = null;
    this.newItem = { productCode: '', description: '', quantityAvailable: 0, quantityRequested: 1};
    this.errorMessage = null;
    this.close.emit();
  }

    /** Limpa os campos do modal */
  clearFields() {
    this.newItem.description = '';
    this.newItem.quantityAvailable = 0;
    this.newItem.quantityRequested = 1;
  }

  /** Quando digitar o código, buscar produto */
  onProductCodeChange() {
    this.errorMessage = null;

    const code = this.newItem.productCode;
    if (!code || code.trim().length === 0) {
      this.clearFields();
      return;
    }

    this.productService.getProductByCode(code).subscribe({
      next: (product: ProductDto) => {
        if (product.status !== 'True') {
          this.errorMessage = "Produto está desativado.";
          this.clearFields();
          return;
        }

        // Produto ativo, preenche campos
        this.newItem.description = product.description;
        this.newItem.quantityAvailable = product.quantity;
        this.newItem.quantityRequested = 1;
        this.errorMessage = '';
      },
      error: () => {
        this.errorMessage = "Produto não encontrado.";
        this.clearFields();
      }
    });
  }

  /** Adicionar item à lista */
  onAddItem() {
    this.errorMessage = null;

    if (!this.newItem.productCode || !this.newItem.description) {
      this.errorMessage = "Código ou descrição inválidos.";
      return;
    }

    if (this.newItem.quantityRequested < 1) {
      this.errorMessage = "Quantidade deve ser maior que zero.";
      return;
    }

    if (this.newItem.quantityRequested > this.newItem.quantityAvailable) {
      this.errorMessage = "Quantidade solicitada maior que o estoque.";
      return;
    }

    // adiciona item compatível com InvoiceItemDto
    this.items.push({
      productCode: this.newItem.productCode,
      description: this.newItem.description,
      quantity: this.newItem.quantityRequested
    });

    // reset
    this.newItem = {
      productCode: '',
      description: '',
      quantityAvailable: 0,
      quantityRequested: 1
    };
  }

  /** Remover item da lista */
  removeItem(index: number) {
    this.items.splice(index, 1);
  }

  /** Cancelar nota */
  cancelInvoice(seqNumber: number | undefined) {
    if (!seqNumber) return;

    this.invoiceService.cancelInvoice(seqNumber).subscribe({
      next: (updatedInvoice) => {
        console.log('Nota cancelada:', updatedInvoice);
        // Atualiza o status local para refletir na tela
        if (this.invoiceCreated) {
          this.invoiceCreated.status = updatedInvoice.status;
        }
      },
      error: (err) => console.error('Erro ao cancelar nota', err)
    });
  }

  /** Criar nota */
  onSave() {
    if (this.items.length === 0) {
      this.errorMessage = "Adicione pelo menos um item.";
      return;
    }

    this.invoiceService.createInvoice(this.items).subscribe({
      next: (invoice) => {
        this.invoiceCreated = invoice;
        console.log('Nota criada:', invoice);

        this.loadInvoice(invoice.sequentialNumber);

      },
      error: (err) => {
        console.error(err);
        this.errorMessage = "Erro ao criar nota.";
      }
    });
  }

  loadInvoice(invoiceNumber:number) {
    if (!invoiceNumber) return;
    this.invoiceService.getInvoice(invoiceNumber).subscribe({
      next: (data) => {
        this.invoice = data;
        console.log("NF carregada:", data);
      },
      error: (err) => {
        console.error("Erro ao carregar nota:", err);
      }
    });
  }
}
