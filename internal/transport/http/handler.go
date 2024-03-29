package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"go-rest-api-boilerplate/config"
	"go-rest-api-boilerplate/internal/domain"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func NewHandler(userService domain.UserService) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Use(otelmux.Middleware(config.App.ServiceName))

	//Registered handler
	NewUserHandlerRegister(r, userService)

	return r
}
