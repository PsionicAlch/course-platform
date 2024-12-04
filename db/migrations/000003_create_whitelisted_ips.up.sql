CREATE TABLE IF NOT EXISTS whitelisted_ips (
    id TEXT PRIMARY KEY,

    user_id TEXT NOT NULL,
    ip_address TEXT NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_whitelisted_ips_user_ip ON whitelisted_ips(user_id, ip_address);

CREATE TRIGGER IF NOT EXISTS trigger_update_whitelisted_ips_updated_at
AFTER UPDATE ON whitelisted_ips
FOR EACH ROW
BEGIN
    UPDATE whitelisted_ips SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
