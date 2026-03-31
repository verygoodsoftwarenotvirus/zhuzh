package authentication

import (
	"context"

	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
)

var (
	// ErrInvalidTOTPToken indicates that a provided two-factor code is invalid.
	ErrInvalidTOTPToken = platformerrors.New("invalid two factor code")
	// ErrTOTPRequired indicates that the user has TOTP enabled but did not provide a code.
	ErrTOTPRequired = platformerrors.New("TOTP code required but not provided")
)

type (
	// Authenticator authenticates users.
	Authenticator interface {
		Hasher

		CredentialsAreValid(ctx context.Context, hash, password, totpSecret, totpCode string) (bool, error)
	}
)
