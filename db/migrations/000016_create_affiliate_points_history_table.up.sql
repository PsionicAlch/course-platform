CREATE TABLE IF NOT EXISTS affiliate_points_history (
    id TEXT PRIMARY KEY,                                                                              -- ULID as the primary key (stored as TEXT)

    user_id TEXT NOT NULL,                                                                            -- Reference to the 'users' table
    course_id TEXT NOT NULL,                                                                          -- Reference to the 'courses' table
    points_change INTEGER NOT NULL,                                                                   -- Positive or negative change in affiliate points
    reason TEXT NOT NULL,                                                                             -- Description or reason for the change (e.g., "Course purchase", "Refund")

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                                    -- Timestamp of the points change

    FOREIGN KEY (user_id) REFERENCES users(id),                                                       -- Ensures the user exists in the 'users' table
    FOREIGN KEY (course_id) REFERENCES courses(id)
);

-- Create an index on user_id to allow quick lookups by user
CREATE INDEX IF NOT EXISTS idx_affiliate_points_user_id ON affiliate_points_history(user_id);

-- Trigger to calculate user's affiliate points when a new affiliate_points_history row gets inserted.
CREATE TRIGGER IF NOT EXISTS trigger_recalculate_user_affiliate_points_after_insert
AFTER INSERT ON affiliate_points_history
FOR EACH ROW
BEGIN
    UPDATE users
    SET affiliate_points = (
        SELECT COALESCE(SUM(points_change), 0)
        FROM affiliate_points_history
        WHERE user_id = NEW.user_id
    )
    WHERE id = NEW.user_id;
END;

-- Trigger to calculate user's affiliate points when a row in affiliate_points_history gets updated.
CREATE TRIGGER IF NOT EXISTS trigger_recalculate_user_affiliate_points_after_update
AFTER UPDATE ON affiliate_points_history
FOR EACH ROW
BEGIN
    UPDATE users
    SET affiliate_points = (
        SELECT COALESCE(SUM(points_change), 0)
        FROM affiliate_points_history
        WHERE user_id = NEW.user_id
    )
    WHERE id = NEW.user_id;
END;

-- Trigger to calculate user's affiliate points when a row in affiliate_points_history gets deleted.
CREATE TRIGGER IF NOT EXISTS trigger_recalculate_user_affiliate_points_after_delete
AFTER DELETE ON affiliate_points_history
FOR EACH ROW
BEGIN
    UPDATE users
    SET affiliate_points = (
        SELECT COALESCE(SUM(points_change), 0)
        FROM affiliate_points_history
        WHERE user_id = OLD.user_id
    )
    WHERE id = OLD.user_id;
END;
