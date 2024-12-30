CREATE TABLE IF NOT EXISTS refunds (
    id TEXT PRIMARY KEY,

    user_id TEXT NOT NULL,
    course_purchase_id TEXT NOT NULL,
    refund_status TEXT NOT NULL DEFAULT 'Refund Pending' CHECK (
        refund_status in ('Refund Pending', 'Refund Requires Action', 'Refund Succeeded', 'Refund Failed', 'Refund Cancelled', 'Dispute Warning Needs Response', 'Dispute Warning Under Review', 'Dispute Warning Closed', 'Dispute Needs Response', 'Dispute Under Review', 'Dispute Won', 'Dispute Lost')
    ),

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_purchase_id) REFERENCES course_purchases(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_refunds_user_id_course_purchase_id ON refunds(user_id, course_purchase_id);

CREATE TRIGGER IF NOT EXISTS trigger_update_refunds_updated_at
AFTER UPDATE ON refunds
FOR EACH ROW
BEGIN
    UPDATE refunds SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
