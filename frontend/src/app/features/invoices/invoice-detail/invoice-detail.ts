import { Component, DestroyRef, inject, OnInit, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTableModule } from '@angular/material/table';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { finalize, map, switchMap } from 'rxjs';
import { Invoice, InvoiceStatus } from '../../../core/models/invoice.model';
import { getApiErrorMessage } from '../../../core/services/api-error';
import { BillingService } from '../../../core/services/billing.service';

@Component({
  selector: 'app-invoice-detail',
  imports: [
    MatButtonModule,
    MatCardModule,
    MatChipsModule,
    MatDividerModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTableModule,
    RouterLink,
  ],
  templateUrl: './invoice-detail.html',
  styleUrl: './invoice-detail.scss',
})
export class InvoiceDetail implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly billingService = inject(BillingService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly destroyRef = inject(DestroyRef);

  readonly displayedColumns = ['product', 'quantity'];
  readonly invoice = signal<Invoice | null>(null);
  readonly loading = signal(true);
  readonly printing = signal(false);

  ngOnInit(): void {
    this.route.paramMap
      .pipe(
        map((params) => Number(params.get('id'))),
        switchMap((id) => {
          if (!Number.isInteger(id) || id <= 0) {
            throw new Error('invalid invoice id');
          }
          this.loading.set(true);
          return this.billingService.getInvoice(id).pipe(finalize(() => this.loading.set(false)));
        }),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe({
        next: (invoice) => this.invoice.set(invoice),
        error: (error) => {
          this.snackBar.open(getApiErrorMessage(error, 'Nota fiscal não encontrada.'), 'Fechar', {
            duration: 6000,
          });
          void this.router.navigate(['/invoices']);
        },
      });
  }

  print(): void {
    const invoiceToPrint = this.invoice();
    if (!invoiceToPrint || invoiceToPrint.status !== 'open' || this.printing()) {
      return;
    }

    this.printing.set(true);

    this.billingService
      .printInvoice(invoiceToPrint.id)
      .pipe(
        finalize(() => this.printing.set(false)),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe({
        next: (invoice) => {
          this.invoice.set(invoice);
          this.snackBar.open('Nota fechada e saldo do estoque atualizado.', 'Fechar', {
            duration: 4000,
          });
          requestAnimationFrame(() => window.print());
        },
        error: (error) =>
          this.snackBar.open(
            getApiErrorMessage(error, 'Não foi possível imprimir a nota fiscal.'),
            'Fechar',
            {
              duration: 7000,
            },
          ),
      });
  }

  statusLabel(status: InvoiceStatus): string {
    return status === 'open' ? 'Aberta' : 'Fechada';
  }
}
