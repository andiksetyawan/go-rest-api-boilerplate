package http

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-rest-api-boilerplate/internal/domain"
	"go-rest-api-boilerplate/internal/model/reqres"
	"go-rest-api-boilerplate/pkg/httputil"
	"go-rest-api-boilerplate/pkg/util"
)

type userHandler struct {
	userSvc domain.UserService
}

func NewUserHandlerRegister(r *mux.Router, service domain.UserService) {
	handler := userHandler{userSvc: service}
	v1 := r.PathPrefix("/api/v1").Subrouter()
	{
		v1.HandleFunc("/user", handler.Create).Methods(http.MethodPost)
		v1.HandleFunc("/user", handler.FindAll).Methods(http.MethodGet)
		v1.HandleFunc("/user/{id}", handler.FindByID).Methods(http.MethodGet)
		v1.HandleFunc("/user/{id}", handler.DeleteByID).Methods(http.MethodDelete)
		v1.HandleFunc("/user/{id}", handler.UpdateByID).Methods(http.MethodPatch)
	}
}

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUserReq reqres.CreateUserReq
	err := json.NewDecoder(r.Body).Decode(&createUserReq)
	if err != nil {
		log.WithContext(r.Context()).WithError(err).Warn("error decoding json payload")
		httputil.RespondWithError(w, http.StatusUnprocessableEntity, "cannot receive the payload schema")
		return
	}

	err = validator.New().Struct(&createUserReq)
	if err != nil {
		log.WithContext(r.Context()).WithError(err).Warn("error validator")
		httputil.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.userSvc.Create(r.Context(), &createUserReq)
	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, httputil.ApiResponse{
		Error:   false,
		Message: "OK",
		Data:    nil,
	})
}

func (h *userHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := util.StringToInt64(mux.Vars(r)["id"])
	if err != nil {
		log.WithContext(r.Context()).WithError(err).Warn("param id not valid")
		httputil.RespondWithError(w, http.StatusUnprocessableEntity, "id not valid")
		return
	}

	var updateUserReq reqres.UpdateUserReq
	err = json.NewDecoder(r.Body).Decode(&updateUserReq)
	if err != nil {
		log.WithContext(r.Context()).WithError(err).Warn("error decoding json payload")
		httputil.RespondWithError(w, http.StatusUnprocessableEntity, "cannot receive the payload schema")
		return
	}

	err = validator.New().Struct(&updateUserReq)
	if err != nil {
		log.WithContext(r.Context()).WithError(err).Warn("error validator")
		httputil.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.userSvc.UpdateByID(r.Context(), id, &updateUserReq)
	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, httputil.ApiResponse{
		Error:   false,
		Message: "OK",
		Data:    nil,
	})
}

func (h *userHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := util.StringToInt64(mux.Vars(r)["id"])
	if err != nil {
		log.WithContext(r.Context()).WithError(err).Warn("param id not valid")
		httputil.RespondWithError(w, http.StatusUnprocessableEntity, "id not valid")
		return
	}

	err = h.userSvc.DeleteByID(r.Context(), id)
	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, httputil.ApiResponse{
		Error:   false,
		Message: "OK",
		Data:    nil,
	})
}

func (h *userHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userSvc.FindAll(r.Context())
	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, httputil.ApiResponse{
		Error:   false,
		Message: "OK",
		Data:    users,
	})
}

func (h *userHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := util.StringToInt64(mux.Vars(r)["id"])
	if err != nil {
		log.WithContext(r.Context()).WithError(err).Warn("param id not valid")
		httputil.RespondWithError(w, http.StatusUnprocessableEntity, "id not valid")
		return
	}

	user, err := h.userSvc.FindByID(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
		}

		httputil.RespondWithError(w, status, "")
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, httputil.ApiResponse{
		Error:   false,
		Message: "OK",
		Data:    user,
	})
}
