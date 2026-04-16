import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { interval, Subscription, catchError, of, map } from 'rxjs';
import { ServiceStatus } from '../models/models';

export interface ServicesHealth {
  estoque: ServiceStatus;
  faturamento: ServiceStatus;
}

@Injectable({ providedIn: 'root' })
export class HealthService {
  private readonly http = inject(HttpClient);

  readonly health = signal<ServicesHealth>({ estoque: 'ok', faturamento: 'ok' });

  private sub?: Subscription;

  startMonitoring(): void {
    this.verificar();
    this.sub = interval(30_000).subscribe(() => this.verificar());
  }

  stopMonitoring(): void {
    this.sub?.unsubscribe();
  }

  private verificar(): void {
    this.http.get<{ status: string }>('/health/estoque').pipe(
      map(r => this.normalizar(r.status)),
      catchError(() => of('down' as ServiceStatus))
    ).subscribe(status => {
      this.health.update(h => ({ ...h, estoque: status }));
    });

    this.http.get<{ status: string }>('/health/faturamento').pipe(
      map(r => this.normalizar(r.status)),
      catchError(() => of('down' as ServiceStatus))
    ).subscribe(status => {
      this.health.update(h => ({ ...h, faturamento: status }));
    });
  }

  /**
   * Normaliza diferentes formatos de status para o padrão interno.
   * .NET retorna: Healthy / Degraded / Unhealthy
   * Go retorna:   ok / degraded / down
   */
  private normalizar(raw: string): ServiceStatus {
    switch (raw?.toLowerCase()) {
      case 'healthy':
      case 'ok':
        return 'ok';
      case 'degraded':
        return 'degraded';
      case 'unhealthy':
      case 'down':
      default:
        return 'down';
    }
  }
}
