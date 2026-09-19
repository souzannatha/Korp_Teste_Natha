import { Routes } from '@angular/router';
import { AppShell } from './layout/app-shell/app-shell';

export const routes: Routes = [
  {
    path: '',
    component: AppShell,
    children: [
      {
        path: 'products',
        loadComponent: () =>
          import('./features/products/product-list/product-list').then(
            (component) => component.ProductList,
          ),
      },
      {
        path: 'products/new',
        loadComponent: () =>
          import('./features/products/product-form/product-form').then(
            (component) => component.ProductForm,
          ),
      },
      {
        path: 'invoices',
        loadComponent: () =>
          import('./features/invoices/invoice-list/invoice-list').then(
            (component) => component.InvoiceList,
          ),
      },
      {
        path: 'invoices/new',
        loadComponent: () =>
          import('./features/invoices/invoice-form/invoice-form').then(
            (component) => component.InvoiceForm,
          ),
      },
      {
        path: 'invoices/:id',
        loadComponent: () =>
          import('./features/invoices/invoice-detail/invoice-detail').then(
            (component) => component.InvoiceDetail,
          ),
      },
      { path: '', pathMatch: 'full', redirectTo: 'invoices' },
      { path: '**', redirectTo: 'invoices' },
    ],
  },
];
