import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';
import { HealthService } from '../healthService';
import { map, catchError, timeout, tap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';
import { environment } from '../enviroment';

@Injectable({ providedIn: 'root' })
export class ProductOnlineGuard implements CanActivate {
   constructor(
    private healthService: HealthService,
    private router: Router
  ) {}

 canActivate(): Observable<boolean> {
    return this.healthService.getServiceStatus(`${environment.inventoryServiceUrl}/health`, 1000)
      .pipe(
        timeout(2000),
        tap(isOnline => console.log('Invoice Guard: serviço online?', isOnline)),
        map(isOnline => {
          if (!isOnline) {
            console.log('Product Guard: redirecionando para service-offline');
            this.router.navigate(['/erro-service-offline']);
          }
          return isOnline;
        }),
        catchError(err => {
          console.log('Product Guard: erro ao checar health', err);
          this.router.navigate(['/erro-service-offline']);
          return of(false);
        })
      );
  }
}
