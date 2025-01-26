package server

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/romakorinenko/task-manager/internal/constant"
	"github.com/romakorinenko/task-manager/internal/controller"
	"github.com/romakorinenko/task-manager/internal/repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed templates/*
var templates embed.FS

var Router *gin.Engine

func RegisterServerAndHandlers(
	userController controller.IUserController,
	taskController controller.ITaskController,
	port int,
) {
	Router = gin.Default()
	tmpl := template.Must(template.ParseFS(templates, "templates/*.html"))
	Router.SetHTMLTemplate(tmpl)
	store := sessions.NewCookieStore([]byte("secret"))
	Router.Use(sessions.Sessions("sessions", store))

	RegisterUserHandlers(userController)
	RegisterTaskHandlers(taskController)
	RegisterSwaggerAndMetricsHandlers()

	slog.Info("Swagger available on http://localhost:8080/swagger/index.html")
	err := Router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("server closed", slog.Any("error", err))
	}
}

func RegisterUserHandlers(userController controller.IUserController) {
	Router.GET("/", userController.GetMainPage)
	Router.POST("/login", userController.Login)
	Router.GET("/logout", userController.Logout)

	usersRouterGroup := Router.Group("/users")
	{
		usersRouterGroup.POST("", AdminSessionMiddleware, userController.Create)
		usersRouterGroup.PUT("/:id/block", AdminSessionMiddleware, userController.Block)
		usersRouterGroup.GET("", AdminSessionMiddleware, userController.GetAll)
	}
}

func RegisterTaskHandlers(taskController controller.ITaskController) {
	tasksRouterGroup := Router.Group("/tasks")
	{
		tasksRouterGroup.GET("/create", taskController.CreateTemplate)
		tasksRouterGroup.POST("", UserSessionMiddleware, taskController.Create)
		tasksRouterGroup.POST("/:id", UserSessionMiddleware, taskController.Update)
		tasksRouterGroup.POST("/:id/delete", UserSessionMiddleware, taskController.Delete)
		tasksRouterGroup.GET("/:id", UserSessionMiddleware, taskController.GetByID)
		tasksRouterGroup.GET("/:id/edit", UserSessionMiddleware, taskController.Edit)
		tasksRouterGroup.GET("/user/:login", UserSessionMiddleware, taskController.GetByUserLogin)
		tasksRouterGroup.GET("", taskController.GetAll)
		tasksRouterGroup.GET("/by-status/:status", AdminSessionMiddleware, taskController.GetByStatus)
		tasksRouterGroup.GET("/by-priority/:priority", AdminSessionMiddleware, taskController.GetByPriority)
	}
}

func RegisterSwaggerAndMetricsHandlers() {
	Router.GET("/swagger/*any", AdminSessionMiddleware, ginSwagger.WrapHandler(swaggerFiles.Handler))
	Router.GET("/metrics", AdminSessionMiddleware, gin.WrapH(promhttp.Handler()))
}

func UserSessionMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user == nil {
		c.Redirect(http.StatusFound, "/")
	}
}

func AdminSessionMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(constant.UserSessionKey)
	if user == nil {
		c.Redirect(http.StatusFound, "/")
	}

	sessionUser, ok := user.(*repository.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if sessionUser.Role != constant.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "you should be admin for the action"})
		return
	}
}
