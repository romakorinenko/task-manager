package controller

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/romakorinenko/task-manager/internal/constant"
	"github.com/romakorinenko/task-manager/internal/repository"
	"github.com/romakorinenko/task-manager/internal/service"
)

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

func (u *UserController) GetMainPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user != nil {
		c.Redirect(http.StatusFound, "/tasks")
		return
	}

	c.HTML(http.StatusOK, "login.html", nil)
}

// @Summary Login user // TODO
// @Description Log in a user and create a session
// @Accept json
// @Produce json
// @Param creds body Creds true "creds"
// @Success 200 {string} string "Logged in successfully"
// @Router /login [post]
func (u *UserController) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := u.UserService.GetByLogin(c.Request.Context(), username)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	}

	if username == user.Login && password == user.Password {
		session := sessions.Default(c)
		session.Set(constant.UserSessionKey, username)
		if err := session.Save(); err != nil {
			slog.Error("error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		c.Redirect(http.StatusFound, "/tasks")
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	}
}

// @Summary Logout user // TODO
// @Description Log out a user and destroy the session
// @Success 200 {string} string "Logged out successfully"
// @Router /logout [post]
func (u *UserController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(constant.UserSessionKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

// TODO
func (u *UserController) Create(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userLogin, ok := user.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}
	dbUser := u.UserService.GetByLogin(c.Request.Context(), userLogin)
	if dbUser.Role != constant.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admins can create users"})
		return
	}

	var newUser repository.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if u.UserService.GetByLogin(c.Request.Context(), newUser.Login) == nil {
		if err := u.UserService.Create(c.Request.Context(), &newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("user '%s' created", newUser.Login)})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("user '%s' already exists", newUser.Login)})
	}
}

// todo
func (u *UserController) Block(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userLogin, ok := user.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}
	dbUser := u.UserService.GetByLogin(c.Request.Context(), userLogin)
	if dbUser.Role != constant.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admins can block users"})
		return
	}

	userID := c.Param("id")

	if u.UserService.BlockByID(c.Request.Context(), userID) {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user '%s' blocked", userID)})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}
}

// todo
func (u *UserController) GetAll(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userLogin, ok := user.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}
	dbUser := u.UserService.GetByLogin(c.Request.Context(), userLogin)
	if dbUser.Role != constant.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"message": "only admins can receive all users"})
		return
	}

	users := u.UserService.GetAll(c.Request.Context())
	c.JSON(http.StatusOK, users)
}
