DROP TRIGGER IF EXISTS trigger_update_whitelisted_ips_updated_at;

DROP INDEX IF EXISTS idx_whitelisted_ips_user_ip;

DROP TABLE IF EXISTS whitelisted_ips;
