package controller

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/romakorinenko/task-manager/internal/repository"
	mockRepository "github.com/romakorinenko/task-manager/internal/repository/mocks"
	"github.com/romakorinenko/task-manager/internal/service"
	"github.com/romakorinenko/task-manager/test"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTaskController_Create_UserCreated(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.POST("/tasks", taskController.Create)

	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	values := url.Values{}
	values.Set("Title", "Title")
	values.Set("Description", "DescriptionDescription")
	values.Set("UserLogin", "user")
	values.Set("Priority", "1")
	req.PostForm = values

	w := httptest.NewRecorder()

	userRepo.EXPECT().GetByLogin(gomock.Any(), gomock.Any()).Return(&repository.User{ID: 2, Login: "user"}, nil)
	taskRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(1, nil)

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusFound, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, "", string(respBodyBytes))
}

func TestTaskController_Create_InvalidPriority(t *testing.T) {
	router := test.SetUpTestRouter()

	taskService := service.NewTaskService(nil, nil)
	taskController := NewTaskController(taskService, nil)

	router.POST("/tasks", taskController.Create)

	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	values := url.Values{}
	values.Set("Title", "Title")
	values.Set("Description", "DescriptionDescription")
	values.Set("UserLogin", "user")
	values.Set("Priority", "r")
	req.PostForm = values

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `{"error":"priority is not a number"}`, string(respBodyBytes))
}

func TestTaskController_Create_InvalidTitle(t *testing.T) {
	router := test.SetUpTestRouter()

	taskService := service.NewTaskService(nil, nil)
	taskController := NewTaskController(taskService, nil)

	router.POST("/tasks", taskController.Create)

	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	values := url.Values{}
	values.Set("Title", "")
	values.Set("Description", "DescriptionDescription")
	values.Set("UserLogin", "user")
	values.Set("Priority", "1")
	req.PostForm = values

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `{"error":"bad request"}`, string(respBodyBytes))
}

func TestTaskController_Create_InternalErrorFromTaskRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.POST("/tasks", taskController.Create)

	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	values := url.Values{}
	values.Set("Title", "Title")
	values.Set("Description", "DescriptionDescription")
	values.Set("UserLogin", "user")
	values.Set("Priority", "1")
	req.PostForm = values

	w := httptest.NewRecorder()

	userRepo.EXPECT().GetByLogin(gomock.Any(), gomock.Any()).Return(&repository.User{ID: 2, Login: "user"}, nil)
	taskRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(0, errors.New("error"))

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, "{\"error\":\"internal server error\"}", string(respBodyBytes))
}

func TestTaskController_Update_TaskUpdated(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.POST("/tasks/:id", taskController.Update)

	req := httptest.NewRequest(http.MethodPost, "/tasks/1", nil)
	values := url.Values{}
	values.Set("Title", "Title")
	values.Set("Description", "DescriptionDescription")
	values.Set("Status", "OPEN")
	values.Set("Priority", "1")
	req.PostForm = values

	w := httptest.NewRecorder()

	taskFromDB := &repository.Task{
		ID:          1,
		Title:       "Title",
		Description: "Description",
		Status:      "OPEN",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      2,
	}
	taskRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(taskFromDB, nil)
	taskRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusFound, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, "", string(respBodyBytes))
}

func TestTaskController_Update_BadRequest(t *testing.T) {
	router := test.SetUpTestRouter()

	taskService := service.NewTaskService(nil, nil)
	taskController := NewTaskController(taskService, nil)

	router.POST("/tasks/:id", taskController.Update)

	req := httptest.NewRequest(http.MethodPost, "/tasks/1", nil)
	values := url.Values{}
	values.Set("Title", "")
	values.Set("Description", "DescriptionDescription")
	values.Set("Status", "OPEN")
	values.Set("Priority", "1")
	req.PostForm = values

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, "{\"error\":\"bad request\"}", string(respBodyBytes))
}

func TestTaskController_Update_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.POST("/tasks/:id", taskController.Update)

	req := httptest.NewRequest(http.MethodPost, "/tasks/1", nil)
	values := url.Values{}
	values.Set("Title", "Title")
	values.Set("Description", "DescriptionDescription")
	values.Set("Status", "OPEN")
	values.Set("Priority", "1")
	req.PostForm = values

	w := httptest.NewRecorder()

	taskFromDB := &repository.Task{
		ID:          1,
		Title:       "Title",
		Description: "Description",
		Status:      "OPEN",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      2,
	}
	taskRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(taskFromDB, nil)
	taskRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New(""))

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `{"error":"internal server error"}`, string(respBodyBytes))
}

func TestTaskController_Delete_UserDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.POST("/:id/delete", taskController.Delete)

	req := httptest.NewRequest(http.MethodPost, "/1/delete", nil)
	values := url.Values{}
	values.Set("Title", "Title")
	values.Set("Description", "DescriptionDescription")
	values.Set("Status", "OPEN")
	values.Set("Priority", "1")
	req.PostForm = values

	w := httptest.NewRecorder()

	taskRepo.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(nil)

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusFound, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, "", string(respBodyBytes))
}

