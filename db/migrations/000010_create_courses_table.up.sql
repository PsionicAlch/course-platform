-- Courses table holds the content and meta data for each individual course to remove
-- the need for a RAM based cache of all the compiled courses (viable option though).
CREATE TABLE IF NOT EXISTS courses (
    id TEXT PRIMARY KEY,

    title TEXT NOT NULL,                                                                    -- The title of the course.
    slug TEXT NOT NULL,                                                                     -- URL friendly slug for the course.
    description TEXT NOT NULL,                                                              -- A short description of the course.
    thumbnail_url TEXT NOT NULL,                                                            -- URL for the thumbnail image of the course.
    banner_url TEXT NOT NULL,                                                               -- URL for the banner image of the course.
    content TEXT NOT NULL,                                                                  -- HTML based contents of the course.
    published INTEGER DEFAULT 0 CHECK (published >= 0 AND published <= 1),                  -- BOOLEAN to represent whether or not the course has been published.

    author_id TEXT DEFAULT NULL,                                                            -- The user ID who published the course.

    file_checksum TEXT NOT NULL,                                                            -- A SHA256 checksum to speed up the process of checking if a file has changed.
    file_key TEXT NOT NULL,                                                                 -- Unique file key.

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                          -- Timestamp for when the course was initially created.
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                          -- Timestamp for when the course was updated.

    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL
);

-- A unique index to help speed up searching via slugs and/or file_key.
CREATE UNIQUE INDEX IF NOT EXISTS idx_courses_slug ON courses(slug);

CREATE UNIQUE INDEX IF NOT EXISTS idx_courses_file_key ON courses(file_key);

CREATE TRIGGER IF NOT EXISTS trigger_update_courses_updated_at
AFTER UPDATE ON courses
FOR EACH ROW
BEGIN
    UPDATE courses SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
