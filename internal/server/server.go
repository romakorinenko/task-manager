package server

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/romakorinenko/task-manager/internal/constant"
	"github.com/romakorinenko/task-manager/internal/controller"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"net/http"
)

var Router *gin.Engine

func RegisterServerAndHandlers(
	userController controller.IUserController,
	taskController controller.ITaskController,
	port int,
) {
	Router = gin.Default()
	Router.LoadHTMLGlob("templates/*")
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

	Router.POST("/users", UserSessionMiddleware, userController.Create)
	Router.PUT("/users/:id/block", UserSessionMiddleware, userController.Block)
	Router.GET("/users", UserSessionMiddleware, userController.GetAll)
}

func RegisterTaskHandlers(taskController controller.ITaskController) {
	Router.GET("/tasks/create", taskController.CreateTemplate)
	Router.POST("/tasks", UserSessionMiddleware, taskController.Create)
	Router.POST("/tasks/:id", UserSessionMiddleware, taskController.Update)
	Router.POST("/tasks/:id/delete", UserSessionMiddleware, taskController.Delete)
	Router.GET("/tasks/:id", UserSessionMiddleware, taskController.GetByID)
	Router.GET("/tasks/:id/edit", UserSessionMiddleware, taskController.Edit)
	Router.GET("/tasks/user/:login", UserSessionMiddleware, taskController.GetByUserLogin)
	Router.GET("/tasks", taskController.GetAll)
	Router.GET("/tasks/by-status/:status", taskController.GetByStatus)
	Router.GET("/tasks/by-priority/:priority", taskController.GetByPriority)
}

func RegisterSwaggerAndMetricsHandlers() {
	Router.GET("/swagger/*any", UserSessionMiddleware, ginSwagger.WrapHandler(swaggerFiles.Handler))
	Router.GET("/metrics", UserSessionMiddleware, gin.WrapH(promhttp.Handler()))
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
}
