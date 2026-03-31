package emaildeliverabilitytest

import (
	"errors"
	"testing"

	emailmock "github.com/verygoodsoftwarenotvirus/platform/v4/email/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/reflection"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewJob(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		job, err := NewJob(
			&emailmock.Emailer{},
			logging.NewNoopLogger(),
			tracing.NewNoopTracerProvider(),
			&JobParams{
				RecipientEmailAddress: "test@example.com",
				ServiceEnvironment:    "test",
			},
		)

		require.NoError(t, err)
		assert.NotNil(t, job)
	})

	T.Run("defaults service environment to prod when empty", func(t *testing.T) {
		t.Parallel()

		job, err := NewJob(
			&emailmock.Emailer{},
			logging.NewNoopLogger(),
			tracing.NewNoopTracerProvider(),
			&JobParams{
				RecipientEmailAddress: "test@example.com",
			},
		)

		require.NoError(t, err)
		assert.NotNil(t, job)
	})

	T.Run("returns error when recipient email is empty", func(t *testing.T) {
		t.Parallel()

		job, err := NewJob(
			&emailmock.Emailer{},
			logging.NewNoopLogger(),
			tracing.NewNoopTracerProvider(),
			&JobParams{},
		)

		require.Error(t, err)
		assert.Nil(t, job)
	})
}

func TestJob_Do(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		emailer := &emailmock.Emailer{}
		emailer.On(reflection.GetMethodName(emailer.SendEmail), mock.Anything, mock.Anything).Return(nil)

		job, err := NewJob(
			emailer,
			logging.NewNoopLogger(),
			tracing.NewNoopTracerProvider(),
			&JobParams{
				RecipientEmailAddress: "test@example.com",
				ServiceEnvironment:    "test",
			},
		)
		require.NoError(t, err)

		err = job.Do(ctx)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, emailer)
	})

	T.Run("returns error when emailer fails", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		expectedErr := errors.New("email send failure")

		emailer := &emailmock.Emailer{}
		emailer.On(reflection.GetMethodName(emailer.SendEmail), mock.Anything, mock.Anything).Return(expectedErr)

		job, err := NewJob(
			emailer,
			logging.NewNoopLogger(),
			tracing.NewNoopTracerProvider(),
			&JobParams{
				RecipientEmailAddress: "test@example.com",
				ServiceEnvironment:    "test",
			},
		)
		require.NoError(t, err)

		err = job.Do(ctx)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, emailer)
	})
}
