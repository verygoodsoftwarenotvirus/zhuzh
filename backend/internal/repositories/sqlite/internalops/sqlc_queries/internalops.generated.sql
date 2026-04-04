-- name: DeleteExpiredOAuth2ClientTokens :execrows
DELETE FROM oauth2_client_tokens WHERE code_expires_at < strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '-1 day') AND access_expires_at < strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '-1 day') AND refresh_expires_at < strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '-1 day');

-- name: DestroyAllData :exec
DELETE FROM account_invitations;
DELETE FROM account_user_memberships;
DELETE FROM accounts;
DELETE FROM audit_log_entries;
DELETE FROM comments;
DELETE FROM issue_reports;
DELETE FROM oauth2_client_tokens;
DELETE FROM oauth2_clients;
DELETE FROM password_reset_tokens;
DELETE FROM payment_transactions;
DELETE FROM permissions;
DELETE FROM products;
DELETE FROM purchases;
DELETE FROM queue_test_messages;
DELETE FROM service_setting_configurations;
DELETE FROM service_settings;
DELETE FROM subscriptions;
DELETE FROM uploaded_media;
DELETE FROM user_avatars;
DELETE FROM user_data_disclosures;
DELETE FROM user_notifications;
DELETE FROM user_role_assignments;
DELETE FROM user_role_hierarchy;
DELETE FROM user_role_permissions;
DELETE FROM user_roles;
DELETE FROM user_sessions;
DELETE FROM users;
DELETE FROM waitlist_signups;
DELETE FROM waitlists;
DELETE FROM webhook_trigger_configs;
DELETE FROM webhook_trigger_events;
DELETE FROM webhooks;

-- name: CreateQueueTestMessage :exec
INSERT INTO queue_test_messages (id, queue_name) VALUES (sqlc.arg(id), sqlc.arg(queue_name));

-- name: AcknowledgeQueueTestMessage :exec
UPDATE queue_test_messages SET acknowledged_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = sqlc.arg(id);

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
