-- OAuth Domain Migration
-- OAuth2 client and token management

CREATE TABLE IF NOT EXISTS oauth2_clients (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    client_id TEXT NOT NULL,
    client_secret TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    archived_at TEXT,
    UNIQUE(client_id),
    UNIQUE(client_secret)
);

CREATE TABLE IF NOT EXISTS oauth2_client_tokens (
    id TEXT NOT NULL PRIMARY KEY,
    client_id TEXT NOT NULL REFERENCES oauth2_clients("client_id") ON DELETE CASCADE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    redirect_uri TEXT DEFAULT '' NOT NULL,
    code TEXT DEFAULT '' NOT NULL,
    code_challenge TEXT DEFAULT '' NOT NULL,
    code_challenge_method TEXT DEFAULT '' NOT NULL,
    code_created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    code_expires_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '+1 hour')) NOT NULL,
    access TEXT DEFAULT '' NOT NULL,
    access_created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    access_expires_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '+1 hour')) NOT NULL,
    refresh TEXT DEFAULT '' NOT NULL,
    refresh_created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    refresh_expires_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '+1 hour')) NOT NULL,
    UNIQUE(belongs_to_user, client_id, code_expires_at, access_expires_at, refresh_expires_at)
);

-- =============================================================================
-- INDEXES FOR OAUTH TABLES
-- =============================================================================

-- OAuth2 clients indexes
CREATE INDEX idx_oauth2_clients_archived_at ON oauth2_clients (archived_at) WHERE archived_at IS NULL;

-- OAuth2 client tokens indexes
CREATE INDEX idx_oauth2_tokens_client_id ON oauth2_client_tokens (client_id);
CREATE INDEX idx_oauth2_tokens_user ON oauth2_client_tokens (belongs_to_user);
CREATE INDEX idx_oauth2_tokens_user_client ON oauth2_client_tokens (belongs_to_user, client_id);
CREATE INDEX idx_oauth2_tokens_code_expires ON oauth2_client_tokens (code_expires_at);
CREATE INDEX idx_oauth2_tokens_access_expires ON oauth2_client_tokens (access_expires_at);
CREATE INDEX idx_oauth2_tokens_refresh_expires ON oauth2_client_tokens (refresh_expires_at);
