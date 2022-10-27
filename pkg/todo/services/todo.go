package services

import (
	"context"
	"net/http"
	"time"
	"todo-example/pkg/db"
	"todo-example/pkg/todo/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	H db.Handler
}

var validate = validator.New()

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

	c.JSON(http.StatusOK, todos)
	// Si se quiere mandar un json
	// c.JSON(http.StatusOK, gin.H{"todos": todos})
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

	c.JSON(http.StatusOK, todo)

}

func (s *Server) Create(c *gin.Context) {
	collection := s.H.DB.Database("todos").Collection("todo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var todo models.Todo

	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTodo := models.Todo{
		ID:     primitive.NewObjectID(),
		Title:  todo.Title,
		Status: todo.Status,
	}

	result, err := collection.InsertOne(ctx, newTodo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (s *Server) Update(c *gin.Context) {
	collection := s.H.DB.Database("todos").Collection("todo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	defer cancel()

	var todo models.Todo

	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedData := bson.M{"title": todo.Title, "status": todo.Status}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updatedData})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)

}

func (s *Server) Delete() {}
