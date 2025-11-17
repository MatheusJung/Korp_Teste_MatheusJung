import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { InvoiceService } from '../../../../../core/Invoice/invoiceService';
import { ProductDto } from '../../../../../core/product/productModel';
import { ProductService } from '../../../../../core/product/productService';

@Component({
  selector: 'app-invoice-add-item-modal',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './invoice-add-item-modal.html',
  styleUrls: ['./invoice-add-item-modal.scss'],
})

export class InvoiceAddItemModal{
  @Input() visible = false;
  @Input() invoiceNumber: number | null = null;
  @Output() close = new EventEmitter<void>();
  @Output() itemAdded = new EventEmitter<void>();

  // item completo
  newItem = {
    productCode: '',
    description: '',
    quantityAvailable: 0,
    quantityRequested: 1
};

  // exibir mensagens de erro
  errorMessage: string | null = null;

  constructor(
    private invoiceService: InvoiceService,
    private productService: ProductService
  ) {}

  /** Fecha modal */
  onCloseModal() {
    this.resetModal();
    this.close.emit();
  }

  /** Limpa o modal */
  resetModal(){
    this.newItem.productCode = '';
    this.newItem.description = '';
    this.newItem.quantityAvailable = 0;
    this.newItem.quantityRequested = 1;
    this.errorMessage = null;
  }

  /** Quando digitar o código, buscar produto */
  onProductCodeChange() {
    this.errorMessage = null;

    const code = this.newItem.productCode;
    if (!code || code.trim().length === 0) {
      return;
    }

    this.productService.getProductByCode(code).subscribe({
      next: (product: ProductDto) => {
        if (product.status !== 'True') {
          this.errorMessage = "Produto está desativado.";
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
      }
    });
  }


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

    if (!this.invoiceNumber) {
      this.errorMessage = "NF inválida. Reabra o modal.";
      return;
    }

    const payload = {
      productCode: this.newItem.productCode,
      description: this.newItem.description,
      quantity: this.newItem.quantityRequested
    };

    this.invoiceService.addItem(this.invoiceNumber, payload).subscribe({
      next: () => {
        this.itemAdded.emit();
        this.onCloseModal();
      },
      error: err => {
        this.errorMessage = err.error?.message ?? "Erro ao adicionar item.";
      }
    });
  }
}
