package test

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"path"
	"runtime"
)

func SetUpTestRouter() *gin.Engine {
	_, filename, line, ok := runtime.Caller(0)
	_ = line
	if !ok {
		log.Fatalln("cannot receive root project filename")
	}
	projectDir := path.Join(path.Dir(filename), "..")

	router := gin.Default()
	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("sessions", store))
	router.LoadHTMLGlob(projectDir + "/internal/server/templates/*.html")

	return router
}
