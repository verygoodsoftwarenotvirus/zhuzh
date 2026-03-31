package manager

import (
	"context"
	"testing"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/sessions"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth/fakes"
	oauthkeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth/keys"
	oauthmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth/mock"

	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	mockpublishers "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/random"
	randommock "github.com/verygoodsoftwarenotvirus/platform/v4/random/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/reflection"
	"github.com/verygoodsoftwarenotvirus/platform/v4/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildOAuthManagerForTest(t *testing.T) *manager {
	t.Helper()

	ctx := t.Context()
	repo := &oauthmock.RepositoryMock{}
	queueCfg := &msgconfig.QueuesConfig{DataChangesTopicName: t.Name()}

	mpp := &mockpublishers.PublisherProvider{}
	mpp.On(reflection.GetMethodName(mpp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	sessionData := &sessions.ContextData{}
	sessionData.ActiveAccountID = "account-1"
	sessionData.Requester.UserID = "user-1"
	sessionFetcher := func(context.Context) (*sessions.ContextData, error) {
		return sessionData, nil
	}

	secretGen := random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider())

	m, err := NewOAuth2Manager(
		ctx,
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		secretGen,
		sessionFetcher,
		mpp,
		repo,
		queueCfg,
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mpp)

	return m.(*manager)
}

func setupExpectationsForOAuthManager(
	manager *manager,
	repoSetupFunc func(repo *oauthmock.RepositoryMock),
	secretGenSetupFunc func(gen *randommock.Generator),
	eventTypeMaps ...map[string][]string,
) []any {
	repo := &oauthmock.RepositoryMock{}
	if repoSetupFunc != nil {
		repoSetupFunc(repo)
	}
	manager.oauthRepository = repo

	expectations := []any{repo}
	if secretGenSetupFunc != nil {
		secretGen := &randommock.Generator{}
		secretGenSetupFunc(secretGen)
		manager.secretGenerator = secretGen
		expectations = append(expectations, secretGen)
	}

	mp := &mockpublishers.Publisher{}
	for _, eventTypeMap := range eventTypeMaps {
		for eventType, payload := range eventTypeMap {
			mp.On(reflection.GetMethodName(mp.PublishAsync), testutils.ContextMatcher, eventMatches(eventType, payload)).Return()
		}
	}
	manager.dataChangesPublisher = mp
	expectations = append(expectations, mp)

	return expectations
}

func TestOAuth2Manager_CreateOAuth2Client(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		om := buildOAuthManagerForTest(t)

		input := fakes.BuildFakeOAuth2ClientCreationRequestInput()
		expected := fakes.BuildFakeOAuth2Client()

		expectations := setupExpectationsForOAuthManager(
			om,
			func(repo *oauthmock.RepositoryMock) {
				repo.On(reflection.GetMethodName(repo.CreateOAuth2Client), testutils.ContextMatcher, mock.MatchedBy(func(in *oauth.OAuth2ClientDatabaseCreationInput) bool {
					return in.Name == input.Name && in.Description == input.Description && in.ClientID != "" && in.ClientSecret != ""
				})).Return(expected, nil)
			},
			func(gen *randommock.Generator) {
				gen.On(reflection.GetMethodName(gen.GenerateHexEncodedString), testutils.ContextMatcher, 16).Return("aabbccdd11223344aabbccdd11223344", nil).Once()
				gen.On(reflection.GetMethodName(gen.GenerateHexEncodedString), testutils.ContextMatcher, 16).Return("eeddccbb55443322eeddccbb55443322", nil).Once()
			},
			map[string][]string{
				oauth.OAuth2ClientCreatedServiceEventType: {oauthkeys.OAuth2ClientIDKey},
			},
		)

		actual, err := om.CreateOAuth2Client(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestOAuth2Manager_ArchiveOAuth2Client(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		om := buildOAuthManagerForTest(t)

		oauth2ClientID := fakes.BuildFakeID()

		expectations := setupExpectationsForOAuthManager(
			om,
			func(repo *oauthmock.RepositoryMock) {
				repo.On(reflection.GetMethodName(repo.ArchiveOAuth2Client), testutils.ContextMatcher, oauth2ClientID).Return(nil)
			},
			nil,
			map[string][]string{
				oauth.OAuth2ClientArchivedServiceEventType: {oauthkeys.OAuth2ClientIDKey},
			},
		)

		err := om.ArchiveOAuth2Client(ctx, oauth2ClientID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}
