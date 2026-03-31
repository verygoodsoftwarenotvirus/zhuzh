package identity

import (
	"context"
	"time"
)

// WebAuthnCredential represents a stored passkey credential for a user.
type WebAuthnCredential struct {
	CreatedAt     time.Time
	LastUsedAt    *time.Time
	ArchivedAt    *time.Time
	ID            string
	BelongsToUser string
	Transports    string
	FriendlyName  string
	CredentialID  []byte
	PublicKey     []byte
	SignCount     uint32
}

// WebAuthnCredentialDataManager describes a structure which can manage WebAuthn credentials in persistent storage.
type WebAuthnCredentialDataManager interface {
	CreateWebAuthnCredential(ctx context.Context, input *WebAuthnCredentialCreationInput) (*WebAuthnCredential, error)
	GetWebAuthnCredentialsForUser(ctx context.Context, userID string) ([]*WebAuthnCredential, error)
	GetWebAuthnCredentialByCredentialID(ctx context.Context, credentialID []byte) (*WebAuthnCredential, error)
	UpdateWebAuthnCredentialSignCount(ctx context.Context, id string, signCount uint32) error
	ArchiveWebAuthnCredential(ctx context.Context, id string) error
	// ArchiveWebAuthnCredentialForUser archives a credential only if it belongs to the given user. Returns nil if no rows updated.
	ArchiveWebAuthnCredentialForUser(ctx context.Context, id, userID string) error
}

// WebAuthnCredentialCreationInput is used when creating a new WebAuthn credential.
type WebAuthnCredentialCreationInput struct {
	ID            string
	BelongsToUser string
	Transports    string
	FriendlyName  string
	CredentialID  []byte
	PublicKey     []byte
	SignCount     uint32
}
