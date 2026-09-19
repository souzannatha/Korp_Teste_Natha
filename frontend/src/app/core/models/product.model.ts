export interface Product {
  id: number;
  code: string;
  description: string;
  balance: number;
}

export interface CreateProductRequest {
  code: string;
  description: string;
  balance: number;
}
