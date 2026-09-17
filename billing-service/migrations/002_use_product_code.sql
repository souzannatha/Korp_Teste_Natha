-- Converts the column structure only. Existing product IDs are not mapped
-- to Inventory codes: ID 1 becomes '1', not '001'.
BEGIN;

ALTER TABLE invoice_items
    RENAME COLUMN product_id TO product_code;

ALTER TABLE invoice_items
    ALTER COLUMN product_code TYPE VARCHAR(50)
    USING product_code::VARCHAR(50);

COMMIT;
