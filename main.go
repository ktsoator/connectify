package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ktsoator/connectify/internal/repository"
	"github.com/ktsoator/connectify/internal/repository/dao"
	"github.com/ktsoator/connectify/internal/service"
	"github.com/ktsoator/connectify/internal/web"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	db := dao.InitDB()
	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepo)
	userHandler := web.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)

	router.Run(":8080")
}
