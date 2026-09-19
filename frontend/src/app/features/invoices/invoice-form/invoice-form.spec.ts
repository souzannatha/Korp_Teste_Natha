import { TestBed } from '@angular/core/testing';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { of } from 'rxjs';
import { BillingService } from '../../../core/services/billing.service';
import { InventoryService } from '../../../core/services/inventory.service';
import { InvoiceForm } from './invoice-form';

describe('InvoiceForm', () => {
  const inventoryService = {
    getProducts: vi.fn(() => of([{ id: 1, code: '001', description: 'Teclado', balance: 10 }])),
  };
  const billingService = {
    createInvoice: vi.fn(() =>
      of({
        id: 5,
        number: 5,
        status: 'open' as const,
        items: [{ product_code: '001', quantity: 3 }],
      }),
    ),
  };
  const router = { navigate: vi.fn(() => Promise.resolve(true)) };
  const snackBar = { open: vi.fn() };

  beforeEach(() => {
    vi.clearAllMocks();
    TestBed.configureTestingModule({
      imports: [InvoiceForm],
      providers: [
        { provide: InventoryService, useValue: inventoryService },
        { provide: BillingService, useValue: billingService },
        { provide: Router, useValue: router },
        { provide: MatSnackBar, useValue: snackBar },
      ],
    });
  });

  it('should add and remove invoice items', () => {
    const fixture = TestBed.createComponent(InvoiceForm);
    const component = fixture.componentInstance;
    component.ngOnInit();

    expect(component.items.length).toBe(1);
    component.addItem();
    expect(component.items.length).toBe(2);
    component.removeItem(1);
    expect(component.items.length).toBe(1);
  });

  it('should send the selected product and quantity', () => {
    const fixture = TestBed.createComponent(InvoiceForm);
    const component = fixture.componentInstance;
    component.ngOnInit();
    component.items.at(0).setValue({ product_code: '001', quantity: 3 });

    component.save();

    expect(billingService.createInvoice).toHaveBeenCalledWith({
      items: [{ product_code: '001', quantity: 3 }],
    });
    expect(router.navigate).toHaveBeenCalledWith(['/invoices', 5]);
  });
});
