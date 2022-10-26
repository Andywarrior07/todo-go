package services

import (
	"context"
	"net/http"
	"time"
	"todo-example/pkg/db"
	"todo-example/pkg/todo/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Server struct {
	H db.Handler
}

func (s *Server) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var todos []models.Todo

	collection := s.H.DB.Database("todos").Collection("todo")
	result, err := collection.Find(ctx, bson.D{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}

	defer result.Close(ctx)

	for result.Next(ctx) {
		var todo models.Todo

		if err = result.Decode(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})

}

func (s *Server) GetById() {}

func (s *Server) Create() {}

func (s *Server) Update() {}

func (s *Server) Delete() {}
