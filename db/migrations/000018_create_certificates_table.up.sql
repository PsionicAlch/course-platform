CREATE TABLE IF NOT EXISTS certificates (
    id TEXT PRIMARY KEY,

    user_id TEXT NOT NULL,
    course_id TEXT NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_certificates_user_id_course_id ON certificates(user_id, course_id);
