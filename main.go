package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ktsoator/connectify/internal/web"
)

func main() {
	router := gin.Default()
	userHandler := web.NewUserHandler()
	userHandler.RegisterRoutes(router)
	router.Run(":8080")
}
