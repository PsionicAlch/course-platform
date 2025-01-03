DROP TRIGGER IF EXISTS trigger_update_refunds_updated_at;

DROP INDEX IF EXISTS idx_refunds_user_id_course_purchase_id;

DROP TABLE IF EXISTS refunds;
