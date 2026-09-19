# Korp_Teste_Natha

Sistema de produtos e notas fiscais com front-end Angular e dois microsserviços em Go, Gin e PostgreSQL.

- Inventory: produtos e baixa de estoque, API na porta 8080.
- Billing: notas e finalização da impressão, API na porta 8081.
- Angular: interface de produtos e notas fiscais, executada na porta 4200.

## Organização

Cada serviço mantém `cmd/main.go`, `docs` e os packages em `internal`: `controller`, `model`, `repository`, `routes` e `usecases`. A pasta `infra` fica na raiz do serviço, separada de `internal`, contendo somente `db` e `migrations`. O Billing também possui `internal/client` para HTTP.

O fluxo permanece controller → use case → repository. Cada serviço tem seu próprio banco e pool; o Billing não acessa o banco do Inventory.

## Executar

Requisitos: Go 1.27.1 e Docker Compose. Execute os comandos a partir da raiz de cada serviço.

No Inventory:

```bash
cd inventory-service
cp .env.example .env
```

Preencha `DB_PASSWORD` no arquivo local e execute:

```bash
docker compose up -d
go run ./cmd
```

Em outro terminal, faça o mesmo no Billing:

```bash
cd billing-service
cp .env.example .env
```

Preencha `DB_PASSWORD` e execute:

```bash
docker compose up -d
go run ./cmd
```

Em um terceiro terminal, inicie o front-end:

```bash
cd frontend
npm install
npm start
```

Acesse http://localhost:4200. Não é necessário instalar o Angular CLI globalmente.

Copie o exemplo somente se ainda não existir `.env`; não sobrescreva suas credenciais. O exemplo do Billing usa banco `billing` e porta 5433, enquanto o Inventory usa `inventory` e porta 5432. Os containers usam internamente 5432.

Arquivos `.env` não são versionados. Também é possível fornecer todas as variáveis pelo ambiente, sem arquivo `.env`. Não altere `DB_USER`, `DB_PASSWORD` ou `DB_NAME` de um banco existente esperando que o Docker atualize as credenciais.

## Configuração

Variáveis obrigatórias: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD` e `DB_NAME`.

Variáveis opcionais e padrões:

- `DB_SSLMODE=disable`: desenvolvimento local; configure SSL conforme o banco do ambiente.
- `DB_MAX_OPEN_CONNS=10`.
- `DB_MAX_IDLE_CONNS=5`.
- `DB_CONN_MAX_LIFETIME=30m`.
- `DB_CONN_MAX_IDLE_TIME=5m`.
- `CORS_ORIGIN=http://localhost:4200`.
- `INVENTORY_URL=http://localhost:8080`: somente Billing.

Um único `*sql.DB` é criado por serviço e reutilizado pelos repositories. Operações SQL possuem timeout de 5 segundos. O client HTTP também possui timeout de 5 segundos. A transação de impressão é limitada a 15 segundos; requisições têm limite de 25 segundos e corpo máximo de 1 MiB.

## Rotas e Swagger

Inventory:

- `POST /product`
- `GET /product`
- `GET /product/:codeProduct`
- `POST /stock-deductions`
- `GET /ping`

Billing:

- `POST /invoice`
- `GET /invoice`
- `GET /invoice/:id`
- `POST /invoice/:id/print`
- `GET /healthy`

Swagger:

- Inventory: http://localhost:8080/swagger/index.html
- Billing: http://localhost:8081/swagger/index.html

Cada serviço tem `docs/openapi.yaml` separado do código. Os assets locais são do `swagger-ui-dist` 5.33.0, com licença e aviso distribuídos em `docs/assets`. A UI não precisa de CDN.

`/ping` e `/healthy` verificam que o processo responde; não são verificações de disponibilidade do banco.

## Fluxo de exemplo

Cadastre um produto no Inventory:

```bash
curl -i http://localhost:8080/product \
  -H 'Content-Type: application/json' \
  -d '{"code":"001","description":"Teclado mecânico","balance":10}'
```

Saldo é obrigatório, mas zero é válido. Código é uma string única e preserva zeros à esquerda.

Crie uma nota no Billing:

```bash
curl -i http://localhost:8081/invoice \
  -H 'Content-Type: application/json' \
  -d '{"items":[{"product_code":"001","quantity":3}]}'
```

O Billing consulta o produto por código antes de salvar. A nota nasce `open`, com ID e número gerados pelo banco. O saldo continua 10.

Substitua 1 pelo ID retornado e finalize:

```bash
curl -i -X POST http://localhost:8081/invoice/1/print
curl -i http://localhost:8080/product/001
```

A nota fica `closed` e o saldo fica 7. Essa rota finaliza a operação de negócio; ela não gera PDF nem envia trabalho à impressora.

## Falhas, concorrência e idempotência

Na impressão, o Billing bloqueia a nota, envia seus itens ao Inventory e fecha somente após confirmação.

O Inventory agrupa produtos repetidos, ordena os códigos e executa toda a baixa em uma transação. Saldo insuficiente ou produto ausente desfazem todos os descontos.

