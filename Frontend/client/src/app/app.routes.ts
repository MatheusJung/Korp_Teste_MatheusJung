import { Routes } from '@angular/router';
import { InvoiceOnlineGuard } from './core/Invoice/invoice-online-guard';
import { ProductOnlineGuard } from './core/product/product-online-guard';

export const routes: Routes = [
    {
    path: '',
    loadComponent: () =>import('./pages/home/home').then(m => m.Home)
  },
  { path: 'invoice',
    canActivate: [InvoiceOnlineGuard],
    loadComponent: () => import('./pages/invoice/invoice').then(m => m.InvoiceComponent)
  },
  { path: 'products',
    canActivate: [ProductOnlineGuard],
    loadComponent: () => import('./pages/product/product').then(m => m.ProductComponent)
  },
  {
    path: 'erro-service-offline',
    loadComponent: () => import('./pages/service-offline/service-offline').then(m => m.ServiceOfflineComponent)
  },
  { path: '**', redirectTo: '' }];
