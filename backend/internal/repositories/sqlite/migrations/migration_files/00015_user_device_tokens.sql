-- User Device Tokens Migration
-- Push notification device token storage (APNs for iOS, FCM for Android)

CREATE TABLE IF NOT EXISTS user_device_tokens (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    device_token TEXT NOT NULL,
    platform TEXT NOT NULL,
    created_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')) NOT NULL,
    last_updated_at TEXT,
    archived_at TEXT
);

-- Partial unique: only one active (non-archived) token per user+device_token.
-- Allows re-registration after archive by inserting new row or un-archiving.
CREATE UNIQUE INDEX idx_user_device_tokens_user_token_active
    ON user_device_tokens (belongs_to_user, device_token)
    WHERE archived_at IS NULL;

CREATE INDEX idx_user_device_tokens_user ON user_device_tokens (belongs_to_user) WHERE archived_at IS NULL;
