import { Component, OnInit, Input } from '@angular/core';
import { ProductService } from '../../../core/product/productService';
import { CommonModule } from '@angular/common';
import { ProductDto } from '../../../core/product/productModel';
import { PagedProductResult } from '../../../core/product/productModel';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ProductAddModal } from "./product-add-modal/product-add-modal";
import { ProductAddStockModal } from "./product-add-stock-modal/product-add-stock-modal";


@Component({
  imports: [CommonModule, FormsModule, ProductAddModal,ProductAddStockModal],

  selector: 'app-product-table',
  templateUrl:'./product-table.html',
  styleUrl: './product-table.scss',
})
export class ProductTable implements OnInit {
 pagedProducts: PagedProductResult | null = null;
  isLoading = false;
  modalAddStockVisible = false;
  modalAddVisible = false;
  selectedProduct: string | null = null;

  pageNumber = 1;
  pageSize = 10;

  constructor(private productService: ProductService) {}

  ngOnInit(): void {
    this.OnloadProducts();
  }

  OnloadProducts() {
    this.isLoading = true;
    this.productService.getProducts(this.pageNumber, this.pageSize).subscribe({
      next: (res: PagedProductResult) => {
        this.pagedProducts = res;
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Erro ao carregar produtos:', err);
        this.isLoading = false;
      }
    });
  }

  prevPage() {
    if (this.pageNumber > 1) {
      this.pageNumber--;
      this.OnloadProducts();
    }
  }

  nextPage() {
    if (this.pagedProducts && this.pageNumber < this.pagedProducts.totalPages) {
      this.pageNumber++;
      this.OnloadProducts();
    }
  }

  changePageSize(size: number) {
    this.pageSize = size;
    this.pageNumber = 1; // reset para a primeira página
    this.OnloadProducts();
  }

  deactivateProduct(code: string) {
    this.productService.deactivateProduct(code).subscribe({
      next: () => {
        alert('Produto desativado');
        this.OnloadProducts();
      },
      error: (err) => alert(err.error?.error || 'Erro ao desativar produto')
    });
  }

  addStock(code: string, quantity: number) {
    this.productService.addStock({ productCode: code, quantity, operationKey:Date.now().toString() }).subscribe({
      next: (res) => {
        alert(res.message);
        this.OnloadProducts();
      },
      error: (err) => alert(err.error?.error || 'Erro ao adicionar estoque')
    });
  }

  removeStock(code: string, quantity: number) {
    this.productService.removeStock({ productCode: code, quantity, operationKey:Date.now().toString() }).subscribe({
      next: (res) => {
        alert(res.message);
        this.OnloadProducts();
      },
      error: (err) => alert(err.error?.error || 'Erro ao remover estoque')
    });
  }

    openAddStockModal(productCode: string) {
    this.selectedProduct = productCode;
    this.modalAddStockVisible = true;
    }

    closeAddStockModal(productCode: string | null) {
      this.modalAddStockVisible = false;
      this.selectedProduct = null;
      this.OnloadProducts();
    }

    onProductLoaded(invoice: ProductDto) {
      console.log('Produto carregada:', invoice);
    }

    openAddModal() {
    this.modalAddVisible = true;
    }

    closeAddModal() {
      this.modalAddVisible = false;
       this.OnloadProducts();
    }

    isActivated(status: string): boolean {
      return status === 'True';
    }

    onProductsLoaded(product: ProductDto) {
      console.log('Produto carregado:', product);
    }
}
