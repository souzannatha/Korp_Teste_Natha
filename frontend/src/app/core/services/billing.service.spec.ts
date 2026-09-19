import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { BillingService } from './billing.service';

describe('BillingService', () => {
  let service: BillingService;
  let httpTesting: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [BillingService, provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(BillingService);
    httpTesting = TestBed.inject(HttpTestingController);
  });

  afterEach(() => httpTesting.verify());

  it('should create an invoice using only its items', () => {
    const invoice = { items: [{ product_code: '001', quantity: 3 }] };

    service.createInvoice(invoice).subscribe();

    const request = httpTesting.expectOne('http://localhost:8081/invoice');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual(invoice);
    request.flush({ id: 1, number: 1, status: 'open', ...invoice });
  });

  it('should print an invoice without a request body', () => {
    service.printInvoice(7).subscribe();

    const request = httpTesting.expectOne('http://localhost:8081/invoice/7/print');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toBeNull();
    request.flush({ id: 7, number: 7, status: 'closed', items: [] });
  });
});
