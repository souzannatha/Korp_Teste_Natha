CREATE TABLE invoices (
  id SERIAL PRIMARY KEY,
  number INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
  status VARCHAR(10) NOT NULL DEFAULT 'open',

  CONSTRAINT invoice_status
    CHECK (status IN ('open', 'closed'))
);

CREATE TABLE invoice_items (
  id SERIAL PRIMARY KEY,
  invoice_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,

  CONSTRAINT invoice_items_invoice_fk
    FOREIGN KEY (invoice_id) REFERENCES invoices(id),

  CONSTRAINT invoice_items_quantity
    CHECK (quantity > 0)
);