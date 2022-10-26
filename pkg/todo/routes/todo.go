package routes

import (
	"todo-example/pkg/todo/services"

	"github.com/gin-gonic/gin"
)

func TodoRoute(router *gin.Engine, s services.Server) {
	router.GET("/todos", s.GetAll)
}