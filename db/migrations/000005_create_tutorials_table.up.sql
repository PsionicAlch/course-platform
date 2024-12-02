-- Tutorials table holds the content and meta data for each individual tutorial to remove
-- the need for a RAM based cache of all the compiled tutorials (viable option though).
CREATE TABLE IF NOT EXISTS tutorials (
    id TEXT PRIMARY KEY,                                                                        -- The ID of the individual tutorial.

    title TEXT NOT NULL,                                                                        -- The title of the tutorial.
    slug TEXT NOT NULL,                                                                         -- URL friendly slug for the tutorial.
    description TEXT NOT NULL,                                                                  -- A short description of the tutorial.
    thumbnail_url TEXT NOT NULL,                                                                -- URL for the thumbnail image of the tutorial.
    banner_url TEXT NOT NULL,                                                                   -- URL for the banner image of the tutorial.
    content TEXT NOT NULL,                                                                      -- HTML based contents of the tutorial.
    published INTEGER DEFAULT 0 CHECK (published >= 0 AND published <= 1),                      -- BOOLEAN to represent whether or not the tutorial has been published.

    author_id TEXT DEFAULT NULL,                                                                -- The user ID who published the tutorial (for when I have multiple authors).

    file_checksum TEXT NOT NULL,                                                                -- A SHA256 checksum to speed up the process of checking if a file has changed.
    file_key TEXT NOT NULL,                                                                     -- Unique file key.

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                              -- Timestamp for when the tutorial was initially created.
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                              -- Timestamp for when the tutorial was updated.

    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL                             -- Foreign key for users table. Set to null if the author is deleted.
);

-- A unique index to help speed up searching via slugs and/or file_key.
CREATE UNIQUE INDEX idx_tutorials_slug ON tutorials(slug);
CREATE UNIQUE INDEX idx_tutorials_file_key ON tutorials(file_key);