A mesma transação grava `invoice_id` e um hash dos itens normalizados. Repetir a operação não desconta novamente; conteúdo diferente para a mesma nota retorna 409.

Dois pedidos simultâneos da mesma nota são serializados pelo Billing. Notas diferentes competem pelo estoque com atualização condicional, impedindo saldo negativo.

Para demonstrar indisponibilidade, interrompa somente o processo do Inventory, mantendo seus dados. Tente imprimir uma nota aberta: o Billing retorna 503 e mantém a nota aberta. Reinicie o Inventory e tente novamente.

Se o Inventory já tiver confirmado a baixa, mas a resposta se perder ou o Billing falhar ao fechar, repita a impressão da mesma nota. A baixa existente é reconhecida e o fechamento pode ser concluído. Não crie outra nota para repetir a mesma operação.

Respostas de erro usam `{"message":"..."}`:

- 400: entrada inválida.
- 404: nota ou produto ausente.
- 409: código duplicado, nota fechada, saldo insuficiente ou baixa conflitante.
- 500: falha interna.
- 502: resposta inválida do Inventory.
- 503: Inventory indisponível.
- 504: timeout do Inventory; a baixa pode ter sido confirmada, então repita a mesma nota.

Não existem retries automáticos, filas ou reconciliação em segundo plano. A recuperação é feita por nova tentativa do usuário. A deduplicação assume uma única instância lógica do Billing com IDs persistidos; não reutilize esses IDs em outro banco de Billing apontando para o mesmo Inventory.

## Migrations

As migrations ficam em `infra/migrations`. O PostgreSQL executa os arquivos em ordem somente na primeira inicialização de um volume vazio. Reiniciar containers não aplica novas migrations.

Para um clone com bancos novos, `docker compose up -d` executa todas as migrations automaticamente. Não é necessário apagar volumes.

Para bancos existentes, aplique somente os arquivos ainda não executados, uma vez, em ordem. Faça backup antes de alterações de schema.

No Inventory, se a tabela de produtos já existe e ainda não há tabela `stock_deductions`:

```bash
docker compose exec -T inventory_db \
  sh -c 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' \
  < infra/migrations/002_stock_deductions.sql
```

No Billing, primeiro confira se `invoice_items` ainda possui `product_id`. Somente nesse caso execute:

```bash
docker compose exec -T billing_db \
  sh -c 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' \
  < infra/migrations/002_use_product_code.sql
```

A migration 002 converte ID 1 para código `"1"`, não `"001"`. Dados legados precisam ser conferidos com os códigos reais do Inventory; não há conversão automática entre os bancos.

Depois, se a constraint `invoice_items_code_not_empty` ainda não existe:

```bash
docker compose exec -T billing_db \
  sh -c 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' \
  < infra/migrations/003_invoice_items_constraints.sql
```

As novas migrations rejeitam dados com códigos ou descrições vazios antes de adicionar restrições. Nesse caso, revise e corrija esses registros e repita a migration; ela não apaga dados automaticamente.

Não há runner com histórico de versões nesta etapa. Não execute novamente os scripts sobre um schema já migrado.

## Verificação

Na raiz de cada serviço:

```bash
go build ./...
go vet ./...
```

Numeração é crescente e pode ter lacunas em caso de falhas. Não reutilizar IDs nem reiniciar sequences para eliminar lacunas.

## Detalhamento técnico

### Angular

O ciclo de vida utilizado foi `ngOnInit`, responsável pelos carregamentos
iniciais das listagens, dos produtos do formulário e dos detalhes da nota.

Não foi necessário implementar `ngOnDestroy` manualmente. As inscrições são
encerradas com `takeUntilDestroyed` e `DestroyRef`.

O RxJS foi utilizado nos Observables do HttpClient. Os principais operadores
utilizados foram:

- `map`: conversão do parâmetro da rota;
- `switchMap`: consulta da nota a partir do ID da rota;
- `finalize`: finalização dos indicadores de carregamento;
- `takeUntilDestroyed`: encerramento automático das inscrições.

Também foram utilizados Angular Router, Reactive Forms e HttpClient.

A interface utiliza Angular Material/CDK, incluindo toolbar, cards, tabelas,
inputs, selects, botões, chips, spinner, snackbar, tooltip e divider.

### Golang

As dependências são gerenciadas por Go Modules, por meio dos arquivos
`go.mod` e `go.sum`.

Os principais pacotes utilizados são:

- Gin: criação das APIs HTTP;
- lib/pq: driver de conexão com PostgreSQL;
- godotenv: carregamento opcional das variáveis do arquivo `.env`;
- database/sql: acesso ao banco e gerenciamento do pool de conexões.

Os erros são validados nas camadas controller e use case e convertidos em
respostas HTTP apropriadas. Operações críticas utilizam transações, rollback,
contexts e timeouts.

O projeto não utiliza C#. Portanto, LINQ não se aplica a esta implementação.
