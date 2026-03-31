-- name: CreateWebAuthnCredential :exec
INSERT INTO webauthn_credentials (
    id,
    belongs_to_user,
    credential_id,
    public_key,
    sign_count,
    transports,
    friendly_name
) VALUES (
    sqlc.arg(id),
    sqlc.arg(belongs_to_user),
    sqlc.arg(credential_id),
    sqlc.arg(public_key),
    sqlc.arg(sign_count),
    sqlc.arg(transports),
    sqlc.arg(friendly_name)
);

-- name: GetWebAuthnCredentialsForUser :many
SELECT
    id,
    belongs_to_user,
    credential_id,
    public_key,
    sign_count,
    transports,
    friendly_name,
    created_at,
    last_used_at
FROM webauthn_credentials
WHERE archived_at IS NULL
    AND belongs_to_user = sqlc.arg(belongs_to_user);

-- name: GetWebAuthnCredentialByCredentialID :one
SELECT
    id,
    belongs_to_user,
    credential_id,
    public_key,
    sign_count,
    transports,
    friendly_name,
    created_at,
    last_used_at
FROM webauthn_credentials
WHERE archived_at IS NULL
    AND credential_id = sqlc.arg(credential_id);

-- name: UpdateWebAuthnCredentialSignCount :execrows
UPDATE webauthn_credentials SET
    sign_count = sqlc.arg(sign_count),
    last_used_at = NOW()
WHERE archived_at IS NULL
    AND id = sqlc.arg(id);

-- name: ArchiveWebAuthnCredential :execrows
UPDATE webauthn_credentials SET archived_at = NOW()
WHERE archived_at IS NULL
    AND id = sqlc.arg(id);

-- name: ArchiveWebAuthnCredentialForUser :execrows
UPDATE webauthn_credentials SET archived_at = NOW()
WHERE archived_at IS NULL
    AND id = sqlc.arg(id)
    AND belongs_to_user = sqlc.arg(belongs_to_user);
