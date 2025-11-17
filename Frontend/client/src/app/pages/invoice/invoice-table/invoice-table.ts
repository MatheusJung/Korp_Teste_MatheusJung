import { InvoiceDto, PagedInvoiceResult } from '../../../core/Invoice/invoiceModel';
import { Component, OnInit} from '@angular/core';
import { CommonModule } from '@angular/common';
import { InvoiceService } from '../../../core/Invoice/invoiceService';
import { FormsModule } from '@angular/forms';
import { InvoiceViewModal } from "./Invoice-view-modal/invoice-view-modal";
import { InvoiceAddInvoiceModal } from "./invoice-add-Invoice-modal/invoice-add-Invoice-modal";

@Component({
  selector: 'app-invoice-table',
  imports: [CommonModule, FormsModule, InvoiceViewModal, InvoiceAddInvoiceModal],
  templateUrl: './invoice-table.html',
  styleUrl: './invoice-table.scss',
})

export class InvoiceTable implements OnInit {
  pagedInvoices?: PagedInvoiceResult;
  pageNumber = 1;
  pageSize = 10;
  isLoading = false;
  invoices: InvoiceDto[] = [];
  modalViewVisible = false;
  modalAddInvoiceVisible = false;
  selectedInvoice: number | null = null;

  constructor(private invoiceService: InvoiceService) {}

  ngOnInit() {
    this.loadInvoices();
  }

  loadInvoices() {
    this.isLoading = true;
    this.invoiceService.getInvoicesPaged(this.pageNumber, this.pageSize).subscribe({
      next: (data) => {
        this.pagedInvoices = data;
        this.isLoading = false;
        console.log('Dados recebidos do backend:', data);
      },
      error: (err) => {
        console.error('Erro ao carregar notas', err);
        this.isLoading = false;
      },
    });
  }

  nextPage() {
    if (this.pageNumber < (this.pagedInvoices?.totalPages ?? 1)) {
      this.pageNumber++;
      this.loadInvoices();
    }
  }

  prevPage() {
    if (this.pageNumber > 1) {
      this.pageNumber--;
      this.loadInvoices();
    }
  }

  changePageSize(size: number) {
    this.pageSize = size;
    this.pageNumber = 1;
    this.loadInvoices();
  }

  cancelInvoice(seqNumber: number) {
    this.invoiceService.cancelInvoice(seqNumber).subscribe({
      next: (updatedInvoice) => {
        const index = this.pagedInvoices?.items.findIndex(i => i.sequentialNumber === updatedInvoice.sequentialNumber);
        if (index !== undefined && index >= 0) this.pagedInvoices!.items[index] = updatedInvoice;
      },
      error: (err) => console.error('Erro ao cancelar nota', err),
    });
  }

  closeInvoice(seqNumber: number) {
    this.invoiceService.closeInvoice(seqNumber).subscribe({
      next: (updatedInvoice) => {
        const index = this.pagedInvoices?.items.findIndex(i => i.sequentialNumber === updatedInvoice.sequentialNumber);
        if (index !== undefined && index >= 0) this.pagedInvoices!.items[index] = updatedInvoice;
      },
      error: (err) => console.error('Erro ao fechar nota', err),
    });
  }

  openViewModal(seqNumber: number) {
  this.selectedInvoice = seqNumber;
  this.modalViewVisible = true;
  }

  closeViewModal(seqNumber: number | null) {
    this.modalViewVisible = false;
    this.selectedInvoice = null;
    this.loadInvoices();
  }

  onInvoiceLoaded(invoice: InvoiceDto) {
    console.log('Invoice carregada:', invoice);
  }

  openAddInvoiceModal() {
  this.modalAddInvoiceVisible = true;
  }

  closeAddInvoiceModal() {
    this.modalAddInvoiceVisible = false;
    this.loadInvoices();
  }

  isPrintable(status: string): boolean {
    return status === 'Aberta';
  }

  //Imprimir e Fechar
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
        this.loadInvoices();

      },
      error: (err) => {
        console.error('Erro ao fechar e imprimir NF', err);
        alert('Não foi possível gerar o PDF da NF.');
      },
    });
  }
}
