package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/romakorinenko/task-manager/internal/constant"
	"github.com/romakorinenko/task-manager/internal/dto"
	"github.com/romakorinenko/task-manager/internal/repository"
	"github.com/romakorinenko/task-manager/internal/service"
)

type ITaskController interface {
	GetByUserLogin(c *gin.Context)
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Edit(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Create(c *gin.Context)
	CreateTemplate(c *gin.Context)
	GetByStatus(c *gin.Context)
	GetByPriority(c *gin.Context)
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

// todo ok
func (t *TaskController) GetByUserLogin(c *gin.Context) {
	userLoginParam := c.Param("login")
	tasks, err := t.TaskService.GetByUserLogin(c.Request.Context(), userLoginParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// todo ok, but not fully
func (t *TaskController) GetAll(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userStr, ok := user.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var tasks []repository.Task
	userFromDB := t.UserService.GetByLogin(c.Request.Context(), userStr)
	if userFromDB.Role == constant.AdminRole {
		all, err := t.TaskService.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("")})
			return
		}
		tasks = all
	} else {
		byLogin, err := t.TaskService.GetByUserLogin(c.Request.Context(), userFromDB.Login)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("")})
			return
		}
		tasks = byLogin
	}

	var data []dto.TaskTemplateData
	for _, task := range tasks {
		userFromDB := t.UserService.GetByID(c.Request.Context(), task.UserID) // todo нужна новая структура и новый репо и сервис
		taskTemplateData := dto.TaskToTaskTemplateData(&task, userFromDB.Login)

		data = append(data, *taskTemplateData)
	}

	templateData := dto.TasksTemplateData{data}
	c.HTML(http.StatusOK, "usertasks.html", templateData)
}

// todo ok
func (t *TaskController) GetByID(c *gin.Context) {
	IDParam := c.Param("id")
	taskID, err := strconv.Atoi(IDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not number"})
		return
	}
	task, err := t.TaskService.GetByID(c.Request.Context(), taskID) // todo вынести в сервис
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not exists"})
		return
	}
	userFromDB := t.UserService.GetByID(c.Request.Context(), task.UserID)
	data := dto.TaskToTaskTemplateData(task, userFromDB.Login)

	c.HTML(http.StatusOK, "task.html", data)
}

// todo ok
func (t *TaskController) Edit(c *gin.Context) {
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
	data := dto.TaskToTaskTemplateData(task, userFromDB.Login)

	c.HTML(http.StatusOK, "taskedit.html", data)
}

// todo ok
func (t *TaskController) Update(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	taskForUpdate.Priority = priority
	taskForUpdate.Status = statusForm

	_, err = t.TaskService.Update(c.Request.Context(), taskForUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/tasks/%d", taskID))
}

// todo ok
func (t *TaskController) Delete(c *gin.Context) {
	taskIDParam := c.Param("id")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task ID is not number"})
		return
	}
	err = t.TaskService.DeleteByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Redirect(http.StatusFound, "/tasks")
}

// todo ok
func (t *TaskController) Create(c *gin.Context) {
	titleForm := c.PostForm("Title")
	descriptionForm := c.PostForm("Description")
	priorityForm := c.PostForm("Priority")
	priority, err := strconv.Atoi(priorityForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	userLoginForm := c.PostForm("UserLogin")
	userFromDB := t.UserService.GetByLogin(c.Request.Context(), userLoginForm)

	taskForCreate := &repository.Task{
		Title:       titleForm,
		Description: descriptionForm,
		Priority:    priority,
		UserID:      userFromDB.ID,
		Status:      constant.OpenTaskStatus,
	}
	createdTask, err := t.TaskService.Create(c.Request.Context(), taskForCreate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/tasks/%d", createdTask.ID))
}

// todo ok
func (t *TaskController) CreateTemplate(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userStr, ok := user.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var data dto.UsersTemplateData
	userFromDB := t.UserService.GetByLogin(c.Request.Context(), userStr)
	if userFromDB.Role == constant.AdminRole {
		users := t.UserService.GetAll(c.Request.Context())
		data.Users = users
	} else {
		users := []repository.User{*userFromDB}
		data.Users = users
	}

	c.HTML(http.StatusOK, "taskcreate.html", data)
}

func (t *TaskController) GetByStatus(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
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
	if userFromDB.Role != constant.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admins can watch tasks by status"})
		return
	}

	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is not signed"})
		return
	}
	tasks, err := t.TaskService.GetByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusOK, tasks)
}

func (t *TaskController) GetByPriority(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
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
	if userFromDB.Role != constant.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admins can watch tasks by status"})
	}

	priorityParam := c.Param("priority")
	if priorityParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "priorityParam is not signed"})
		return
	}
	priority, strConvErr := strconv.Atoi(priorityParam)
	if strConvErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if priority < 1 || priority > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "priorityParam should be from 1 to 4"})
		return
	}
	tasks, err := t.TaskService.GetByPriority(c.Request.Context(), priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
