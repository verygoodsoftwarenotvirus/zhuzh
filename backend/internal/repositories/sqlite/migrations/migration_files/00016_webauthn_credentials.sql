-- WebAuthn Credentials Migration
-- Passkey (WebAuthn/FIDO2) credential storage for passwordless authentication

CREATE TABLE IF NOT EXISTS webauthn_credentials (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    credential_id BLOB NOT NULL,
    public_key BLOB NOT NULL,
    sign_count INTEGER NOT NULL DEFAULT 0,
    transports TEXT DEFAULT '' NOT NULL,
    friendly_name TEXT DEFAULT '' NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_used_at TEXT,
    archived_at TEXT
);

CREATE UNIQUE INDEX idx_webauthn_credentials_credential_id_active
    ON webauthn_credentials (credential_id)
    WHERE archived_at IS NULL;

CREATE INDEX idx_webauthn_credentials_user ON webauthn_credentials (belongs_to_user) WHERE archived_at IS NULL;
