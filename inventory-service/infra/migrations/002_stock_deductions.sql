BEGIN;

DO $$
BEGIN
 IF EXISTS (SELECT 1 FROM products WHERE code !~ '[^[:space:]]' OR description !~ '[^[:space:]]') THEN
  RAISE EXCEPTION 'Existing products contain empty code or description; correct them before migrating';
 END IF;
END $$;

ALTER TABLE products
 ADD CONSTRAINT products_code_not_empty CHECK (code ~ '[^[:space:]]'),
 ADD CONSTRAINT products_description_not_empty CHECK (description ~ '[^[:space:]]');

CREATE TABLE stock_deductions (
 invoice_id BIGINT PRIMARY KEY CHECK (invoice_id > 0),
 fingerprint VARCHAR(64) NOT NULL CHECK (fingerprint ~ '^[0-9a-f]{64}$'),
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
