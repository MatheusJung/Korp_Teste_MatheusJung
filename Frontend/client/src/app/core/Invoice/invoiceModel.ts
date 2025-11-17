export interface InvoiceItemDto {
  productCode: string;
  description:string;
  quantity: number;
}

export interface InvoiceDto {
  sequentialNumber: number;
  createdAt: string; // ISO string
  status: string;
  items: InvoiceItemDto[];
}

export interface PagedInvoiceResult {
  pageNumber: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
  items: InvoiceDto[];
}

