CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,                                -- ULID as the primary key (stored as TEXT)
    name TEXT NOT NULL,                                 -- User's name
    surname TEXT NOT NULL,                              -- User's surname
    email TEXT NOT NULL,                                -- User email (unique and required)
    password TEXT NOT NULL,                             -- Hashed password (required)
    is_admin INTEGER DEFAULT 0,                            -- Boolean for whether or not the user is an administrator
    is_author INTEGER DEFAULT 0,                           -- Boolean for whether or not the user is an author

    affiliate_code TEXT NOT NULL,                       -- User's affiliate code for discounts

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,      -- Creation timestamp
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP       -- Update timestamp
);

-- Create a unique index on the email column to allow quick lookups by email and ensure uniqueness.
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- Create a unique index on the affiliate code column to ensure that affiliate codes are unique
-- and to speed up searches based on the user's affiliate code.
CREATE UNIQUE INDEX idx_users_affiliate_code ON users(affiliate_code);
