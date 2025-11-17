import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ProductDto,AdjustStockDto} from './productModel';
import { catchError, map, timeout, retry, of } from 'rxjs';
import { environment } from '../enviroment';

@Injectable({
  providedIn: 'root'
})
export class ProductService {
  private apiUrl = environment.inventoryServiceUrl;

  constructor(private http: HttpClient) {}

  //Verifica se o serviço esta online
  getInventoryServiceStatus(){
    this.http.get(`${this.apiUrl}/health`).pipe(
      timeout(3000),
      retry(2),
      map(() => true),
      catchError(() => of(false))
    );
  }

  // Cadastrar produto
  createProduct(dto: { code: string; name: string; initialStock: number }): Observable<ProductDto> {
    return this.http.post<ProductDto>(`${this.apiUrl}/products`, dto);
  }

  // Listar todos com paginação
  getProducts(pageNumber = 1, pageSize = 10): Observable<any> {
    const params = new HttpParams()
      .set('pageNumber', pageNumber)
      .set('pageSize', pageSize);
    return this.http.get<any>(`${this.apiUrl}/products`, { params });
  }

  // Listar produtos ativos
  getActiveProducts(): Observable<ProductDto[]> {
    return this.http.get<ProductDto[]>(`${this.apiUrl}/products/ActiveProducts`);
  }

  // Consultar por código
  getProductByCode(code: string): Observable<ProductDto> {
    return this.http.get<ProductDto>(`${this.apiUrl}/products/${code}`);
  }

  // Desativar produto
  deactivateProduct(code: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/products/${code}`);
  }

  // Adicionar estoque
  addStock(dto: AdjustStockDto): Observable<any> {
    return this.http.post(`${this.apiUrl}/stock/add`, dto);
  }

  // Remover estoque
  removeStock(dto: AdjustStockDto): Observable<any> {
    return this.http.post(`${this.apiUrl}/stock/remove`, dto);
  }
}

