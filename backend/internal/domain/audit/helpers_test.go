package audit

import (
	"context"
	"testing"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/sessions"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	identityfakes "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/fakes"

	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"

	"github.com/stretchr/testify/assert"
)

func Test_buildDataChangeMessageFromContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		sessionContextData := &sessions.ContextData{
			Requester:       sessions.RequesterInfo{UserID: identityfakes.BuildFakeID()},
			ActiveAccountID: identityfakes.BuildFakeID(),
		}
		ctx = context.WithValue(ctx, sessions.SessionContextDataKey, sessionContextData)

		expected := &DataChangeMessage{
			EventType: identity.UserSignedUpServiceEventType,
			Context: map[string]any{
				"things": "stuff",
			},
			UserID:    sessionContextData.Requester.UserID,
			AccountID: sessionContextData.ActiveAccountID,
		}

		actual := BuildDataChangeMessageFromContext(ctx, logging.NewNoopLogger(), expected.EventType, expected.Context)

		assert.Equal(t, expected, actual)
	})
}
