package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"go-rest-api-boilerplate/internal/domain"
)

func NewHandler(userService domain.UserService) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	NewUserHandlerRegister(r, userService)
	return r
}
