-- Auth Domain Migration
-- Password reset functionality

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id TEXT NOT NULL PRIMARY KEY,
    token TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    redeemed_at TEXT,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE
);

-- =============================================================================
-- INDEXES FOR AUTH TABLES
-- =============================================================================

-- Password reset tokens indexes
CREATE INDEX idx_password_reset_user ON password_reset_tokens (belongs_to_user);
CREATE INDEX idx_password_reset_token ON password_reset_tokens (token);
CREATE INDEX idx_password_reset_expires ON password_reset_tokens (expires_at);
CREATE INDEX idx_password_reset_unredeemed ON password_reset_tokens (belongs_to_user, expires_at) WHERE redeemed_at IS NULL;
