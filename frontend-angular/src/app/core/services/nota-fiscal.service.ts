import { Injectable, inject } from "@angular/core";
import { HttpClient, HttpHeaders, HttpParams } from "@angular/common/http";
import {
  Observable,
  interval,
  switchMap,
  takeWhile,
  startWith,
  tap,
  catchError,
  of,
  throwError,
} from "rxjs";
import {
  NotaFiscal,
  CriarNotaRequest,
  ImprimirResponse,
  PagedResult,
  ListarParams,
} from "../models/models";

@Injectable({ providedIn: "root" })
export class NotaFiscalService {
  private readonly http = inject(HttpClient);
  private readonly base = "/api/faturamento/notas";

  // 🔥 cache em memória
  private cacheNotas: NotaFiscal[] = [];

  private readonly CACHE_KEY = "notas_cache";

  // =========================
  // LISTAR (com cache)
  // =========================
  listar(): Observable<NotaFiscal[]> {
    return this.http.get<NotaFiscal[]>(this.base).pipe(
      tap((notas) => {
        this.cacheNotas = notas;
        localStorage.setItem(this.CACHE_KEY, JSON.stringify(notas));
      }),
      catchError(() => {
        // tenta memória
        if (this.cacheNotas.length > 0) {
          return of(this.cacheNotas);
        }

        // tenta localStorage
        const cache = localStorage.getItem(this.CACHE_KEY);
        if (cache) {
          const notas = JSON.parse(cache) as NotaFiscal[];
          this.cacheNotas = notas;
          return of(notas);
        }

        // fallback vazio
        return of([]);
      }),
    );
  }

  listarPaginado(params?: ListarParams): Observable<PagedResult<NotaFiscal>> {
    let httpParams = new HttpParams();

    if (params?.page != null) {
      httpParams = httpParams.set("page", params.page);
    }

    if (params?.pageSize != null) {
      httpParams = httpParams.set("pageSize", params.pageSize);
    }

    if (params?.search) {
      httpParams = httpParams.set("search", params.search);
    }

    if (params?.sortBy) {
      httpParams = httpParams.set("sortBy", params.sortBy);
    }

    if (params?.sortDirection) {
      httpParams = httpParams.set("sortDirection", params.sortDirection);
    }

    return this.http.get<PagedResult<NotaFiscal>>(`${this.base}/paginado`, {
      params: httpParams,
    });
  }

  // =========================
  // OBTER (com cache leve)
  // =========================
  obter(id: string): Observable<NotaFiscal> {
    return this.http.get<NotaFiscal>(`${this.base}/${id}`).pipe(
      catchError(() => {
        const nota = this.cacheNotas.find((n) => n.id === id);
        if (nota) return of(nota);

        return throwError(() => new Error("Nota não encontrada no cache"));
      }),
    );
  }

  // =========================
  // AÇÕES (SEM CACHE)
  // =========================
  criar(req: CriarNotaRequest): Observable<NotaFiscal> {
    return this.http.post<NotaFiscal>(this.base, req);
  }

  imprimir(id: string): Observable<ImprimirResponse> {
    const headers = new HttpHeaders({
      "Idempotency-Key": crypto.randomUUID(),
    });

    return this.http.post<ImprimirResponse>(
      `${this.base}/${id}/imprimir`,
      {},
      { headers },
    );
  }

  // =========================
  // POLLING
  // =========================
  aguardarProcessamento(id: string): Observable<NotaFiscal> {
    return interval(3000).pipe(
      startWith(0),
      switchMap(() => this.obter(id)),
      takeWhile((nota) => nota.status === "Processando", true),
    );
  }
}
