package main

import (
	"crowdfunding-minpro-alterra/config"
	"crowdfunding-minpro-alterra/database"
	"crowdfunding-minpro-alterra/handler"
	"crowdfunding-minpro-alterra/modules/campaign"
	"crowdfunding-minpro-alterra/modules/donation"
	"crowdfunding-minpro-alterra/modules/payment"
	"crowdfunding-minpro-alterra/modules/user"
	"crowdfunding-minpro-alterra/utils/auth"
	"crowdfunding-minpro-alterra/utils/helper"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitConfigDB()

	db := database.ConnectDB(config.InitConfigDB())

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	donationRepository := donation.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	donationService := donation.NewService(donationRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	donationHandler := handler.NewDonationHandler(donationService)

	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images", "./images")

	api := router.Group("/api/v1")

	api.GET("/admin/users", authMiddleware(authService, userService), userHandler.GetAllUsers)

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/donations", authMiddleware(authService, userService), donationHandler.GetCampaignDonations)
	api.GET("/donations", authMiddleware(authService, userService), donationHandler.GetUserDonations)
	api.POST("/donations", authMiddleware(authService, userService), donationHandler.CreateDonation)
	api.POST("/donations/notification", donationHandler.GetNotification)

	router.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})
	
	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
