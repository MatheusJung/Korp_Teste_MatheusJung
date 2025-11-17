import { InvoiceService } from './../../core/Invoice/invoiceService';
import { environment } from './../../core/enviroment';
import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Button } from "../button/button";
import { HealthService } from '../../core/healthService';

@Component({
  selector: 'app-menu-options',
  imports: [Button, CommonModule],
  templateUrl: './menu-options.html',
  styleUrl: './menu-options.scss',
})
export class MenuOptions implements OnInit{
  invoiceOnline:boolean = false;
  productOnline:boolean = false;

  invoiceService:string;
  productService:string;

  constructor(private healthService: HealthService)
  {
    this.invoiceService = `${environment.billingServiceUrl}/health`
    this.productService = `${environment.inventoryServiceUrl}/health`
  }

  ngOnInit() {

    // botão começa desabilitado
    this.invoiceOnline = false;
    this.productOnline = false;

   // Verificar Billing (Invoice)
    this.healthService
      .getServiceStatus(this.invoiceService)
      .subscribe(status => this.invoiceOnline = status);

    // Verificar Inventory (Product)
    this.healthService
      .getServiceStatus(this.productService)
      .subscribe(status => this.productOnline = status);
  }
}
