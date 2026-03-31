-- Idempotently create per-service database users for platform services.
-- In production these already exist (created by Terraform); in localdev they are created here.
-- All services get full access — the goal is per-service identity in pg_stat_activity,
-- not fine-grained privilege isolation.

DO $$
DECLARE
  service_users TEXT[] := ARRAY[
    'api_db_user',
    'async_message_handler',
    'db_cleaner',
    'search_data_index_scheduler',
    'mobile_notification_scheduler',
    'queue_test'
  ];
  u TEXT;
BEGIN
  FOREACH u IN ARRAY service_users LOOP
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = u) THEN
      EXECUTE format('CREATE ROLE %I WITH LOGIN', u);
    END IF;
  END LOOP;
END $$;

-- Grant on all existing tables and sequences.
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO
  api_db_user,
  async_message_handler,
  db_cleaner,
  search_data_index_scheduler,
  mobile_notification_scheduler,
  queue_test;

GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO
  api_db_user,
  async_message_handler,
  db_cleaner,
  search_data_index_scheduler,
  mobile_notification_scheduler,
  queue_test;

-- Ensure future tables/sequences created by the migration user also get grants.
-- This covers tables created by later migrations (e.g. meal planning).
ALTER DEFAULT PRIVILEGES IN SCHEMA public
  GRANT ALL PRIVILEGES ON TABLES TO
    api_db_user,
    async_message_handler,
    db_cleaner,
    search_data_index_scheduler,
    mobile_notification_scheduler,
    queue_test;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
  GRANT ALL PRIVILEGES ON SEQUENCES TO
    api_db_user,
    async_message_handler,
    db_cleaner,
    search_data_index_scheduler,
    mobile_notification_scheduler,
    queue_test;
