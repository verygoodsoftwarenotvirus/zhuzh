package authentication

import (
	"encoding/base64"
	"net/http"
	"testing"

	mockauthn "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/mock"
	tokenscfg "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/tokens/config"
	identitymanagermock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/manager/mock"
	oauthmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth/mock"

	noopanalytics "github.com/verygoodsoftwarenotvirus/platform/v4/analytics/noop"
	"github.com/verygoodsoftwarenotvirus/platform/v4/encoding"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	mockpublishers "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/reflection"
	mockrouting "github.com/verygoodsoftwarenotvirus/platform/v4/routing/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	ctx := t.Context()
	logger := logging.NewNoopLogger()
	encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	cfg := &Config{
		Tokens: tokenscfg.Config{
			Provider:                tokenscfg.ProviderJWT,
			Audience:                "",
			Base64EncodedSigningKey: base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
		},
	}
	queueCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

	pp := &mockpublishers.PublisherProvider{}
	pp.On(reflection.GetMethodName(pp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	rpm := mockrouting.NewRouteParamManager()
	rpm.On(
		"BuildRouteParamStringIDFetcher",
		AuthProviderParamKey,
	).Return(func(*http.Request) string { return "" })

	s, err := ProvideService(
		ctx,
		logger,
		cfg,
		&mockauthn.Authenticator{},
		&oauthmock.RepositoryMock{},
		&identitymanagermock.IdentityDataManager{},
		encoderDecoder,
		tracing.NewNoopTracerProvider(),
		pp,
		noopanalytics.NewEventReporter(),
		rpm,
		queueCfg,
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, pp, rpm)

	return s.(*service)
}

func TestProvideService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		logger := logging.NewNoopLogger()
		encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		cfg := &Config{
			Tokens: tokenscfg.Config{
				Provider:                tokenscfg.ProviderJWT,
				Audience:                "",
				Base64EncodedSigningKey: base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
			},
		}
		queueCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		pp := &mockpublishers.PublisherProvider{}
		pp.On(reflection.GetMethodName(pp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			AuthProviderParamKey,
		).Return(func(*http.Request) string { return "" })

		s, err := ProvideService(
			ctx,
			logger,
			cfg,
			&mockauthn.Authenticator{},
			&oauthmock.RepositoryMock{},
			&identitymanagermock.IdentityDataManager{},
			encoderDecoder,
			tracing.NewNoopTracerProvider(),
			pp,
			noopanalytics.NewEventReporter(),
			rpm,
			queueCfg,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})
}
