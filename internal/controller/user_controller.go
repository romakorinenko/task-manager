package controller

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/romakorinenko/task-manager/internal/constant"
	"github.com/romakorinenko/task-manager/internal/dto"
	"github.com/romakorinenko/task-manager/internal/errs"
	"github.com/romakorinenko/task-manager/internal/repository"
	"github.com/romakorinenko/task-manager/internal/service"
)

func init() {
	gob.Register(&repository.User{})
}

type IUserController interface {
	GetMainPage(c *gin.Context)

	Login(c *gin.Context)
	Logout(c *gin.Context)

	Create(c *gin.Context)
	Block(c *gin.Context)
	GetAll(c *gin.Context)
}

type UserController struct {
	UserService service.IUserService
	TaskService service.ITaskService
}

func NewUserController(userService service.IUserService, taskService service.ITaskService) *UserController {
	return &UserController{UserService: userService, TaskService: taskService}
}

// GetMainPage открывает главную страницу приложения.
// @Summary Get Main Page
// @Description открывает страницу для логина или главную таблицу с задачами, если пользователь уже авторизован
// @Tags pages
// @Produce html
// @Success 302 {string} Redirected to /tasks
// @Success 200 {object} dto.ResponseMap
// @Router / [get]
func (u *UserController) GetMainPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user != nil {
		c.Redirect(http.StatusFound, "/tasks")
		return
	}

	c.HTML(http.StatusOK, "login.html", nil)
}

// Login выполняет аутентификацию пользователя и создает сессию.
// @Summary User Login
// @Description аутентификация пользователя и создание сессии
// @Tags users
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 302 {object} dto.ResponseMap
// @Failure 401 {object} dto.ResponseMap
// @Failure 500 {object} dto.ResponseMap
// @Router /login [post]
func (u *UserController) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := u.UserService.GetByLogin(c.Request.Context(), username)
	if user == nil {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "bad request"})
	}

	if username == user.Login && password == user.Password {
		session := sessions.Default(c)
		session.Set(constant.UserSessionKey, user)
		if err := session.Save(); err != nil {
			slog.Error("error", err.Error())
			c.JSON(http.StatusInternalServerError, dto.ResponseMap{"message": "internal server error"})
			return
		}

		c.Redirect(http.StatusFound, "/tasks")
	} else {
		c.JSON(http.StatusUnauthorized, dto.ResponseMap{"error": "invalid credentials"})
	}
}

// Logout завершает сессию пользователя.
// @Summary User Logout
// @Description завершает сессию пользователя и открывает страницу для логина
// @Tags users
// @Produce json
// @Success 302 {string} Redirected to main page
// @Failure 500 {object} dto.ResponseMap
// @Router /logout [get]
func (u *UserController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(constant.UserSessionKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"message": "internal server error"})
	}
	c.Redirect(http.StatusFound, "/")
}

// Create создает нового пользователя.
// @Summary Create User
// @Description создает нового пользователя, только для администраторов
// @Tags users-admins
// @Accept json
// @Produce json
// @Param user body repository.User true "New User Data"
// @Success 201 {object} dto.ResponseMap
// @Failure 400 {object} dto.ResponseMap
// @Router /users [post]
func (u *UserController) Create(c *gin.Context) {
	var newUser repository.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": err.Error()})
		return
	}

	err := u.UserService.Create(c.Request.Context(), &newUser)
	if err == nil {
		c.JSON(http.StatusCreated, dto.ResponseMap{"message": fmt.Sprintf("user '%s' created", newUser.Login)})
	} else if err != nil && errors.Is(err, errs.UserExistsErr{}) {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": fmt.Sprintf("user '%s' already exists", newUser.Login)})
	} else {
		c.JSON(http.StatusInternalServerError, dto.ResponseMap{"error": err.Error()})
	}
}

// Block блокирует пользователя по идентификатору.
// @Summary Block User
// @Description блокирует пользователя по идентификатору, только для администраторов
// @Tags users-admins
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.ResponseMap
// @Failure 400 {object} dto.ResponseMap
// @Router /users/{id}/block [post]
func (u *UserController) Block(c *gin.Context) {
	userID := c.Param("id")

	if u.UserService.GetUserRepository().BlockByID(c.Request.Context(), userID) {
		c.JSON(http.StatusOK, dto.ResponseMap{"message": fmt.Sprintf("user '%s' blocked", userID)})
	} else {
		c.JSON(http.StatusBadRequest, dto.ResponseMap{"error": "incorrect user ID"})
	}
}

// GetAll возвращает список всех пользователей.
// @Summary Get All Users
// @Description возвращает список всех пользователей, только для администраторов
// @Tags users
// @Produce json
// @Success 200 {array} repository.User "List of users"
// @Failure 500 {object} dto.ResponseMap
// @Router /users [get]
func (u *UserController) GetAll(c *gin.Context) {
	users := u.UserService.GetAll(c.Request.Context())
	c.JSON(http.StatusOK, users)
}
