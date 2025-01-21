package server

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/romakorinenko/task-manager/internal/controller"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	slog.Info("server ")
	slog.Info("Swagger available on http://localhost:8080/swagger/index.html")
	_ = Router.Run(fmt.Sprintf(":%d", port))
}

func RegisterUserHandlers(userController controller.IUserController) {
	Router.GET("/", userController.GetMainPage)
	Router.POST("/login", userController.Login)
	Router.GET("/logout", userController.Logout)

	Router.POST("/users", userController.Create)
	Router.PUT("/users/:id/block", userController.Block)
	Router.GET("/users", userController.GetAll)
}

func RegisterTaskHandlers(taskController controller.ITaskController) {
	Router.GET("/tasks/create", taskController.CreateTemplate)
	Router.POST("/tasks", taskController.Create)
	Router.POST("/tasks/:id", taskController.Update)
	Router.POST("/tasks/:id/delete", taskController.Delete)
	Router.GET("/tasks/:id", taskController.GetByID)
	Router.GET("/tasks/:id/edit", taskController.Edit)
	Router.GET("/tasks/user/:login", taskController.GetByUserLogin)
	Router.GET("/tasks", taskController.GetAll)
}

func RegisterSwaggerAndMetricsHandlers() {
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	Router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
