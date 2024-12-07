DROP TRIGGER IF EXISTS trigger_update_course_purchases_updated_at;

DROP INDEX IF EXISTS idx_course_purchases_payment_key;

DROP INDEX IF EXISTS idx_course_purchases_stripe_checkout_session_id;

DROP INDEX IF EXISTS idx_course_purchases_discount_code;

DROP TABLE IF EXISTS course_purchases;
