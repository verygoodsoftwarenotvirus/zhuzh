package sessions

import (
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
)

var (
	ErrAuthenticationNotFound = platformerrors.New("authentication not found")
)
