package repository_test

import (
	"context"
	"go-project-template/internal/domain"
	"go-project-template/internal/repository"
	"go-project-template/internal/testhelper"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	user_create_uuid = "79f8aa8e-f7ed-4e47-b9e4-4cd5db68a297"
	user_get_uuid    = "3fe82b1f-ab3d-40a1-8bd8-bccd4dd166f8"
)

func NewTestPostgresUser(t *testing.T) domain.UserRepository {
	t.Helper()

	ctx := context.Background()
	conn := testhelper.NewTestPgxConn(t)

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)

	repo := repository.NewUserRepository(tx)

	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
	})

	return repo
}

func TestPostgresUser_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := NewTestPostgresUser(t)

	parsedUserId, err := uuid.Parse(user_create_uuid)
	if err != nil {
		t.Error("error parsing uuid", err)
	}

	testCases := map[string]struct {
		have *domain.User
		err  bool
	}{
		"valid":        {&domain.User{UUID: parsedUserId}, false},
		"invalid uuid": {&domain.User{UUID: uuid.New()}, true},
	}

	for scenario, tc := range testCases { //nolint:paralleltest
		t.Run(scenario, func(t *testing.T) {
			_, err := repo.CreateOrUpdate(ctx, tc.have)

			if tc.err {
				assert.Error(t, err)
				return
			}

			assert.NotEqual(t, 0, tc.have.UUID)
		})
	}
}

func TestPostgresUser_GetByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo := NewTestPostgresUser(t)

	parsedUserId, err := uuid.Parse(user_get_uuid)
	if err != nil {
		t.Error("error parsing uuid from env", err)
	}

	user := &domain.User{UUID: parsedUserId, FirstName: "john", LastName: "wick"}
	// require.NoError(t, repo.CreateOrUpdate(ctx, user))

	testCases := map[string]struct {
		id   uuid.UUID
		want *domain.User
		err  error
	}{
		"valid ID":   {user.UUID, user, nil},
		"invalid ID": {uuid.New(), nil, domain.ErrNotFound},
	}

	for scenario, tc := range testCases { //nolint:paralleltest
		t.Run(scenario, func(t *testing.T) {
			dev, err := repo.GetByID(ctx, tc.id)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want.UUID, dev.UUID)
			assert.Equal(t, tc.want.FirstName, dev.FirstName)
			assert.Equal(t, tc.want.LastName, dev.LastName)
		})
	}
}
