CREATE TABLE IF NOT EXISTS course_refunds (
    id TEXT PRIMARY KEY,                                                                -- ULID for the refund request

    user_id TEXT NOT NULL,                                                              -- Reference to the user
    purchase_id TEXT NOT NULL,                                                          -- Reference to the purchase
    status TEXT NOT NULL CHECK (status IN ('Requested', 'Processed', 'Failed')),        -- Refund status

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                      -- Request creation timestamp
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                      -- Update timestamp

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (purchase_id) REFERENCES course_purchases(id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS trigger_update_course_refunds_updated_at
AFTER UPDATE ON course_refunds
FOR EACH ROW
BEGIN
    UPDATE course_refunds SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
