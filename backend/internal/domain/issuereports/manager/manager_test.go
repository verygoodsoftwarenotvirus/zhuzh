package manager

import (
	"errors"
	"testing"

	types "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports/fakes"
	issuereportkeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports/keys"
	issuereportsmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports/mock"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
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

func buildIssueReportsManagerForTest(t *testing.T) (*issueReportsManager, *issuereportsmock.Repository) {
	t.Helper()

	ctx := t.Context()
	repo := &issuereportsmock.Repository{}
	queueCfg := &msgconfig.QueuesConfig{DataChangesTopicName: t.Name()}

	mpp := &mockpublishers.PublisherProvider{}
	mpp.On(reflection.GetMethodName(mpp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	m, err := NewIssueReportsDataManager(ctx, tracing.NewNoopTracerProvider(), logging.NewNoopLogger(), repo, queueCfg, mpp)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mpp)

	return m.(*issueReportsManager), repo
}

func setupExpectationsForIssueReportsManager(
	manager *issueReportsManager,
	repoSetupFunc func(repo *issuereportsmock.Repository),
	eventTypeMaps ...map[string][]string,
) []any {
	repo := &issuereportsmock.Repository{}
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

func TestIssueReportsDataManager_CreateIssueReport(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		manager, _ := buildIssueReportsManagerForTest(t)

		dbInput := fakes.BuildFakeIssueReportDatabaseCreationInput()
		createdReport := fakes.BuildFakeIssueReport()
		createdReport.ID = dbInput.ID

		expectations := setupExpectationsForIssueReportsManager(
			manager,
			func(repo *issuereportsmock.Repository) {
				repo.On(reflection.GetMethodName(repo.CreateIssueReport), testutils.ContextMatcher, mock.Anything).Return(createdReport, nil)
			},
			map[string][]string{
				types.IssueReportCreatedServiceEventType: {issuereportkeys.IssueReportIDKey},
			},
		)

		created, err := manager.CreateIssueReport(ctx, dbInput)

		require.NoError(t, err)
		assert.NotNil(t, created)
		assert.Equal(t, dbInput.ID, created.ID)
		mock.AssertExpectationsForObjects(t, expectations...)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		manager, _ := buildIssueReportsManagerForTest(t)

		dbInput := fakes.BuildFakeIssueReportDatabaseCreationInput()

		expectations := setupExpectationsForIssueReportsManager(
			manager,
			func(repo *issuereportsmock.Repository) {
				repo.On(reflection.GetMethodName(repo.CreateIssueReport), testutils.ContextMatcher, mock.Anything).Return((*types.IssueReport)(nil), errors.New("db error"))
			},
		)

		created, err := manager.CreateIssueReport(ctx, dbInput)

		assert.Error(t, err)
		assert.Nil(t, created)
		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestIssueReportsDataManager_GetIssueReport(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		manager, repo := buildIssueReportsManagerForTest(t)

		expected := fakes.BuildFakeIssueReport()
		repo.On(reflection.GetMethodName(repo.GetIssueReport), testutils.ContextMatcher, expected.ID).Return(expected, nil)

		result, err := manager.GetIssueReport(ctx, expected.ID)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mock.AssertExpectationsForObjects(t, repo)
	})
}

func TestIssueReportsDataManager_GetIssueReports(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		manager, repo := buildIssueReportsManagerForTest(t)

		filter := filtering.DefaultQueryFilter()
		report := fakes.BuildFakeIssueReport()
		expected := &filtering.QueryFilteredResult[types.IssueReport]{
			Data: []*types.IssueReport{report},
		}
		repo.On(reflection.GetMethodName(repo.GetIssueReports), testutils.ContextMatcher, filter).Return(expected, nil)

		result, err := manager.GetIssueReports(ctx, filter)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mock.AssertExpectationsForObjects(t, repo)
	})
}

func TestIssueReportsDataManager_UpdateIssueReport(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		manager, _ := buildIssueReportsManagerForTest(t)

		issueReport := fakes.BuildFakeIssueReport()

		expectations := setupExpectationsForIssueReportsManager(
			manager,
			func(repo *issuereportsmock.Repository) {
				repo.On(reflection.GetMethodName(repo.UpdateIssueReport), testutils.ContextMatcher, issueReport).Return(nil)
			},
			map[string][]string{
				types.IssueReportUpdatedServiceEventType: {issuereportkeys.IssueReportIDKey},
			},
		)

		err := manager.UpdateIssueReport(ctx, issueReport)

		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestIssueReportsDataManager_ArchiveIssueReport(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		manager, _ := buildIssueReportsManagerForTest(t)

		issueReportID := fakes.BuildFakeID()

		expectations := setupExpectationsForIssueReportsManager(
			manager,
			func(repo *issuereportsmock.Repository) {
				repo.On(reflection.GetMethodName(repo.ArchiveIssueReport), testutils.ContextMatcher, issueReportID).Return(nil)
			},
			map[string][]string{
				types.IssueReportArchivedServiceEventType: {issuereportkeys.IssueReportIDKey},
			},
		)

		err := manager.ArchiveIssueReport(ctx, issueReportID)

		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, expectations...)
	})
}
