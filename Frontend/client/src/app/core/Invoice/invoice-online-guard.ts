import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';
import { HealthService } from '../healthService';
import { map, timeout, catchError, tap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';
import { environment } from '../enviroment';

@Injectable({ providedIn: 'root' })
export class InvoiceOnlineGuard implements CanActivate {
  constructor(
    private healthService: HealthService,
    private router: Router
  ) {}

 canActivate(): Observable<boolean> {
    return this.healthService.getServiceStatus(`${environment.billingServiceUrl}/health`, 1000)
      .pipe(
        timeout(2000),
        tap(isOnline => console.log('Invoice Guard: serviço online?', isOnline)),
        map(isOnline => {
          if (!isOnline) {
            this.router.navigate(['/erro-service-offline']);
          }
          return isOnline;
        }),
        catchError(err => {
          console.log('Invoice Guard: erro ao checar health', err);
          this.router.navigate(['/erro-service-offline']);
          return of(false);
        })
      );
  }
}
