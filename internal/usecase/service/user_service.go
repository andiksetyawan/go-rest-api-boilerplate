package service

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go-rest-api-boilerplate/config"
	"go-rest-api-boilerplate/internal/domain"
	"go-rest-api-boilerplate/internal/model/reqres"
	"go.opentelemetry.io/otel"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (u *userService) Create(ctx context.Context, req *reqres.CreateUserReq) error {
	ctx, span := otel.Tracer(config.App.ServiceName).Start(ctx, "user.service.Create")
	defer span.End()

	newUser := domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := u.repo.Save(ctx, &newUser)
	if err != nil {
		return err
	}

	log.Debugf("last insertid :%v", newUser.ID)
	return nil
}

func (u *userService) UpdateByID(ctx context.Context, id int64, req *reqres.UpdateUserReq) error {
	ctx, span := otel.Tracer(config.App.ServiceName).Start(ctx, "user.service.UpdateByID")
	defer span.End()

	newUser := domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	return u.repo.UpdateByID(ctx, id, &newUser)
}

func (u *userService) DeleteByID(ctx context.Context, id int64) error {
	ctx, span := otel.Tracer(config.App.ServiceName).Start(ctx, "user.service.DeleteByID")
	defer span.End()

	return u.repo.DeleteByID(ctx, id)
}

func (u *userService) FindAll(ctx context.Context) (*[]domain.User, error) {
	ctx, span := otel.Tracer(config.App.ServiceName).Start(ctx, "user.service.FindAll")
	defer span.End()

	return u.repo.FindAll(ctx)
}

func (u *userService) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	ctx, span := otel.Tracer(config.App.ServiceName).Start(ctx, "user.service.FindById")
	defer span.End()

	return u.repo.FindByID(ctx, id)
}
