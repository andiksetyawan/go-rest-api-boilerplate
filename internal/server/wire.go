//go:build wireinject
// +build wireinject

package server

import (
	"database/sql"
	"net/http"

	"github.com/google/wire"
	httpTransport "go-rest-api-boilerplate/internal/transport/http"
	"go-rest-api-boilerplate/internal/usecase/repository"
	"go-rest-api-boilerplate/internal/usecase/service"
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
)

//var postSet = wire.NewSet(
//	repository.NewPostRepository,
//	service.NewPostService,
//)

func InitializedHandlerServer(db *sql.DB) http.Handler {
	wire.Build(
		userSet,
		httpTransport.NewHandler,
	)
	return nil
}
