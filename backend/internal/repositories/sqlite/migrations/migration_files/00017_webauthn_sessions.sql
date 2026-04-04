-- WebAuthn Sessions Migration
-- Passkey (WebAuthn/FIDO2) session storage for registration/authentication ceremonies

CREATE TABLE IF NOT EXISTS webauthn_sessions (
    challenge TEXT NOT NULL PRIMARY KEY,
    session_data TEXT NOT NULL,
    expires_at TEXT NOT NULL
);

CREATE INDEX idx_webauthn_sessions_expires_at ON webauthn_sessions (expires_at);
