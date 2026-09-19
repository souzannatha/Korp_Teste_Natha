import { Component, DestroyRef, inject, OnInit, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTableModule } from '@angular/material/table';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';
import { Invoice, InvoiceStatus } from '../../../core/models/invoice.model';
import { getApiErrorMessage } from '../../../core/services/api-error';
import { BillingService } from '../../../core/services/billing.service';

@Component({
  selector: 'app-invoice-list',
  imports: [
    MatButtonModule,
    MatCardModule,
    MatChipsModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTableModule,
    RouterLink,
  ],
  templateUrl: './invoice-list.html',
  styleUrl: './invoice-list.scss',
})
export class InvoiceList implements OnInit {
  private readonly billingService = inject(BillingService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly destroyRef = inject(DestroyRef);

  readonly displayedColumns = ['number', 'status', 'items', 'actions'];
  readonly invoices = signal<Invoice[]>([]);
  readonly loading = signal(true);

  ngOnInit(): void {
    this.billingService
      .getInvoices()
      .pipe(
        finalize(() => this.loading.set(false)),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe({
        next: (invoices) => this.invoices.set(invoices),
        error: (error) =>
          this.snackBar.open(
            getApiErrorMessage(error, 'Não foi possível carregar as notas fiscais.'),
            'Fechar',
            {
              duration: 6000,
            },
          ),
      });
  }

  statusLabel(status: InvoiceStatus): string {
    return status === 'open' ? 'Aberta' : 'Fechada';
  }
}
