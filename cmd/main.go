package main

import (
	"fmt"
	"log"
	"todo-example/pkg/config"
	"todo-example/pkg/db"
	"todo-example/pkg/todo/routes"
	"todo-example/pkg/todo/services"

	"github.com/gin-gonic/gin"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	s := services.Server{
		H: h,
	}

	router := gin.Default()

	routes.TodoRoute(router, s)

	router.Run(fmt.Sprintf("localhost:%v", c.Port))
}
