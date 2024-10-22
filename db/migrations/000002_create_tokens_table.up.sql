CREATE TABLE IF NOT EXISTS tokens (
    id TEXT PRIMARY KEY,                            -- ULID as the primary key
    token TEXT NOT NULL,                            -- Actual token (required)
    token_type TEXT NOT NULL,                       -- Type of token (e.g., 'authentication', 'password_reset')
    valid_until DATETIME NOT NULL,                  -- Expiration date of the token
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Timestamp when token was created
    user_id TEXT NOT NULL,                          -- Foreign key referencing users table
    ip_address TEXT NOT NULL,                       -- The IP address of the client that requested this token

    -- Foreign key constraint: user_id must exist in users table
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create a composite index on the token and token_type columns for faster lookups
CREATE UNIQUE INDEX idx_tokens_token_type ON tokens(token, token_type);
