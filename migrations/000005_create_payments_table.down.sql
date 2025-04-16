DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP FUNCTION IF EXISTS update_payments_modified_column;
DROP TABLE IF EXISTS payments CASCADE;