# Front-end Angular

Interface do sistema de produtos e notas fiscais. O projeto usa Angular, Angular Material, Reactive Forms, Router, HttpClient e RxJS.

## Executar

Antes de iniciar o front-end, deixe o Inventory disponível em `http://localhost:8080` e o Billing em `http://localhost:8081`.

```bash
cd frontend
npm install
npm start
```

Acesse `http://localhost:4200`.

Não é necessário instalar o Angular CLI globalmente. Os scripts do npm usam a versão instalada no próprio projeto.

## Rotas

- `/products`: lista produtos e saldos.
- `/products/new`: cadastra um produto.
- `/invoices`: lista notas fiscais.
- `/invoices/new`: cria uma nota com um ou mais produtos.
- `/invoices/:id`: mostra os detalhes e permite imprimir uma nota aberta.

## Verificação

```bash
npm run build
npm test -- --watch=false
```

As URLs das APIs ficam em `src/environments/environment.ts`.

O cache persistente do Angular está desabilitado em `angular.json` porque a versão nativa do LMDB apresentou falha no macOS usado no desenvolvimento. Isso não altera o funcionamento da aplicação; apenas evita o cache entre builds.
