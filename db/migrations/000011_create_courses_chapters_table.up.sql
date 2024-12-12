-- Courses Chapters is a table to hold the individual chapters of each course along with
-- all their meta data and content.
CREATE TABLE IF NOT EXISTS course_chapters (
    id TEXT PRIMARY KEY,

    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    chapter INTEGER NOT NULL,
    content TEXT NOT NULL,

    course_id TEXT NOT NULL,

    file_checksum TEXT NOT NULL,
    file_key TEXT NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Ensure that you can't have multiple entries with the same chapter number.
CREATE UNIQUE INDEX IF NOT EXISTS idx_course_chapters_course_id_chapter ON course_chapters(course_id, chapter);

-- Ensure that all titles are unique per course. It's okay if multiple courses have a chapter with the same title,
-- we just don't want the same course to have multiple chapters with the same title.
CREATE UNIQUE INDEX IF NOT EXISTS idx_course_chapters_course_id_title ON course_chapters(course_id, title);

-- Ensure that all slugs are unique per course. It's okay if multiple courses have a chapter with the same slug,
-- we just don't want the same course to have multiple chapters with the same slug.
CREATE UNIQUE INDEX IF NOT EXISTS idx_course_chapters_course_id_slug ON course_chapters(course_id, slug);

CREATE TRIGGER IF NOT EXISTS trigger_update_course_chapters_updated_at
AFTER UPDATE ON course_chapters
FOR EACH ROW
BEGIN
    UPDATE course_chapters SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
