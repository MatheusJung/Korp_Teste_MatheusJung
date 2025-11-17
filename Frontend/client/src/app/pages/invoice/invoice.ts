import { Component , OnInit } from '@angular/core';
import { Footer } from "../../shared/footer/footer";
import { Sidebar } from "../../shared/sidebar/sidebar";
import { Topbar } from "../../shared/topbar/topbar";
import { InvoiceService } from '../../core/Invoice/invoiceService';
import { InvoiceDto } from '../../core/Invoice/invoiceModel';
import { CommonModule } from '@angular/common';
import { InvoiceTable } from "./invoice-table/invoice-table";

@Component({
  selector: 'app-invoice',
  imports: [Footer, Sidebar, Topbar, CommonModule, InvoiceTable],
  templateUrl: './invoice.html',
  styleUrl: './invoice.scss',
})
export class InvoiceComponent implements OnInit{
    sidebarOpen = false;
    invoices: InvoiceDto[] = [];
  showModal = false;

  constructor(private invoiceService: InvoiceService) {}

  ngOnInit() {
  }

    handleToggleSidebar() {
    this.sidebarOpen = !this.sidebarOpen;
  }

    openModal() {
    this.showModal = true;
  }

  closeModal() {
    this.showModal = false;
  }
}
