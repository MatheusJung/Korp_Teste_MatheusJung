import { Injectable, inject } from "@angular/core";
import { HttpClient, HttpParams } from "@angular/common/http";
import { Observable, of } from "rxjs";
import { tap, catchError } from "rxjs/operators";
import {
  Produto,
  CriarProdutoRequest,
  Movimentacao,
  AdicionarEntradaRequest,
  PagedResult,
  ListarParams,
} from "../models/models";

const CACHE_KEY = "produtos_cache";

@Injectable({ providedIn: "root" })
export class ProdutoService {
  private readonly http = inject(HttpClient);
  private readonly base = "/api/estoque/produtos";

  listarTodos(): Observable<Produto[]> {
    return this.http.get<Produto[]>(`${this.base}/`).pipe(
      tap((produtos) => this.salvarCache(produtos)),
      catchError(() => {
        const cache = this.obterCache();
        return of(cache ?? []);
      }),
    );
  }

  listarPaginado(params?: ListarParams): Observable<PagedResult<Produto>> {
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

    return this.http.get<PagedResult<Produto>>(`${this.base}/paginado`, {
      params: httpParams,
    });
  }

  obter(id: string): Observable<Produto> {
    return this.http.get<Produto>(`${this.base}/${id}`);
  }

  criar(req: CriarProdutoRequest): Observable<Produto> {
    return this.http
      .post<Produto>(this.base, req)
      .pipe(tap(() => this.invalidarCache()));
  }

  listarMovimentacoes(produtoId: string): Observable<Movimentacao[]> {
    return this.http.get<Movimentacao[]>(
      `/api/estoque/movimentacoes/produto/${produtoId}`,
    );
  }

  adicionarEntrada(req: AdicionarEntradaRequest): Observable<Produto> {
    return this.http
      .post<Produto>("/api/estoque/movimentacoes/entrada", req)
      .pipe(tap(() => this.invalidarCache()));
  }

  private salvarCache(produtos: Produto[]) {
    localStorage.setItem(CACHE_KEY, JSON.stringify(produtos));
  }

  private obterCache(): Produto[] | null {
    const raw = localStorage.getItem(CACHE_KEY);
    return raw ? JSON.parse(raw) : null;
  }

  private invalidarCache() {
    localStorage.removeItem(CACHE_KEY);
  }
}
