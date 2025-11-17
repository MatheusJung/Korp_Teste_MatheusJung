import { environment } from './../../core/enviroment';
import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Button } from "../button/button";
import { HealthService } from '../../core/healthService';
import { startWith,combineLatest,interval,switchMap } from 'rxjs';

@Component({
  selector: 'app-menu-options',
  imports: [Button, CommonModule],
  templateUrl: './menu-options.html',
  styleUrl: './menu-options.scss',
})

export class MenuOptions implements OnInit {
  invoiceOnline: boolean = false;
  productOnline: boolean = false;

  invoiceService: string;
  productService: string;

  constructor(private healthService: HealthService) {
    this.invoiceService = `${environment.billingServiceUrl}/health`
    this.productService = `${environment.inventoryServiceUrl}/health`
  }

ngOnInit() {
  this.healthService.getServiceStatusPoll(this.invoiceService)
    .subscribe(status => this.invoiceOnline = status);

  this.healthService.getServiceStatusPoll(this.productService)
    .subscribe(status => this.productOnline = status);
  }
}
