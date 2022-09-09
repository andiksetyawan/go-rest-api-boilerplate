package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-boilerplate/internal/domain"
	"go-rest-api-boilerplate/internal/domain/mocks"
	"go-rest-api-boilerplate/internal/model/reqres"
	"go-rest-api-boilerplate/internal/usecase/service"
)

func TestNewUserService(t *testing.T) {
	svc := service.NewUserService(nil)
	assert.NotNil(t, svc)
}

func TestUserService_FindAll(t *testing.T) {
	mockUsersResult := []domain.User{
		{
			ID:        1,
			FirstName: "john",
			LastName:  "due",
			Email:     "john@email.test",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("FindAll", mock.Anything).Return(&mockUsersResult, nil)

		svc := service.NewUserService(repo)
		users, err := svc.FindAll(context.TODO())
		assert.NoError(t, err)

		assert.Equal(t, "john", (*users)[0].FirstName)
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("FindAll", mock.Anything).Return(nil, errors.New("Unexpexted Error"))

		svc := service.NewUserService(repo)
		users, err := svc.FindAll(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, users)
	})
}

func TestUserService_FindByID(t *testing.T) {
	mockUserResult := domain.User{
		ID:        1,
		FirstName: "john",
		LastName:  "due",
		Email:     "john@email.test",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(&mockUserResult, nil)

		svc := service.NewUserService(repo)
		user, err := svc.FindByID(context.TODO(), mockUserResult.ID)
		assert.NoError(t, err)
		assert.Equal(t, "john", user.FirstName)
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpexted Error"))

		svc := service.NewUserService(repo)
		user, err := svc.FindByID(context.TODO(), mockUserResult.ID)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_Create(t *testing.T) {
	req := reqres.CreateUserReq{
		FirstName: "john",
		LastName:  "due",
		Email:     "email",
	}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(nil)

		svc := service.NewUserService(repo)
		err := svc.Create(context.TODO(), &req)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(errors.New("Unexpexted Error"))

		svc := service.NewUserService(repo)
		err := svc.Create(context.TODO(), &req)
		assert.Error(t, err)
	})
}

func TestUserService_UpdateByID(t *testing.T) {
	req := reqres.UpdateUserReq{
		FirstName: "john",
		LastName:  "due",
		Email:     "email",
	}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("UpdateByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*domain.User")).
			Return(nil)

		svc := service.NewUserService(repo)
		err := svc.UpdateByID(context.TODO(), 1, &req)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("UpdateByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*domain.User")).
			Return(errors.New("Unexpexted Error"))

		svc := service.NewUserService(repo)
		err := svc.UpdateByID(context.TODO(), 1, &req)
		assert.Error(t, err)
	})
}

func TestUserService_DeleteByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("DeleteByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(nil)

		svc := service.NewUserService(repo)
		err := svc.DeleteByID(context.TODO(), 1)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		repo.On("DeleteByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(errors.New("Unexpexted Error"))

		svc := service.NewUserService(repo)
		err := svc.DeleteByID(context.TODO(), 1)
		assert.Error(t, err)
	})
}
