import { ProductDto } from '../../../../core/product/productModel';
import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ProductService } from '../../../../core/product/productService';


@Component({
  selector: 'app-product-add-modal',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './product-add-modal.html',
  styleUrl: './product-add-modal.scss',
})

export class ProductAddModal {
  @Input() visible = false;
  @Output() close = new EventEmitter<void>();

  // Produto a ser cadastrado
  newItem: ProductDto = {productCode: '', description: '', quantity: 0, status: 'False'
  };

  productCreated: ProductDto | null = null;
  errorMessage: string | null = null;

  constructor(private productService: ProductService) {}

  /** Abre modal e limpa tudo */
  onOpenModal() {
    this.visible = true;
    this.resetFields();
  }

  /** Fecha modal e avisa componente pai */
  onCloseModal() {
    this.visible = false;
    this.resetFields();
    this.close.emit();
  }

  /** Reseta formulário */
  private resetFields() {
    this.newItem = {productCode: '',description: '',quantity: 0, status: 'False'
    };
    this.productCreated = null;
    this.errorMessage = null;
  }

  /** Envia novo produto ao backend */
 onRegisterNewProduct() {
    if (!this.newItem.productCode || !this.newItem.description) {
      this.errorMessage = "Preencha todos os campos obrigatórios.";
      return;
    }

    this.errorMessage = null; // limpa erro anterior

    this.productService.createProduct({
      code: this.newItem.productCode,
      name: this.newItem.description,
      initialStock: this.newItem.quantity
    }).subscribe({
      next: (productDto) => {
        console.log('Produto criado:', productDto);
        this.productCreated = productDto;
        this.onCloseModal();
      },
      error: (err) => {
        console.error('Erro ao criar produto:', err);

        // ⬇️ PEGAR ERRO DA API
        this.errorMessage = "Erro inesperado ao criar o produto.";
      }
    });
  }
}
