package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-rest-api-boilerplate/internal/domain"
	"go-rest-api-boilerplate/internal/domain/mocks"
	"go-rest-api-boilerplate/internal/model/reqres"
	"go-rest-api-boilerplate/pkg/httputil"
)

func TestUserHandler_FindAll(t *testing.T) {
	mockUsers := []domain.User{
		{
			ID:        1,
			FirstName: "john",
			LastName:  "due",
			Email:     "john@email.local",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}

	t.Run("success", func(t *testing.T) {
		//r := mux.NewRouter()
		mockUserSvc := mocks.NewUserService(t)
		//NewUserHandlerRegister(r, mockUserSvc)
		mockUserSvc.On("FindAll", mock.Anything).Return(&mockUsers, nil)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user", strings.NewReader(""))
		assert.NoError(t, err)

		handler := userHandler{userSvc: mockUserSvc}
		handler.FindAll(w, req)
		//r.ServeHTTP(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "OK", response.Message)

		b, err := json.Marshal(response.Data)
		assert.NoError(t, err)

		var responseUsers []domain.User
		err = json.Unmarshal(b, &responseUsers)
		assert.NoError(t, err)

		assert.Equal(t, len(mockUsers), len(responseUsers))
		assert.Equal(t, mockUsers[0].FirstName, responseUsers[0].FirstName)
	})

	t.Run("error", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)
		mockUserSvc.On("FindAll", mock.Anything).Return(nil, errors.New("Unexpexted Error"))

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user", strings.NewReader(""))
		assert.NoError(t, err)

		handler := userHandler{userSvc: mockUserSvc}
		handler.FindAll(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, true, response.Error)
		assert.Nil(t, response.Data)
	})
}

func TestUserHandler_FindByID(t *testing.T) {
	mockUser := domain.User{
		ID:        1,
		FirstName: "john",
		LastName:  "due",
		Email:     "john@email.local",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	t.Run("success", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)
		mockUserSvc.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockUser, nil)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user/1", strings.NewReader(""))
		assert.NoError(t, err)

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler := userHandler{userSvc: mockUserSvc}
		handler.FindByID(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "OK", response.Message)

		b, err := json.Marshal(response.Data)
		assert.NoError(t, err)

		var resUser domain.User
		err = json.Unmarshal(b, &resUser)
		assert.NoError(t, err)

		assert.Equal(t, mockUser.ID, resUser.ID)
		assert.Equal(t, mockUser.FirstName, resUser.FirstName)
	})

	t.Run("error", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)
		mockUserSvc.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpexted Error"))

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/v1/user/1", strings.NewReader(""))
		assert.NoError(t, err)

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler := userHandler{userSvc: mockUserSvc}
		handler.FindByID(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, true, response.Error)
		assert.Nil(t, response.Data)
	})
}

func TestUserHandler_Create(t *testing.T) {
	userReq := reqres.CreateUserReq{
		FirstName: "john",
		LastName:  "due",
		Email:     "john@m.co",
	}

	b, err := json.Marshal(userReq)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)
		mockUserSvc.On("Create", mock.Anything, mock.AnythingOfType("*reqres.CreateUserReq")).Return(nil)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/user", strings.NewReader(string(b)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		handler := userHandler{userSvc: mockUserSvc}
		handler.Create(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "OK", response.Message)
	})

	t.Run("error:service", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)
		mockUserSvc.On("Create", mock.Anything, mock.AnythingOfType("*reqres.CreateUserReq")).Return(errors.New("Unexpexted Error"))

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/user", strings.NewReader(string(b)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		handler := userHandler{userSvc: mockUserSvc}
		handler.Create(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error:validator", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/v1/user", strings.NewReader("{\"first_name\":\"\"}"))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		handler := userHandler{userSvc: mockUserSvc}
		handler.Create(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}

func TestUserHandler_UpdateByID(t *testing.T) {
	userReq := reqres.UpdateUserReq{
		FirstName: "john",
		LastName:  "due",
		Email:     "john@m.co",
	}

	b, err := json.Marshal(userReq)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)
		mockUserSvc.On("UpdateByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*reqres.UpdateUserReq")).
			Return(nil)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/api/v1/user", strings.NewReader(string(b)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler := userHandler{userSvc: mockUserSvc}
		handler.UpdateByID(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "OK", response.Message)
	})

	t.Run("error:service", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)
		mockUserSvc.On("UpdateByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*reqres.UpdateUserReq")).
			Return(errors.New("Unexpexted Error"))

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/api/v1/user", strings.NewReader(string(b)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler := userHandler{userSvc: mockUserSvc}
		handler.UpdateByID(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error:validator", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/api/v1/user", strings.NewReader("{\"first_name\":\"\"}"))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler := userHandler{userSvc: mockUserSvc}
		handler.UpdateByID(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}

func TestUserHandler_DeleteByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)
		mockUserSvc.On("DeleteByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/1", strings.NewReader(""))
		assert.NoError(t, err)

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler := userHandler{userSvc: mockUserSvc}
		handler.DeleteByID(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "OK", response.Message)
	})

	t.Run("error", func(t *testing.T) {
		mockUserSvc := mocks.NewUserService(t)
		mockUserSvc.On("DeleteByID", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("Unexpexted Error"))

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/1", strings.NewReader(""))
		assert.NoError(t, err)

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler := userHandler{userSvc: mockUserSvc}
		handler.DeleteByID(w, req)

		var response httputil.ApiResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, true, response.Error)
		assert.Nil(t, response.Data)
	})
}
