import { Injectable } from '@angular/core';
import { environment } from './enviroment';

@Injectable({providedIn: 'root'})
export class ApiConfigService  {
  billingHealthUrl = environment.billingServiceUrl;
  inventoryHealthUrl = environment.inventoryServiceUrl;
}