func TestTaskController_Delete_BadRequest(t *testing.T) {
	router := test.SetUpTestRouter()

	taskService := service.NewTaskService(nil, nil)
	taskController := NewTaskController(taskService, nil)

	router.POST("/:id/delete", taskController.Delete)

	req := httptest.NewRequest(http.MethodPost, "/r/delete", nil)
	values := url.Values{}
	values.Set("Title", "Title")
	values.Set("Description", "DescriptionDescription")
	values.Set("Status", "OPEN")
	values.Set("Priority", "1")
	req.PostForm = values

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `{"error":"task ID is not number"}`, string(respBodyBytes))
}

func TestTaskController_Delete_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.POST("/:id/delete", taskController.Delete)

	req := httptest.NewRequest(http.MethodPost, "/1/delete", nil)
	values := url.Values{}
	values.Set("Title", "Title")
	values.Set("Description", "DescriptionDescription")
	values.Set("Status", "OPEN")
	values.Set("Priority", "1")
	req.PostForm = values

	w := httptest.NewRecorder()

	taskRepo.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(errors.New(""))

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	respBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `{"error":"internal server error"}`, string(respBodyBytes))
}

func TestTaskController_GetByID_TaskReturned(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := service.NewTaskService(taskRepo, nil)
	taskController := NewTaskController(taskService, nil)

	router.POST("/tasks/:id", taskController.GetByID)

	req := httptest.NewRequest(http.MethodPost, "/tasks/1", nil)

	w := httptest.NewRecorder()

	taskFromDB := &repository.TaskWithLogin{
		ID:          1,
		Title:       "Title",
		Description: "Description",
		Status:      "OPEN",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserLogin:   "user",
	}

	taskRepo.EXPECT().GetTaskWithLoginByID(gomock.Any(), gomock.Any()).Return(taskFromDB, nil)

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestTaskController_GetByID_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := service.NewTaskService(taskRepo, nil)
	taskController := NewTaskController(taskService, nil)

	router.POST("/tasks/:id", taskController.GetByID)

	req := httptest.NewRequest(http.MethodPost, "/tasks/1", nil)

	w := httptest.NewRecorder()

	taskRepo.EXPECT().GetTaskWithLoginByID(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestTaskController_GetByUserLogin_TasksReturned(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.GET("/tasks/user/:login", taskController.GetByUserLogin)

	req := httptest.NewRequest(http.MethodGet, "/tasks/user/user", nil)

	w := httptest.NewRecorder()

	taskRepo.EXPECT().GetByUserLogin(gomock.Any(), gomock.Any()).Return([]repository.Task{}, nil)

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestTaskController_GetByPriority_TaskReturned(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.GET("/tasks/by-priority/:priority", taskController.GetByPriority)

	req := httptest.NewRequest(http.MethodGet, "/tasks/by-priority/1", nil)

	w := httptest.NewRecorder()

	taskRepo.EXPECT().GetByPriority(gomock.Any(), gomock.Any()).Return([]repository.Task{}, nil)

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestTaskController_GetByPriority_BadRequest(t *testing.T) {
	router := test.SetUpTestRouter()

	taskService := service.NewTaskService(nil, nil)
	taskController := NewTaskController(taskService, nil)

	router.GET("/tasks/by-priority/:priority", taskController.GetByPriority)

	req := httptest.NewRequest(http.MethodGet, "/tasks/by-priority/r", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestTaskController_GetByStatus_TaskReturned(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.GET("/tasks/by-status/:status", taskController.GetByStatus)

	req := httptest.NewRequest(http.MethodGet, "/tasks/by-status/OPEN", nil)

	w := httptest.NewRecorder()

	taskRepo.EXPECT().GetByStatus(gomock.Any(), gomock.Any()).Return([]repository.Task{}, nil)

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestTaskController_GetByStatus_BadRequest(t *testing.T) {
	router := test.SetUpTestRouter()

	taskService := service.NewTaskService(nil, nil)
	taskController := NewTaskController(taskService, nil)

	router.GET("/tasks/by-status/:status", taskController.GetByStatus)

	req := httptest.NewRequest(http.MethodGet, "/tasks/by-status/1", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestTaskController_Edit_TemplateReturned(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.GET("/tasks/:id/edit", taskController.Edit)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1/edit", nil)

	w := httptest.NewRecorder()

	taskRepo.EXPECT().GetTaskWithLoginByID(gomock.Any(), gomock.Any()).Return(&repository.TaskWithLogin{}, nil)

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestTaskController_Edit_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	router := test.SetUpTestRouter()

	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)
	taskController := NewTaskController(taskService, userService)

	router.GET("/tasks/:id/edit", taskController.Edit)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1/edit", nil)

	w := httptest.NewRecorder()

	taskRepo.EXPECT().GetTaskWithLoginByID(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestTaskController_CreateTemplate_Unauthorized(t *testing.T) {
	router := test.SetUpTestRouter()

	taskService := service.NewTaskService(nil, nil)
	taskController := NewTaskController(taskService, nil)

	router.GET("/tasks/create", taskController.CreateTemplate)

	req := httptest.NewRequest(http.MethodGet, "/tasks/create", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestTaskController_GetAll(t *testing.T) {
	router := test.SetUpTestRouter()

	taskService := service.NewTaskService(nil, nil)
	taskController := NewTaskController(taskService, nil)

	router.GET("/tasks", taskController.CreateTemplate)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	require.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
