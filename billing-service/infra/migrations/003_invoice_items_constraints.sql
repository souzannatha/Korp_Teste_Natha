BEGIN;

DO $$
BEGIN
 IF EXISTS (SELECT 1 FROM invoice_items WHERE product_code !~ '[^[:space:]]') THEN
  RAISE EXCEPTION 'Existing invoice items contain empty product codes; correct them before migrating';
 END IF;
END $$;

ALTER TABLE invoice_items
 ADD CONSTRAINT invoice_items_code_not_empty CHECK (product_code ~ '[^[:space:]]');

CREATE INDEX invoice_items_invoice_id_idx ON invoice_items (invoice_id);

COMMIT;
