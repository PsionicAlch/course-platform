DROP TRIGGER IF EXISTS trigger_recalculate_user_affiliate_points_after_delete;

DROP TRIGGER IF EXISTS trigger_recalculate_user_affiliate_points_after_update;

DROP TRIGGER IF EXISTS trigger_recalculate_user_affiliate_points_after_insert;

DROP INDEX IF EXISTS idx_affiliate_points_user_id;

DROP TABLE IF EXISTS affiliate_points_history;
