package domain_test

import (
	"go-project-template/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func TestUser(t *testing.T) {
	t.Parallel()

	new_uuid := uuid.New()
	var first_name, last_name, email string = "John", "Wick", "jwick@gggggmail.com"

	user := &domain.User{
		UUID:      new_uuid,
		FirstName: first_name,
		LastName:  last_name,
		Email:     &email,
	}

	// User creation
	assert.Equal(t, new_uuid, user.UUID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Wick", user.LastName)
	assert.Equal(t, email, *user.Email)

	// User normalization
	assert.Equal(t, "john", user.NormalizedFirstName())
	assert.Equal(t, "wick", user.NormalizedLastName())

	// User validation
	assert.Equal(t, nil, user.Validate())
}
