import { AdjustStockDto, ProductDto } from './../../../../core/product/productModel';
import { Component, EventEmitter, Input, Output, OnChanges, SimpleChanges} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ProductService } from '../../../../core/product/productService';
import { FormsModule } from "@angular/forms";
import { finalize } from 'rxjs';

@Component({
  selector: 'app-product-add-stock-modal',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl:'./product-add-stock-modal.html',
  styleUrl: './product-add-stock-modal.scss'
})
export class ProductAddStockModal implements OnChanges{
  @Input() visible = false;
  @Input() productCode: string | null = null;
  @Output() close = new EventEmitter<void>();
  @Output() productLoaded = new EventEmitter<ProductDto>();

  newItem = {
    productCode: '',
    description: '',
    quantityAvailable: 0,
    quantityRequested: 0,
    movementkey:'',
    operationType:'Add'
  };

  loading = false;
  errorMessage: string | null = null;

  constructor(private productService: ProductService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['visible'] && this.visible === true) {
      if (this.productCode) {
        this.loadProduct(this.productCode);
      }
    }
  }

  onOpenModal() {
    this.visible = true;
    if (this.productCode) {
      this.loadProduct(this.productCode);
    }
  }

  onCloseModal() {
    this.visible = false;
    this.resetFields();
    this.close.emit();
  }

  private resetFields() {
    this.newItem = {
      productCode: '',
      description: '',
      quantityAvailable: 0,
      quantityRequested: 0,
      movementkey:'',
      operationType: 'Add'
    };
    this.errorMessage = null;
  }

  private loadProduct(code: string) {
    this.loading = true;

    this.productService.getProductByCode(code).subscribe({
      next: (product) => {
        this.newItem.productCode = product.productCode;
        this.newItem.description = product.description;
        this.newItem.quantityAvailable = product.quantity;

        this.productLoaded.emit(product);
        this.loading = false;
      },
      error: (err) => {
        this.errorMessage = "Não foi possível carregar o produto.";
        this.loading = false;
        console.error(err);
      }
    });
  }

  /** Busca produto pelo código ao digitar */
  onProductCodeChange() {
    const code = this.newItem.productCode.trim();
    if (!code) {
      this.newItem.description = '';
      this.newItem.quantityAvailable = 0;
      return;
    }

    this.loading = true;
    this.productService.getProductByCode(code).subscribe({
      next: (product) => {
        this.newItem.description = product?.description || '';
        this.newItem.quantityAvailable = product?.quantity || 0;
        this.loading = false;
      },
      error: () => {
        this.newItem.description = '';
        this.newItem.quantityAvailable = 0;
        this.loading = false;
      }
    });
  }

/** Salvar movimentação no estoque */
onRegisterStockMovement() {
  if (!this.newItem.productCode || !this.newItem.description) {
    this.errorMessage = 'Preencha todos os campos obrigatórios.';
    return;
  }

  if (this.newItem.quantityRequested <= 0) {
    this.errorMessage = 'A quantidade solicitada deve ser maior que zero.';
    return;
  }

  // Validação para remoção
  if (this.newItem.operationType === 'Remove' &&
      this.newItem.quantityRequested > this.newItem.quantityAvailable) {
    this.errorMessage = `Não é possível remover mais do que ${this.newItem.quantityAvailable} unidades disponíveis.`;
    return;
  }

  this.errorMessage = null;
  this.loading = true;

  let request$;

  switch (this.newItem.operationType) {
    case 'Add':
      request$ = this.productService.addStock({
        productCode: this.newItem.productCode,
        quantity: this.newItem.quantityRequested,
        operationKey: `${Date.now()}`
      });
      break;

    case 'Remove':
      request$ = this.productService.removeStock({
        productCode: this.newItem.productCode,
        quantity: this.newItem.quantityRequested,
        operationKey: `${Date.now()}`
      });
      break;

    default:
      this.errorMessage = 'Operação inválida';
      this.loading = false;
      return;
  }

  request$.pipe(finalize(() => this.loading = false))
    .subscribe({
      next: () => {
        // Atualiza saldo no modal
        this.loadProduct(this.newItem.productCode);

        // Limpa a quantidade solicitada
        this.newItem.quantityRequested = 0;
      },
      error: () => {
        this.errorMessage = 'Erro inesperado ao movimentar item no estoque.';
      }
    });
  }
}
