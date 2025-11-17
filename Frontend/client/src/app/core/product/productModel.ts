export interface ProductDto {
  productCode: string;
  description: string;
  quantity: number;
  status: string;
}

export interface AdjustStockDto {
  productCode: string;
  quantity: number;
  operationKey:string;
}

export interface PagedProductResult {
  pageNumber: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
  items: ProductDto[];
}

