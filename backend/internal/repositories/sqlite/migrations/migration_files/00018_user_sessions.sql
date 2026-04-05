CREATE TABLE IF NOT EXISTS user_sessions (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    session_token_id TEXT NOT NULL,
    refresh_token_id TEXT NOT NULL,
    client_ip TEXT DEFAULT '' NOT NULL,
    user_agent TEXT DEFAULT '' NOT NULL,
    device_name TEXT DEFAULT '' NOT NULL,
    login_method TEXT DEFAULT '' NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_active_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    expires_at TEXT NOT NULL,
    revoked_at TEXT
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user ON user_sessions (belongs_to_user) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_sessions_token_id ON user_sessions (session_token_id) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_token_id ON user_sessions (refresh_token_id) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires ON user_sessions (expires_at);
