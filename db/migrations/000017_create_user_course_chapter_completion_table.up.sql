CREATE TABLE IF NOT EXISTS user_course_chapter_completion (
    id TEXT PRIMARY KEY,                                         -- ULID as the primary key

    user_id TEXT NOT NULL,                                       -- References the `users` table
    course_id TEXT NOT NULL,                                     -- References the `courses` table
    chapter_id TEXT NOT NULL,                                    -- References the `course_chapters` table

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,               -- Timestamp for when the chapter was marked as completed

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (chapter_id) REFERENCES course_chapters(id) ON DELETE CASCADE
);

-- Ensure a user can only have one completion entry per chapter in a course
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_course_chapter_completion ON user_course_chapter_completion(user_id, course_id, chapter_id);
