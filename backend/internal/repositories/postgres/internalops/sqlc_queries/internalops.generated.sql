-- name: DeleteExpiredOAuth2ClientTokens :execrows
DELETE FROM oauth2_client_tokens WHERE code_expires_at < (SELECT NOW() - '1 day'::INTERVAL) AND access_expires_at < (SELECT NOW() - '1 day'::INTERVAL) AND refresh_expires_at < (SELECT NOW() - '1 day'::INTERVAL);

-- name: DestroyAllData :exec
TRUNCATE account_invitations, account_user_memberships, accounts, audit_log_entries, comments, issue_reports, oauth2_client_tokens, oauth2_clients, password_reset_tokens, payment_transactions, permissions, products, purchases, queue_test_messages, service_setting_configurations, service_settings, subscriptions, uploaded_media, user_avatars, user_data_disclosures, user_notifications, user_role_assignments, user_role_hierarchy, user_role_permissions, user_roles, user_sessions, users, waitlist_signups, waitlists, webhook_trigger_configs, webhook_trigger_events, webhooks CASCADE;

-- name: CreateQueueTestMessage :exec
INSERT INTO queue_test_messages (id, queue_name) VALUES (sqlc.arg(id), sqlc.arg(queue_name));

-- name: AcknowledgeQueueTestMessage :exec
UPDATE queue_test_messages SET acknowledged_at = NOW() WHERE id = sqlc.arg(id);

-- name: GetQueueTestMessage :one
SELECT id, queue_name, created_at, acknowledged_at FROM queue_test_messages WHERE id = sqlc.arg(id);

-- name: PruneQueueTestMessages :exec
DELETE FROM queue_test_messages AS qtm
WHERE qtm.queue_name = sqlc.arg(queue_name)
  AND qtm.id NOT IN (
      SELECT keep.id FROM queue_test_messages AS keep
      WHERE keep.queue_name = sqlc.arg(queue_name)
      ORDER BY keep.created_at DESC
      LIMIT 100
  );
