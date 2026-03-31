package manager

import (
	"testing"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/comments"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/comments/fakes"
	commentskeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/comments/keys"
	commentsmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/comments/mock"

	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	mockpublishers "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/reflection"
	"github.com/verygoodsoftwarenotvirus/platform/v4/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildCommentsManagerForTest(t *testing.T) *commentsManager {
	t.Helper()

	ctx := t.Context()
	queueCfg := &msgconfig.QueuesConfig{
		DataChangesTopicName: t.Name(),
	}

	mpp := &mockpublishers.PublisherProvider{}
	mpp.On(reflection.GetMethodName(mpp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	m, err := NewCommentsDataManager(
		ctx,
		tracing.NewNoopTracerProvider(),
		logging.NewNoopLogger(),
		&commentsmock.Repository{},
		queueCfg,
		mpp,
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mpp)

	return m.(*commentsManager)
}

func setupExpectationsForCommentsManager(
	manager *commentsManager,
	repoSetupFunc func(repo *commentsmock.Repository),
	eventTypeMaps ...map[string][]string,
) []any {
	repo := &commentsmock.Repository{}
	if repoSetupFunc != nil {
		repoSetupFunc(repo)
	}
	manager.repo = repo

	mp := &mockpublishers.Publisher{}
	for _, eventTypeMap := range eventTypeMaps {
		for eventType, payload := range eventTypeMap {
			mp.On(reflection.GetMethodName(mp.PublishAsync), testutils.ContextMatcher, eventMatches(eventType, payload)).Return()
		}
	}
	manager.dataChangesPublisher = mp

	return []any{repo, mp}
}

func TestCommentsManager_CreateComment(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		cm := buildCommentsManagerForTest(t)

		input := fakes.BuildFakeCommentCreationRequestInput()
		expected := fakes.BuildFakeComment()

		expectations := setupExpectationsForCommentsManager(
			cm,
			func(repo *commentsmock.Repository) {
				repo.On(reflection.GetMethodName(repo.CreateComment), testutils.ContextMatcher, testutils.MatchType[*comments.CommentDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				comments.CommentCreatedServiceEventType: {commentskeys.CommentIDKey},
			},
		)

		actual, err := cm.CreateComment(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestCommentsManager_UpdateComment(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		cm := buildCommentsManagerForTest(t)

		comment := fakes.BuildFakeComment()
		commentID := comment.ID
		belongsToUser := comment.BelongsToUser
		input := &comments.CommentUpdateRequestInput{Content: "Updated content"}

		expectations := setupExpectationsForCommentsManager(
			cm,
			func(repo *commentsmock.Repository) {
				repo.On(reflection.GetMethodName(repo.UpdateComment), testutils.ContextMatcher, commentID, belongsToUser, input.Content).Return(nil)
			},
			map[string][]string{
				comments.CommentUpdatedServiceEventType: {commentskeys.CommentIDKey},
			},
		)

		err := cm.UpdateComment(ctx, commentID, belongsToUser, input)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestCommentsManager_ArchiveComment(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		cm := buildCommentsManagerForTest(t)

		commentID := fakes.BuildFakeID()

		expectations := setupExpectationsForCommentsManager(
			cm,
			func(repo *commentsmock.Repository) {
				repo.On(reflection.GetMethodName(repo.ArchiveComment), testutils.ContextMatcher, commentID).Return(nil)
			},
			map[string][]string{
				comments.CommentArchivedServiceEventType: {commentskeys.CommentIDKey},
			},
		)

		err := cm.ArchiveComment(ctx, commentID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}
