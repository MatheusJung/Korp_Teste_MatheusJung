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

  private readonly CACHE_KEY = "notas_cache";
  private readonly CACHE_PREFIX = "notas_cache_paginado_";

  // cache em memória da listagem simples
  private cacheNotas: NotaFiscal[] = [];

  // cache em memória da paginação
  private cachePaginado = new Map<string, PagedResult<NotaFiscal>>();

  listar(): Observable<NotaFiscal[]> {
    return this.http.get<NotaFiscal[]>(this.base).pipe(
      tap((notas) => {
        this.cacheNotas = notas;
        localStorage.setItem(this.CACHE_KEY, JSON.stringify(notas));
      }),
      catchError(() => {
        if (this.cacheNotas.length > 0) {
          return of(this.cacheNotas);
        }

        const cache = localStorage.getItem(this.CACHE_KEY);
        if (cache) {
          const notas = JSON.parse(cache) as NotaFiscal[];
          this.cacheNotas = notas;
          return of(notas);
        }

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

    const cacheKey = this.getPagedCacheKey(params);

    return this.http
      .get<PagedResult<NotaFiscal>>(`${this.base}/paginado`, {
        params: httpParams,
      })
      .pipe(
        tap((res) => {
          this.cachePaginado.set(cacheKey, res);
          localStorage.setItem(cacheKey, JSON.stringify(res));
        }),
        catchError(() => {
          const emMemoria = this.cachePaginado.get(cacheKey);
          if (emMemoria) {
            return of(emMemoria);
          }

          const local = localStorage.getItem(cacheKey);
          if (local) {
            const parsed = JSON.parse(local) as PagedResult<NotaFiscal>;
            this.cachePaginado.set(cacheKey, parsed);
            return of(parsed);
          }

          return of({
            items: [],
            page: params?.page ?? 1,
            pageSize: params?.pageSize ?? 10,
            totalItems: 0,
            totalPages: 1,
          });
        }),
      );
  }

  obter(id: string): Observable<NotaFiscal> {
    return this.http.get<NotaFiscal>(`${this.base}/${id}`).pipe(
      tap((nota) => this.atualizarCachesComNota(nota)),
      catchError(() => {
        const notaMemoria = this.cacheNotas.find((n) => n.id === id);
        if (notaMemoria) {
          return of(notaMemoria);
        }

        for (const pagina of this.cachePaginado.values()) {
          const nota = pagina.items.find((n) => n.id === id);
          if (nota) {
            return of(nota);
          }
        }

        const cacheLista = localStorage.getItem(this.CACHE_KEY);
        if (cacheLista) {
          const notas = JSON.parse(cacheLista) as NotaFiscal[];
          const nota = notas.find((n) => n.id === id);
          if (nota) {
            return of(nota);
          }
        }

        const notaLocalPaginado =
          this.buscarNotaNosCachesPaginadosDoLocalStorage(id);
        if (notaLocalPaginado) {
          return of(notaLocalPaginado);
        }

        return throwError(() => new Error("Nota não encontrada no cache"));
      }),
    );
  }

  criar(req: CriarNotaRequest): Observable<NotaFiscal> {
    return this.http
      .post<NotaFiscal>(this.base, req)
      .pipe(tap(() => this.invalidarCache()));
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

  aguardarProcessamento(id: string): Observable<NotaFiscal> {
    return interval(3000).pipe(
      startWith(0),
      switchMap(() => this.obter(id)),
      takeWhile((nota) => nota.status === "Processando", true),
    );
  }

  private getPagedCacheKey(params?: ListarParams): string {
    return `${this.CACHE_PREFIX}${JSON.stringify({
      page: params?.page ?? 1,
      pageSize: params?.pageSize ?? 10,
      search: params?.search ?? "",
      sortBy: params?.sortBy ?? "",
      sortDirection: params?.sortDirection ?? "asc",
    })}`;
  }

  private invalidarCache(): void {
    this.cacheNotas = [];
    this.cachePaginado.clear();

    localStorage.removeItem(this.CACHE_KEY);

    Object.keys(localStorage)
      .filter((key) => key.startsWith(this.CACHE_PREFIX))
      .forEach((key) => localStorage.removeItem(key));
  }

  private atualizarCachesComNota(nota: NotaFiscal): void {
    const idx = this.cacheNotas.findIndex((n) => n.id === nota.id);
    if (idx >= 0) {
      this.cacheNotas[idx] = nota;
      localStorage.setItem(this.CACHE_KEY, JSON.stringify(this.cacheNotas));
    }

    for (const [key, pagina] of this.cachePaginado.entries()) {
      const itemIndex = pagina.items.findIndex((n) => n.id === nota.id);
      if (itemIndex >= 0) {
        const atualizada: PagedResult<NotaFiscal> = {
          ...pagina,
          items: pagina.items.map((item) =>
            item.id === nota.id ? nota : item,
          ),
        };

        this.cachePaginado.set(key, atualizada);
        localStorage.setItem(key, JSON.stringify(atualizada));
      }
    }
  }

  private buscarNotaNosCachesPaginadosDoLocalStorage(
    id: string,
  ): NotaFiscal | null {
    for (const key of Object.keys(localStorage)) {
      if (!key.startsWith(this.CACHE_PREFIX)) continue;

      const raw = localStorage.getItem(key);
      if (!raw) continue;

      try {
        const pagina = JSON.parse(raw) as PagedResult<NotaFiscal>;
        const nota = pagina.items.find((n) => n.id === id);
        if (nota) return nota;
      } catch {
        // ignora cache inválido
      }
    }

    return null;
  }
}
