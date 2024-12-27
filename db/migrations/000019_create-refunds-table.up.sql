CREATE TABLE IF NOT EXISTS refunds (
    id TEXT PRIMARY KEY,

    user_id TEXT NOT NULL,
    course_purchase_id TEXT NOT NULL,
    STATUS TEXT
);
