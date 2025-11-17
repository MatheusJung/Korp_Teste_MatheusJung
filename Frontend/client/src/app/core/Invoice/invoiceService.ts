import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { InvoiceDto, InvoiceItemDto, PagedInvoiceResult} from '../Invoice/invoiceModel';
import { Observable } from 'rxjs';
import { catchError, map, timeout, retry, of } from 'rxjs';
import { environment } from '../enviroment';

@Injectable({
  providedIn: 'root',
})export class InvoiceService {
  private apiUrl = environment.billingServiceUrl; // ajuste se necessário

  constructor(private http: HttpClient) {}

  //Verifica se o serviço esta online
  getInvoiceServiceStatus(){
    this.http.get(`${this.apiUrl}/health`).pipe(
      timeout(3000),
      retry(2),
      map(() => true),
      catchError(() => of(false))
    );
  }

  // Listar notas paginadas
  getInvoicesPaged(pageNumber = 1, pageSize = 10): Observable<PagedInvoiceResult> {
    const params = new HttpParams()
      .set('pageNumber', pageNumber.toString())
      .set('pageSize', pageSize.toString());
    return this.http.get<PagedInvoiceResult>(`${this.apiUrl}/invoices`, { params });
  }

  // Consultar nota por sequencial
  getInvoice(seqNumber: number): Observable<InvoiceDto> {
    return this.http.get<InvoiceDto>(`${this.apiUrl}/invoices/${seqNumber}`);
  }

  // Criar nota
  createInvoice(items: InvoiceItemDto[]): Observable<InvoiceDto> {
    return this.http.post<InvoiceDto>(`${this.apiUrl}/invoices`, { items });
  }

  // Adicionar item à nota
  addItem(seqNumber: number, item: InvoiceItemDto): Observable<InvoiceDto> {
    return this.http.put<InvoiceDto>(`${this.apiUrl}/invoices/${seqNumber}/items`, item);
  }

  // Cancelar nota
  cancelInvoice(seqNumber: number): Observable<InvoiceDto> {
    return this.http.post<InvoiceDto>(`${this.apiUrl}/invoices/${seqNumber}/cancel`, null);
  }

  // Fechar nota
  closeInvoice(seqNumber: number): Observable<InvoiceDto> {
    return this.http.post<InvoiceDto>(`${this.apiUrl}/invoices/${seqNumber}/close`, null);
  }
  //Imprimir e Fechar nota
   closeAndPrint(seqNumber: number) {
    return this.http.post(
      `${this.apiUrl}/invoices/${seqNumber}/close-and-print`,
      {}, // body vazio
      { responseType: 'blob' } // importante para receber PDF
    );
  }
}

