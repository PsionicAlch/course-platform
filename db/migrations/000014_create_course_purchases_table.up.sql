CREATE TABLE IF NOT EXISTS course_purchases (
    id TEXT PRIMARY KEY,                                                                                            -- ULID for the purchase record
    user_id TEXT NOT NULL,                                                                                          -- Reference to the user
    course_id TEXT NOT NULL,                                                                                        -- Reference to the purchased course
    payment_key TEXT NOT NULL,                                                                                      -- Used to validate whether a payment was successful.
    stripe_checkout_session_id TEXT NOT NULL,                                                                       -- Stripe Checkout Session ID used for getting more information on a payment and issuing a refund.
    affiliate_code TEXT,                                                                                            -- Optional affiliate code used
    discount_code TEXT,                                                                                             -- Optional discount code used
    affiliate_points_used INTEGER DEFAULT 0 CHECK (affiliate_points_used >= 0),                                     -- Points used
    amount_paid REAL NOT NULL CHECK (amount_paid >= 0.0),                                                           -- Final amount paid in cents
    payment_status TEXT NOT NULL DEFAULT 'Pending' CHECK (
        payment_status IN ('Succeeded', 'Refunded', 'Pending', 'Cancelled', 'Failed', 'Requires Action')
    ),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                                                  -- Purchase timestamp
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                                                  -- Update timestamp

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_course_purchases_discount_code ON course_purchases(discount_code);

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_purchases_stripe_checkout_session_id ON course_purchases(stripe_checkout_session_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_course_purchases_payment_key ON course_purchases(payment_key);

CREATE TRIGGER IF NOT EXISTS trigger_update_course_purchases_updated_at
AFTER UPDATE ON course_purchases
FOR EACH ROW
BEGIN
    UPDATE course_purchases SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
