package domain

import (
	"context"
	"time"

	"go-rest-api-boilerplate/internal/model/reqres"
)

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	UpdateByID(ctx context.Context, id int64, user *User) error
	DeleteByID(ctx context.Context, id int64) error
	FindAll(ctx context.Context) (*[]User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}

type UserService interface {
	Create(ctx context.Context, req *reqres.CreateUserReq) error
	UpdateByID(ctx context.Context, id int64, req *reqres.UpdateUserReq) error
	DeleteByID(ctx context.Context, id int64) error
	FindAll(ctx context.Context) (*[]User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}
