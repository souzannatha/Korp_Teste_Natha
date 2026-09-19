import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { InventoryService } from './inventory.service';

describe('InventoryService', () => {
  let service: InventoryService;
  let httpTesting: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [InventoryService, provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(InventoryService);
    httpTesting = TestBed.inject(HttpTestingController);
  });

  afterEach(() => httpTesting.verify());

  it('should list products from Inventory', () => {
    service.getProducts().subscribe((products) => expect(products[0].code).toBe('001'));

    const request = httpTesting.expectOne('http://localhost:8080/product');
    expect(request.request.method).toBe('GET');
    request.flush([{ id: 1, code: '001', description: 'Teclado', balance: 10 }]);
  });

  it('should send the product fields when creating a product', () => {
    const product = { code: '001', description: 'Teclado', balance: 10 };

    service.createProduct(product).subscribe();

    const request = httpTesting.expectOne('http://localhost:8080/product');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual(product);
    request.flush({ id: 1, ...product });
  });
});
