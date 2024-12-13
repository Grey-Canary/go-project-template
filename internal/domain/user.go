package domain

import (
	"context"
	"go-project-template/internal/utils"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

// User base model
// @Description User base model
type User struct {
	UUID      uuid.UUID `db:"uuid" json:"uuid,omitempty" example:"3fe82b1f-ab3d-40a1-8bd8-bccd4dd166f8"`
	FirstName string    `db:"first_name" json:"first_name,omitempty" example:"John"`
	LastName  string    `db:"last_name" json:"last_name,omitempty" example:"Wick"`
	Email     *string   `db:"email" json:"email,omitempty" example:"johnwick@mail.com"`
}

func (u *User) NormalizedFirstName() string {
	return strings.ToLower(u.FirstName)
}

func (u *User) NormalizedLastName() string {
	return strings.ToLower(u.LastName)
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Email, validation.Required),
	)
}

type UserRepository interface {
	GetByID(ctx context.Context, uuid uuid.UUID) (User, error)
	CreateOrUpdate(context.Context, *User) (*User, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	GetList(ctx context.Context, pq *utils.PaginationQuery) (*utils.PaginationResponse[User], error)
}
