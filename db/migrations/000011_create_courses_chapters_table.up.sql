-- Courses Chapters is a table to hold the individual chapters of each course along with
-- all their meta data and content.
CREATE TABLE IF NOT EXISTS course_chapters (
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

CREATE UNIQUE INDEX IF NOT EXISTS idx_courses_chapters_course_id_chapter ON course_chapters(course_id, chapter);

CREATE TRIGGER IF NOT EXISTS trigger_update_course_chapters_updated_at
AFTER UPDATE ON course_chapters
FOR EACH ROW
BEGIN
    UPDATE course_chapters SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
