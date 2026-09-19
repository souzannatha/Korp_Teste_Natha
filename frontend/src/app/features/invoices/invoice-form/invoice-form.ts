import { Component, DestroyRef, inject, OnInit, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormArray, FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatDividerModule } from '@angular/material/divider';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { Router, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';
import { Product } from '../../../core/models/product.model';
import { getApiErrorMessage } from '../../../core/services/api-error';
import { BillingService } from '../../../core/services/billing.service';
import { InventoryService } from '../../../core/services/inventory.service';

type InvoiceItemForm = FormGroup<{
  product_code: FormControl<string>;
  quantity: FormControl<number>;
}>;

@Component({
  selector: 'app-invoice-form',
  imports: [
    MatButtonModule,
    MatCardModule,
    MatDividerModule,
    MatFormFieldModule,
    MatInputModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    MatSnackBarModule,
    MatTooltipModule,
    ReactiveFormsModule,
    RouterLink,
  ],
  templateUrl: './invoice-form.html',
  styleUrl: './invoice-form.scss',
})
export class InvoiceForm implements OnInit {
  private readonly inventoryService = inject(InventoryService);
  private readonly billingService = inject(BillingService);
  private readonly router = inject(Router);
  private readonly snackBar = inject(MatSnackBar);
  private readonly destroyRef = inject(DestroyRef);

  readonly products = signal<Product[]>([]);
  readonly loadingProducts = signal(true);
  readonly saving = signal(false);

  readonly invoiceForm = new FormGroup({
    items: new FormArray<InvoiceItemForm>([]),
  });

  get items(): FormArray<InvoiceItemForm> {
    return this.invoiceForm.controls.items;
  }

  ngOnInit(): void {
    this.addItem();

    this.inventoryService
      .getProducts()
      .pipe(
        finalize(() => this.loadingProducts.set(false)),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe({
        next: (products) => this.products.set(products),
        error: (error) =>
          this.snackBar.open(
            getApiErrorMessage(error, 'Não foi possível carregar os produtos.'),
            'Fechar',
            {
              duration: 6000,
            },
          ),
      });
  }

  addItem(): void {
    this.items.push(
      new FormGroup({
        product_code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
        quantity: new FormControl(1, {
          nonNullable: true,
          validators: [Validators.required, Validators.min(1), Validators.pattern(/^\d+$/)],
        }),
      }),
    );
  }

  removeItem(index: number): void {
    if (this.items.length > 1) {
      this.items.removeAt(index);
    }
  }

  productBalance(code: string): number | null {
    return this.products().find((product) => product.code === code)?.balance ?? null;
  }

  save(): void {
    if (this.invoiceForm.invalid || this.saving()) {
      this.invoiceForm.markAllAsTouched();
      return;
    }

    this.saving.set(true);

    this.billingService
      .createInvoice(this.invoiceForm.getRawValue())
      .pipe(
        finalize(() => this.saving.set(false)),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe({
        next: (invoice) => {
          this.snackBar.open(`Nota ${invoice.number} criada com status Aberta.`, 'Fechar', {
            duration: 4000,
          });
          void this.router.navigate(['/invoices', invoice.id]);
        },
        error: (error) =>
          this.snackBar.open(
            getApiErrorMessage(error, 'Não foi possível criar a nota fiscal.'),
            'Fechar',
            {
              duration: 6000,
            },
          ),
      });
  }
}
