import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of, interval } from 'rxjs';
import { catchError, map, retry, timeout, startWith, switchMap} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class HealthService {


  constructor(private http: HttpClient) {}

getServiceStatus(healthUrl: string): Observable<boolean> {
  return this.http.get(healthUrl).pipe(
    timeout(3000),       // timeout de 3s
    retry(2),            // tenta 2 vezes em caso de falha
    map(() => true),     // sucesso = online
    catchError(() => of(false)) // falha = offline
  );
}

getServiceStatusPoll(healthUrl: string, intervalMs: number = 3000): Observable<boolean> {
    return interval(intervalMs).pipe(
      startWith(0), // dispara imediatamente
      switchMap(() =>
        this.http.get(healthUrl).pipe(
          timeout(3000),
          retry(2),
          map(() => true),
          catchError(() => of(false))
        )
      )
    );
  }
}
