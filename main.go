package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ktsoator/connectify/internal/repository"
	"github.com/ktsoator/connectify/internal/repository/dao"
	"github.com/ktsoator/connectify/internal/service"
	"github.com/ktsoator/connectify/internal/web"
	"gorm.io/gorm"
)

func main() {
	router := web.InitRouter()
	db := dao.InitDB()
	initUser(db, router)
	router.Run(":8080")
}

func initUser(db *gorm.DB, router *gin.Engine) {
	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepo)
	userHandler := web.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)
}
