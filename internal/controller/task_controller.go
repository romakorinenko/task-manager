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

// GetByUserLogin возвращает задачи пользователя по его логину.
// @Summary Get Tasks by User Login
// @Description возвращает задачи пользователя по его логину
// @Tags tasks
// @Produce json
// @Param login path string true "User Login"
// @Success 200 {array} repository.Task "List of tasks"
// @Failure 500 {object} dto.ResponseMap
// @Router /tasks/user/{login} [get]
func (t *TaskController) GetByUserLogin(c *gin.Context) {
	userLoginParam := c.Param("login")
	tasks, err := t.TaskService.GetByUserLogin(c.Request.Context(), userLoginParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetAll возвращает список задач в зависимости от роли пользователя.
// @Summary Get All Tasks
// @Description возвращает список задач в зависимости от роли:
// @Description для админов - все задачи, для пользователей - задачи пользователя
// @Tags tasks
// @Produce json
// @Success 200 {array} repository.Task "List of tasks"
// @Failure 401 {object} dto.ResponseMap
// @Failure 500 {object} dto.ResponseMap
// @Failure 400 {object} dto.ResponseMap
// @Router /tasks [get]
func (t *TaskController) GetAll(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.ResponseMap{"error": "unauthorized"})
		return
	}

	sessionUser, ok := user.(*repository.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": "internal server error"})
		return
	}

	var tasks []repository.Task
	if sessionUser.Role == constant.AdminRole {
		all, err := t.TaskService.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": fmt.Sprintf("")})
			return
		}
		tasks = all
	} else if sessionUser.Role == constant.UserRole {
		byLogin, err := t.TaskService.GetByUserLogin(c.Request.Context(), sessionUser.Login)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": fmt.Sprintf("")})
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

// GetByID возвращает задачу по идентификатору.
// @Summary Get Task by ID
// @Description возвращает задачу по идентификатору
// @Tags tasks
// @Produce html
// @Param id path string true "Task ID"
// @Success 200 {object} dto.TaskTemplateData
// @Failure 400 {object} dto.ResponseMap
// @Router /tasks/{id} [get]
func (t *TaskController) GetByID(c *gin.Context) {
	IDParam := c.Param("id")
	taskID, err := strconv.Atoi(IDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "task ID is not number"})
		return
	}
	task, err := t.TaskService.GetByID(c.Request.Context(), taskID) // todo вынести в сервис, новый метод репо создать
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "task ID is not exists"})
		return
	}
	userFromDB := t.UserService.GetByID(c.Request.Context(), task.UserID)
	data := dto.TaskToTaskTemplateData(task, userFromDB.Login)

	c.HTML(http.StatusOK, "task.html", data)
}

// Edit отображает форму редактирования задачи по идентификатору.
// @Description отображает форму редактирования задачи по идентификатору
// @Tags pages
// @Produce html
// @Param id path string true "Task ID"
// @Success 200 {object} dto.TaskTemplateData
// @Failure 400 {object} dto.ResponseMap
// @Router /tasks/{id}/edit [get]
func (t *TaskController) Edit(c *gin.Context) {
	taskIDParam := c.Param("id")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "task ID is not number"})
		return
	}
	task, err := t.TaskService.GetByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "task ID is not exists"})
		return
	}
	userFromDB := t.UserService.GetByID(c.Request.Context(), task.UserID)
	data := dto.TaskToTaskTemplateData(task, userFromDB.Login)

	c.HTML(http.StatusOK, "taskedit.html", data)
}

// Update обновляет задачу по идентификатору.
// @Summary Update Task by ID
// @Description обновляет задачу по идентификатору
// @Tags tasks
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path string true "Task ID"
// @Param Title formData string true "Task Title"
// @Param Description formData string true "Task Description"
// @Param Priority formData integer true "Task Priority"
// @Param Status formData string true "Task Status"
// @Success 302 {string} Redirected to updated task
// @Failure 400 {object} dto.ResponseMap
// @Failure 500 {object} dto.ResponseMap
// @Router /tasks/{id} [post]
func (t *TaskController) Update(c *gin.Context) {
	titleForm := c.PostForm("Title")
	descriptionForm := c.PostForm("Description")
	priorityForm := c.PostForm("Priority")
	statusForm := c.PostForm("Status")
	taskIDParam := c.Param("id")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "task ID is not number"})
		return
	}
	taskForUpdate, err := t.TaskService.GetByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "task ID is not exists"})
		return
	}
	taskForUpdate.Title = titleForm
	taskForUpdate.Description = descriptionForm
	priority, err := strconv.Atoi(priorityForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": "internal server error"})
		return
	}
	taskForUpdate.Priority = priority
	taskForUpdate.Status = statusForm

	_, err = t.TaskService.Update(c.Request.Context(), taskForUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": "internal server error"})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/tasks/%d", taskID))
}

