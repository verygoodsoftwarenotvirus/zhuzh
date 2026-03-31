package integration

import (
	"testing"

	commentsgrpc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/comments"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func commentsServiceCreateCommentOnRecipe(t *testing.T, recipeID string, c client.Client, content string) *commentsgrpc.Comment {
	t.Helper()
	ctx := t.Context()

	if content == "" {
		content = "test comment via CommentsService"
	}

	res, err := c.CommentsService().CreateComment(ctx, &commentsgrpc.CreateCommentRequest{
		Input: &commentsgrpc.CommentCreationRequestInput{
			Content:      content,
			TargetType:   "recipes",
			ReferencedId: recipeID,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Comment)

	return res.Comment
}

func TestCommentsService_CreateComment(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(t)
		referencedID := nonexistentID

		res, err := testClient.CommentsService().CreateComment(ctx, &commentsgrpc.CreateCommentRequest{
			Input: &commentsgrpc.CommentCreationRequestInput{
				Content:      "created via CreateComment",
				TargetType:   "recipes",
				ReferencedId: referencedID,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res.Comment)
		assert.Equal(t, "recipes", res.Comment.TargetType)
		assert.Equal(t, referencedID, res.Comment.ReferencedId)
		assert.Equal(t, "created via CreateComment", res.Comment.Content)

		AssertAuditLogContainsFuzzyForUser(t, ctx, testClient, user.ID, 10, []*ExpectedAuditEntry{
			{EventType: "created", ResourceType: "comments", RelevantID: res.Comment.Id},
		})
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		res, err := c.CommentsService().CreateComment(ctx, &commentsgrpc.CreateCommentRequest{
			Input: &commentsgrpc.CommentCreationRequestInput{
				Content:      "test",
				TargetType:   "recipes",
				ReferencedId: nonexistentID,
			},
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCommentsService_GetCommentsForReference(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		referencedID := nonexistentID
		_ = commentsServiceCreateCommentOnRecipe(t, referencedID, testClient, "")

		listRes, err := testClient.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   "recipes",
			ReferencedId: referencedID,
		})
		require.NoError(t, err)
		require.NotNil(t, listRes)
		assert.GreaterOrEqual(t, len(listRes.Data), 1)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		referencedID := nonexistentID
		_ = commentsServiceCreateCommentOnRecipe(t, referencedID, testClient, "")

		c := buildUnauthenticatedGRPCClientForTest(t)
		listRes, err := c.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   "recipes",
			ReferencedId: referencedID,
		})
		assert.Error(t, err)
		assert.Nil(t, listRes)
	})
}

func TestCommentsService_UpdateComment(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(t)
		referencedID := nonexistentID
		createdComment := commentsServiceCreateCommentOnRecipe(t, referencedID, testClient, "original")

		_, err := testClient.CommentsService().UpdateComment(ctx, &commentsgrpc.UpdateCommentRequest{
			CommentId: createdComment.Id,
			Input:     &commentsgrpc.CommentUpdateRequestInput{Content: "updated via CommentsService"},
		})
		require.NoError(t, err)

		AssertAuditLogContainsFuzzyForUser(t, ctx, testClient, user.ID, 15, []*ExpectedAuditEntry{
			{EventType: "created", ResourceType: "comments", RelevantID: createdComment.Id},
			{EventType: "updated", ResourceType: "comments", RelevantID: createdComment.Id},
		})

		listRes, err := testClient.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   "recipes",
			ReferencedId: referencedID,
		})
		require.NoError(t, err)
		for _, c := range listRes.Data {
			if c.Id == createdComment.Id {
				assert.Equal(t, "updated via CommentsService", c.Content)
				break
			}
		}
	})
}

func TestCommentsService_ArchiveComment(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(t)
		referencedID := nonexistentID
		createdComment := commentsServiceCreateCommentOnRecipe(t, referencedID, testClient, "")

		_, err := testClient.CommentsService().ArchiveComment(ctx, &commentsgrpc.ArchiveCommentRequest{
			CommentId: createdComment.Id,
		})
		require.NoError(t, err)

		AssertAuditLogContainsFuzzyForUser(t, ctx, testClient, user.ID, 15, []*ExpectedAuditEntry{
			{EventType: "created", ResourceType: "comments", RelevantID: createdComment.Id},
			{EventType: "archived", ResourceType: "comments", RelevantID: createdComment.Id},
		})

		listRes, err := testClient.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   "recipes",
			ReferencedId: referencedID,
		})
		require.NoError(t, err)
		for _, c := range listRes.Data {
			assert.NotEqual(t, createdComment.Id, c.Id, "archived comment should not appear")
		}
	})
}
