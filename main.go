package main

import (
	"crowdfunding-minpro-alterra/config"
	"crowdfunding-minpro-alterra/database"
	"crowdfunding-minpro-alterra/handler"
	"crowdfunding-minpro-alterra/modules/user"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitConfigDB()

	db := database.ConnectDB(config.InitConfigDB())

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()
}