// Delete удаляет задачу по указанному идентификатору.
// @Summary Delete Task by ID
// @Description Удаляет задачу по указанному идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 302 {object} dto.ResponseMap
// @Failure 400 {object} dto.ResponseMap
// @Failure 500 {object} dto.ResponseMap
// @Router /tasks/{id}/delete [post]
func (t *TaskController) Delete(c *gin.Context) {
	taskIDParam := c.Param("id")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "task ID is not number"})
		return
	}
	err = t.TaskService.DeleteByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": "internal server error"})
		return
	}

	c.Redirect(http.StatusFound, "/tasks")
}

// Create создаёт новую задачу.
// @Summary Create task
// @Description Создаёт новую задачу с указанными параметрами: заголовком, описанием, приоритетом и пользователем.
// @Tags tasks
// @Accept json
// @Produce json
// @Param Title formData string true "Заголовок задачи"
// @Param Description formData string true "Описание задачи"
// @Param Priority formData int true "Приоритет задачи (число)"
// @Param UserLogin formData string true "Логин пользователя, которому назначена задача"
// @Success 302 {object} dto.ResponseMap
// @Failure 400 {object} dto.ResponseMap
// @Failure 500 {object} dto.ResponseMap
// @Router /tasks/create [post]
func (t *TaskController) Create(c *gin.Context) {
	titleForm := c.PostForm("Title")
	descriptionForm := c.PostForm("Description")
	priorityForm := c.PostForm("Priority")
	priority, err := strconv.Atoi(priorityForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": "internal server error"})
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
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/tasks/%d", createdTask.ID))
}

// CreateTemplate отображает форму создания задачи.
// @Description Отображает страницу создания шаблона задачи для пользователя в зависимости от его роли.
// @Tags pages
// @Accept json
// @Produce html
// @Success 200 {object} dto.UsersTemplateData
// @Failure 401 {object} dto.ResponseMap
// @Failure 500 {object} dto.ResponseMap
// @Router /tasks/create [get]
func (t *TaskController) CreateTemplate(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.ResponseMap{"error": "unauthorized"})
		return
	}

	sessionUser, ok := user.(*repository.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": "internal server error"})
		return
	}

	var data dto.UsersTemplateData
	if sessionUser.Role == constant.AdminRole {
		users := t.UserService.GetAll(c.Request.Context())
		data.Users = users
	} else if sessionUser.Role == constant.UserRole {
		users := []repository.User{*sessionUser}
		data.Users = users
	}

	c.HTML(http.StatusOK, "taskcreate.html", data)
}

// GetByStatus Получить задачи по статусу.
// @Summary Get Tasks by status
// @Description Возвращает список задач с указанным статусом. Статус должен быть одним из: OPEN, IN_PROGRESS или DONE.
// @Description Только для админов
// @Tags tasks-admins
// @Accept json
// @Produce json
// @Param status path string true "Статус задачи"
// @Success 200 {array} repository.Task
// @Failure 400 {object} dto.ResponseMap
// @Failure 500 {object} dto.ResponseMap
// @Router /tasks/by-status/{status} [get]
func (t *TaskController) GetByStatus(c *gin.Context) {
	// todo вынести в сервис проверку

	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "status is not signed"})
		return
	}

	var isStatus bool
	for _, taskStatus := range constant.TaskStatuses {
		if taskStatus == status {
			isStatus = true
		}
	}

	if !isStatus {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": fmt.Sprintf("status %s incorrect, use OPEN, IN_PROGRESS or DONE", status)})
		return
	}

	tasks, err := t.TaskService.GetByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": "internal server error"})
	}

	c.JSON(http.StatusOK, tasks)
}

// GetByPriority Получить задачи по приоритету.
// @Summary Get Tasks by priority
// @Description Возвращает список задач с указанным приоритетом. Только для админов
// @Tags tasks-admins
// @Accept json
// @Produce json
// @Param priority path int true "Приоритет задачи"
// @Success 200 {array} repository.Task
// @Failure 400 {object} dto.ResponseMap
// @Failure 500 {object} dto.ResponseMap
// @Router /tasks/by-priority/{priority} [get]
func (t *TaskController) GetByPriority(c *gin.Context) {
	// todo вынести в сервис проверку

	priorityParam := c.Param("priority")
	if priorityParam == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "priority is not signed or incorrect"})
		return
	}
	priority, strConvErr := strconv.Atoi(priorityParam)
	if strConvErr != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": "internal server error"})
		return
	}

	if priority < 1 || priority > 4 {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "priority should be from 1 to 4"})
		return
	}
	tasks, err := t.TaskService.GetByPriority(c.Request.Context(), priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
