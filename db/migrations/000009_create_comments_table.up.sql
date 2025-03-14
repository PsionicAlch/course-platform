CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,

    content TEXT NOT NULL,
    user_id TEXT NOT NULL,
    tutorial_id TEXT NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tutorial_id) REFERENCES tutorials(id) ON DELETE CASCADE
);
