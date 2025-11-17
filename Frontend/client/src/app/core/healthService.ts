import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map, retry, timeout, shareReplay } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class HealthService {
  // Cache por URL
  private cacheMap: Map<string, Observable<boolean>> = new Map();

  constructor(private http: HttpClient) {}

  getServiceStatus(healthUrl: string, cacheMs: number = 5000): Observable<boolean> {
    // Retorna cache se existir
    if (this.cacheMap.has(healthUrl)) {
      return this.cacheMap.get(healthUrl)!;
    }

    const status$ = this.http.get(healthUrl).pipe(
      timeout(3000),  // 3s timeout
      retry(2),       // 2 tentativas
      map(() => true), // sucesso = online
      catchError(() => of(false)), // falha = offline
      shareReplay(1)
    );

    this.cacheMap.set(healthUrl, status$);

    // Limpa o cache após cacheMs
    setTimeout(() => this.cacheMap.delete(healthUrl), cacheMs);

    return status$;
  }
}
