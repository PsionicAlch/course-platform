-- Courses Keywords is a pivot table since a course can have multiple keywords whilst
-- a keyword can belong to multiple tutorials.
CREATE TABLE IF NOT EXISTS courses_keywords (
    id TEXT PRIMARY KEY,                                                                        -- The ID for each courses keywords pair.

    course_id TEXT NOT NULL,                                                                    -- A reference to which course uses this keyword.
    keyword_id TEXT NOT NULL,                                                                   -- A reference to the keyword.

    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,                           -- If the course gets deleted so should this row.
    FOREIGN KEY (keyword_id) REFERENCES keywords(id) ON DELETE CASCADE                          -- If a keyword gets deleted so should this row.
);

-- Ensure that there is ever only one unique pair of tutorial and keyword pairing. It doesn't
-- make sense for a tutorial to have 3 keywords that are all the same.
CREATE UNIQUE INDEX idx_courses_keywords ON courses_keywords(course_id, keyword_id);
