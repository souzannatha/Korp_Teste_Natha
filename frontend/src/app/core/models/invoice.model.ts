export type InvoiceStatus = 'open' | 'closed';

export interface InvoiceItem {
  product_code: string;
  quantity: number;
}

export interface Invoice {
  id: number;
  number: number;
  status: InvoiceStatus;
  items: InvoiceItem[];
}

export interface CreateInvoiceRequest {
  items: InvoiceItem[];
}
