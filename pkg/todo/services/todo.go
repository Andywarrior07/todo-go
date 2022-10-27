package services

import (
	"context"
	"net/http"
	"time"
	"todo-example/pkg/db"
	"todo-example/pkg/todo/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return
	}

	defer result.Close(ctx)

	for result.Next(ctx) {
		var todo models.Todo

		if err = result.Decode(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func (s *Server) GetById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	defer cancel()

	var todo models.Todo

	collection := s.H.DB.Database("todos").Collection("todo")
	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&todo)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})

}

func (s *Server) Create() {}

func (s *Server) Update() {}

func (s *Server) Delete() {}
