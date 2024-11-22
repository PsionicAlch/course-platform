-- Courses Chapters is a table to hold the individual chapters of each course along with
-- all their meta data and content.
CREATE TABLE IF NOT EXISTS courses_chapters (
    id TEXT PRIMARY KEY,

    title TEXT NOT NULL,
    chapter INTEGER NOT NULL,
    content TEXT NOT NULL,

    course_id TEXT NOT NULL,

    file_checksum TEXT NOT NULL,
    file_key TEXT NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_courses_chapters_course_id_chapter ON courses(course_id, chapter);
