DROP TRIGGER IF EXISTS trigger_update_course_chapters_updated_at;

DROP INDEX IF EXISTS idx_course_chapters_course_id_slug;

DROP INDEX IF EXISTS idx_course_chapters_course_id_title;

DROP INDEX IF EXISTS idx_courses_chapters_course_id_chapter;

DROP TABLE IF EXISTS course_chapters;
