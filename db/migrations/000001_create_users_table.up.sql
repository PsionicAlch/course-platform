CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,                                -- ULID as the primary key (stored as TEXT)
    email TEXT NOT NULL,                                -- User email (unique and required)
    password TEXT NOT NULL,                             -- Hashed password (required)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,      -- Creation timestamp
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP       -- Update timestamp
);

-- Create a unique index on the email column to allow quick lookups by email and ensure uniqueness.
CREATE UNIQUE INDEX idx_users_email ON users(email);
