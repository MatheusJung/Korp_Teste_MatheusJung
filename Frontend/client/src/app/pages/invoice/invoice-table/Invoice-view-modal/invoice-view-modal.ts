import { Component, EventEmitter, Input, Output, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { InvoiceDto } from '../../../../core/Invoice/invoiceModel';
import { InvoiceService } from '../../../../core/Invoice/invoiceService';
import { InvoiceAddItemModal } from "./invoice-add-item-modal/invoice-add-item-modal";

@Component({
  selector: 'app-invoice-view-modal',
  standalone: true,
  imports: [CommonModule, InvoiceAddItemModal],
  templateUrl: './invoice-view-modal.html',
  styleUrl: './invoice-view-modal.scss',
})
export class InvoiceViewModal implements OnChanges {
  @Input() visible = false;
  @Input() invoiceNumber: number | null = null;
  @Output() close = new EventEmitter<void>();
  @Output() invoiceLoaded = new EventEmitter<InvoiceDto>();

  invoice: InvoiceDto | null = null;
  loading = false;
  canCancel = false;
  canPrint = false;
  canAddItem = false;
  addItemModalVisible = false;

  // exibir mensagens de erro
  errorMessage: string | null = null;

  constructor(private invoiceService: InvoiceService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['invoiceNumber'] && this.invoiceNumber) {

      this.loadInvoice();
    }
  }

  // Método único para verificar o status e atualizar flags
  private updateButtonStates(status?: string) {
    this.canCancel = status !== 'Fechada' && status !== 'Cancelada';
    this.canPrint = status === 'Aberta';
    this.canAddItem = status === 'Aberta';
  }

  //Carregar nota
  loadInvoice() {
    if (!this.invoiceNumber) return;
    this.loading = true;

    this.invoiceService.getInvoice(this.invoiceNumber).subscribe({
      next: (data) => {
        this.invoice = data;
        this.invoiceLoaded.emit(data);
        this.loading = false;
        this.updateButtonStates(data.status);
      },
      error: (err) => {
        console.error('Erro ao carregar invoice', err);
        this.loading = false;
      },
    });
  }

  //Fechar modal
  onClose() {
    this.close.emit();
  }

  //Abrir modal de adicionar itens
  openAddItemModal() {
    this.addItemModalVisible = true;
  }

  //Fechar modal de adicionar itens
  closeAddItemModal() {
    this.addItemModalVisible = false;
  }

  //Adicionar item a lista
  onItemAdded() {
    this.addItemModalVisible = false; // fecha modal filho
    if (this.invoiceNumber) {
      console.log('onItemAdded: recarregando invoice', this.invoiceNumber);
      this.loadInvoice();
    } else {
      console.warn('onItemAdded: invoiceNumber está null');
    }
  }

  //Cancelar nota
  cancelInvoice(seqNumber: number | null | undefined) {
  if (!seqNumber) {
    console.warn('Número da nota inválido');
    return;
  }

  this.invoiceService.cancelInvoice(seqNumber).subscribe({
    next: (updatedInvoice) => {
      // Atualiza a nota no modal ou na lista, se necessário
      this.invoice = updatedInvoice;
      this.updateButtonStates(updatedInvoice.status);
    },
    error: (err) => console.error('Erro ao cancelar nota', err),
  });
}

  // Fechar nota
  closeInvoice(seqNumber: number | null | undefined) {
    if (!seqNumber) {
      console.warn('Número da nota inválido');
      return;
    }

    this.invoiceService.closeInvoice(seqNumber).subscribe({
      next: (updatedInvoice) => {
        this.invoice = updatedInvoice;
        this.updateButtonStates(updatedInvoice.status);
      },
      error: (err) => console.error('Erro ao fechar nota', err),
    });
  }

  //Imprimir e Fechar nota
  printAndClose(seqNumber: number) {
    this.invoiceService.closeAndPrint(seqNumber).subscribe({
      next: (pdfBlob) => {
        // Criar URL do Blob e abrir para download
        const url = window.URL.createObjectURL(pdfBlob);
        //Abrir em nova aba
        window.open(url, '_blank');

        // Download
        //const a = document.createElement('a');
        // a.href = url;
        //a.download = `NF-${seqNumber}.pdf`;
        //a.click();
        //window.URL.revokeObjectURL(url);
      },
      error: (err) => {
        console.error('Erro ao fechar e imprimir NF', err);
        alert('Não foi possível gerar o PDF da NF.');
      },
    });
  }
}
