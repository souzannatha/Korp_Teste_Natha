import { convertToParamMap, ActivatedRoute, Router } from '@angular/router';
import { TestBed } from '@angular/core/testing';
import { MatSnackBar } from '@angular/material/snack-bar';
import { BehaviorSubject, of } from 'rxjs';
import { BillingService } from '../../../core/services/billing.service';
import { InvoiceDetail } from './invoice-detail';

describe('InvoiceDetail', () => {
  const openInvoice = {
    id: 7,
    number: 7,
    status: 'open' as const,
    items: [{ product_code: '001', quantity: 2 }],
  };
  const closedInvoice = { ...openInvoice, status: 'closed' as const };
  const billingService = {
    getInvoice: vi.fn(() => of(openInvoice)),
    printInvoice: vi.fn(() => of(closedInvoice)),
  };
  const snackBar = { open: vi.fn() };
  const routeParams = new BehaviorSubject(convertToParamMap({ id: '7' }));

  beforeEach(() => {
    vi.clearAllMocks();
    vi.stubGlobal('requestAnimationFrame', (callback: FrameRequestCallback) => {
      callback(0);
      return 1;
    });
    vi.spyOn(window, 'print').mockImplementation(() => undefined);

    TestBed.configureTestingModule({
      imports: [InvoiceDetail],
      providers: [
        { provide: BillingService, useValue: billingService },
        { provide: MatSnackBar, useValue: snackBar },
        { provide: Router, useValue: { navigate: vi.fn() } },
        { provide: ActivatedRoute, useValue: { paramMap: routeParams.asObservable() } },
      ],
    });
  });

  afterEach(() => vi.unstubAllGlobals());

  it('should close the invoice and open browser printing', () => {
    const fixture = TestBed.createComponent(InvoiceDetail);
    const component = fixture.componentInstance;
    component.ngOnInit();

    expect(component.invoice()?.status).toBe('open');
    expect(component.loading()).toBe(false);
    component.print();

    expect(billingService.printInvoice).toHaveBeenCalledWith(7);
    expect(component.invoice()?.status).toBe('closed');
    expect(window.print).toHaveBeenCalled();
  });

  it('should not print a closed invoice again', () => {
    const fixture = TestBed.createComponent(InvoiceDetail);
    const component = fixture.componentInstance;
    component.invoice.set(closedInvoice);

    component.print();

    expect(billingService.printInvoice).not.toHaveBeenCalled();
  });
});
