package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/romakorinenko/task-manager/internal/constant"
	"github.com/romakorinenko/task-manager/internal/repository"
	"github.com/romakorinenko/task-manager/internal/service"
	"github.com/romakorinenko/task-manager/test"
	"github.com/stretchr/testify/require"
)

func TestUserController_Login(t *testing.T) {
	ctx := context.Background()
	DB := test.CreateDBForTest(t, "/cmd/migrations")
	defer DB.Close()

	userRepo := repository.NewUserRepo(DB.DBPool)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	userRepo.Create(
		ctx,
		&repository.User{
			ID:        1,
			Login:     "drug",
			CreatedAt: time.Now(),
			Role:      constant.UserRole,
			Password:  "drug",
			Active:    true,
		},
	)

	r := test.SetUpTestRouter()
	r.POST("/login", userController.Login)

	t.Run("успешное создание сессии", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		data := url.Values{}
		data.Set("username", "drug")
		data.Set("password", "drug")
		req.PostForm = data

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := w.Result()
		require.Equal(t, http.StatusFound, response.StatusCode)
		respBodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t, "", string(respBodyBytes))
	})
	t.Run("пользователь не авторизован", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		data := url.Values{}
		data.Set("username", "wrong")
		data.Set("password", "wrong")
		req.PostForm = data

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := w.Result()
		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
		respBodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t, "{\"error\":\"invalid credentials\"}", string(respBodyBytes))
	})
}

func TestUserController_Logout(t *testing.T) {
	ctx := context.Background()
	DB := test.CreateDBForTest(t, "/cmd/migrations")
	defer DB.Close()

	userRepo := repository.NewUserRepo(DB.DBPool)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	userRepo.Create(
		ctx,
		&repository.User{
			ID:        1,
			Login:     "drug",
			CreatedAt: time.Now(),
			Role:      constant.UserRole,
			Password:  "drug",
			Active:    true,
		},
	)

	r := test.SetUpTestRouter()
	r.POST("/logout", userController.Logout)

	t.Run("пользователь разлогинен, сессия удалена", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := w.Result()
		require.Equal(t, http.StatusFound, response.StatusCode)
		respBodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t, "", string(respBodyBytes))
	})
}

func TestUserController_GetMainPage(t *testing.T) {
	ctx := context.Background()
	DB := test.CreateDBForTest(t, "/cmd/migrations")
	defer DB.Close()

	userRepo := repository.NewUserRepo(DB.DBPool)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	userRepo.Create(
		ctx,
		&repository.User{
			ID:        1,
			Login:     "drug",
			CreatedAt: time.Now(),
			Role:      constant.UserRole,
			Password:  "drug",
			Active:    true,
		},
	)

	r := test.SetUpTestRouter()
	r.GET("/", userController.GetMainPage)

	t.Run("откроет страницу логина", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := w.Result()
		require.Equal(t, http.StatusOK, response.StatusCode)
		respBodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.True(t, strings.Contains(string(respBodyBytes), "<title>Вход в Task Manager</title>"))
	})
}

func TestUserController_Create(t *testing.T) {
	ctx := context.Background()
	DB := test.CreateDBForTest(t, "/cmd/migrations")
	defer DB.Close()

	userRepo := repository.NewUserRepo(DB.DBPool)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	user := userRepo.Create(
		ctx,
		&repository.User{
			ID:        1,
			Login:     "drug",
			CreatedAt: time.Now(),
			Role:      constant.UserRole,
			Password:  "drug",
			Active:    true,
		},
	)

	r := test.SetUpTestRouter()
	r.POST("/users", userController.Create)

	t.Run("получен невалидный JSON в теле", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("[]"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := w.Result()
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		respBodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t,
			`{"error":"json: cannot unmarshal array into Go value of type repository.User"}`,
			string(respBodyBytes),
		)
	})
	t.Run("пользователь с таким логином уже существует", func(t *testing.T) {
		userBytes, err := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userBytes))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := w.Result()
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		respBodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t,
			`{"error":"user 'drug' already exists"}`,
			string(respBodyBytes),
		)
	})
	t.Run("пользователь создан", func(t *testing.T) {
		userBytes, err := json.Marshal(
			&repository.User{
				ID:        1,
				Login:     "user",
				CreatedAt: time.Now(),
				Role:      constant.UserRole,
				Password:  "user",
				Active:    true,
			},
		)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userBytes))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := w.Result()
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		respBodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t,
			`{"error":"user 'user' already exists"}`,
			string(respBodyBytes),
		)
	})
}

func TestUserController_Block(t *testing.T) {
	ctx := context.Background()
	DB := test.CreateDBForTest(t, "/cmd/migrations")
	defer DB.Close()

	userRepo := repository.NewUserRepo(DB.DBPool)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	user := userRepo.Create(
		ctx,
		&repository.User{
			ID:        1,
			Login:     "drug",
			CreatedAt: time.Now(),
			Role:      constant.UserRole,
			Password:  "drug",
			Active:    true,
		},
	)

	r := test.SetUpTestRouter()
	r.PUT("/users/:id/block", userController.Block)

	t.Run("пользователь заблокирован", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%d/block", user.ID), nil)
		req.SetPathValue("id", strconv.Itoa(user.ID))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := w.Result()
		require.Equal(t, http.StatusOK, response.StatusCode)
		respBodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t,
			`{"message":"user '100' blocked"}`,
			string(respBodyBytes),
		)
	})
	t.Run("пользователь не заблокирован, поскольку отсутствует в БД", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/{id}/block", nil)
		req.SetPathValue("id", strconv.Itoa(1000))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		response := w.Result()
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		respBodyBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t,
			`{"error":"incorrect user ID"}`,
			string(respBodyBytes),
		)
	})
}

// TestUserController_GetAll 2 юзера создаются в миграциях
func TestUserController_GetAll(t *testing.T) {
	expectedUsersSize := 2

	DB := test.CreateDBForTest(t, "/cmd/migrations")
	defer DB.Close()

	userRepo := repository.NewUserRepo(DB.DBPool)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	r := test.SetUpTestRouter()
	r.GET("/users", userController.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusOK, response.StatusCode)

	var users []*repository.User
	err := json.NewDecoder(response.Body).Decode(&users)
	require.NoError(t, err)
	require.True(t, len(users) == expectedUsersSize)
}
