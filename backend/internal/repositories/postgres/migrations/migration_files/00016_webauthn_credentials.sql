-- WebAuthn Credentials Migration
-- Passkey (WebAuthn/FIDO2) credential storage for passwordless authentication

CREATE TABLE IF NOT EXISTS webauthn_credentials (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    credential_id BYTEA NOT NULL,
    public_key BYTEA NOT NULL,
    sign_count INTEGER NOT NULL DEFAULT 0,
    transports TEXT DEFAULT ''::TEXT NOT NULL,
    friendly_name TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_webauthn_credentials_credential_id_active
    ON webauthn_credentials (credential_id)
    WHERE archived_at IS NULL;

CREATE INDEX idx_webauthn_credentials_user ON webauthn_credentials (belongs_to_user) WHERE archived_at IS NULL;
