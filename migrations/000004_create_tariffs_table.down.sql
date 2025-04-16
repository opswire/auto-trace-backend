DROP TRIGGER IF EXISTS update_tariffs_updated_at ON tariffs;
DROP FUNCTION IF EXISTS update_tariffs_modified_column;
DROP TABLE IF EXISTS tariffs CASCADE;