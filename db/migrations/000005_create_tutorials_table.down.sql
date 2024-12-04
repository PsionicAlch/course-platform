DROP TRIGGER IF EXISTS trigger_update_tutorials_updated_at;

DROP INDEX IF EXISTS idx_tutorials_slug;

DROP INDEX IF EXISTS idx_tutorials_file_key;

DROP TABLE IF EXISTS tutorials;
