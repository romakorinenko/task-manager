package controller

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/romakorinenko/task-manager/internal/repository"
	"github.com/romakorinenko/task-manager/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type TaskTemplateData struct {
	ID          int
	Title       string
	Description string
	Priority    int
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserLogin   string
}

func TaskToTaskTemplateData(task *repository.Task, userLogin string) *TaskTemplateData {
	return &TaskTemplateData{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		UserLogin:   userLogin,
	}
}

type UsersTemplateData struct {
	Users []repository.User
}

type ITaskController interface {
	GetByUserLogin(c *gin.Context)
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Edit(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Create(c *gin.Context)
	CreateTemplate(c *gin.Context)
}

type TaskController struct {
	TaskService service.ITaskService
	UserService service.IUserService
}

func NewTaskController(taskService service.ITaskService, userService service.IUserService) *TaskController {
	return &TaskController{
		TaskService: taskService,
		UserService: userService,
	}
}

func (t *TaskController) GetByUserLogin(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userLoginParam := c.Param("login")
	tasks, err := t.TaskService.GetByUserLogin(c.Request.Context(), userLoginParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (t *TaskController) GetAll(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userStr, ok := user.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	userFromDB := t.UserService.GetByLogin(c.Request.Context(), userStr)
	if userFromDB.Role != "ADMIN" { // todo в константу
		slog.Error(fmt.Sprintf("cannot receive all users cause user '%s' is not ADMIN", userStr))
	}

	allTasks, err := t.TaskService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("")})
		return
	}

	var data []TaskTemplateData
	for _, task := range allTasks {
		userFromDB := t.UserService.GetByID(c.Request.Context(), task.UserID) // todo нужна новая структура и новый репо и сервис
		taskTemplateData := TaskToTaskTemplateData(&task, userFromDB.Login)

		data = append(data, *taskTemplateData)
	}

	templateData := TasksTemplateData{data}
	c.HTML(http.StatusOK, "usertasks.html", templateData)
}

func (t *TaskController) GetByID(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	IDParam := c.Param("id")
	taskID, err := strconv.Atoi(IDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not number"})
		return
	}
	task, err := t.TaskService.GetByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not exists"})
		return
	}
	userFromDB := t.UserService.GetByID(c.Request.Context(), task.UserID)
	data := TaskToTaskTemplateData(task, userFromDB.Login)

	c.HTML(http.StatusOK, "task.html", data)
}

func (t *TaskController) Edit(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// todo либо админ, либо пользователь

	taskIDParam := c.Param("id")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not number"})
		return
	}
	task, err := t.TaskService.GetByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not exists"})
		return
	}
	userFromDB := t.UserService.GetByID(c.Request.Context(), task.UserID)
	data := TaskToTaskTemplateData(task, userFromDB.Login)

	c.HTML(http.StatusOK, "taskedit.html", data)
}

func (t *TaskController) Update(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	titleForm := c.PostForm("Title")
	descriptionForm := c.PostForm("Description")
	priorityForm := c.PostForm("Priority")
	statusForm := c.PostForm("Status")
	taskIDParam := c.Param("id")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not number"})
		return
	}
	taskForUpdate, err := t.TaskService.GetByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not exists"})
		return
	}
	taskForUpdate.Title = titleForm
	taskForUpdate.Description = descriptionForm
	priority, err := strconv.Atoi(priorityForm)
	if err != nil {
		// todo
	}
	taskForUpdate.Priority = priority
	taskForUpdate.Status = statusForm

	_, err = t.TaskService.Update(c.Request.Context(), taskForUpdate)
	if err != nil {
		// todo
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/tasks/%d", taskID))
}

func (t *TaskController) Delete(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskIDParam := c.Param("id")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not number"})
		return
	}
	err = t.TaskService.DeleteByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.Redirect(http.StatusFound, "/tasks")
}

func (t *TaskController) Create(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	titleForm := c.PostForm("Title")
	descriptionForm := c.PostForm("Description")
	priorityForm := c.PostForm("Priority")
	priority, err := strconv.Atoi(priorityForm)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userLoginForm := c.PostForm("UserLogin")
	userFromDB := t.UserService.GetByLogin(c.Request.Context(), userLoginForm)

	taskForCreate := &repository.Task{
		Title:       titleForm,
		Description: descriptionForm,
		Priority:    priority,
		UserID:      userFromDB.ID,
		Status:      "OPEN", // todo в константы
	}
	createdTask, err := t.TaskService.Create(c.Request.Context(), taskForCreate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/tasks/%d", createdTask.ID))
}

func (t *TaskController) CreateTemplate(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userStr, ok := user.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var data UsersTemplateData
	userFromDB := t.UserService.GetByLogin(c.Request.Context(), userStr)
	if userFromDB.Role == "ADMIN" { //todo  в константу
		users := t.UserService.GetAll(c.Request.Context())
		data.Users = users
	} else {
		users := []repository.User{*userFromDB}
		data.Users = users
	}

	c.HTML(http.StatusOK, "taskcreate.html", data)
}
