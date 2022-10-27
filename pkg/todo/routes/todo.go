package routes

import (
	"todo-example/pkg/todo/services"

	"github.com/gin-gonic/gin"
)

func TodoRoute(router *gin.Engine, s services.Server) {
	router.GET("/todos", s.GetAll)
	router.GET("/todos/:id", s.GetById)
	router.POST("/todos", s.Create)
	router.PUT("/todos/:id", s.Update)
	router.DELETE("/todos/:id", s.Delete)
}
