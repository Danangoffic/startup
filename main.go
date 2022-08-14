package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// untuk set konfigurasi ke database mysql. example: localhost root
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	// untuk open connection ke db mysql dan cek error nya menggunakan gorm
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// untuk print log gagal konek ke db
	if err != nil {
		log.Fatal(err.Error())
	}

	// REPOSITORY LIST
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	// SERVICE LIST
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService(transactionRepository, campaignRepository)
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)
	authService := auth.NewService()

	// userInput, _ := userService.GetUserByID(2)
	// fmt.Println("user : ", userInput)

	// inputTransaction := transaction.CreateTransactionInput{
	// 	CampaignId: 5,
	// 	Amount:     100000000,
	// 	User:       userInput,
	// }
	// fmt.Println("input transaction : ", inputTransaction)

	// newTransaction, err := transactionService.CreateTransaction(inputTransaction)
	// if err != nil {
	// 	fmt.Println("error create transaction", err.Error())
	// }
	// fmt.Println("success create a new transaction", newTransaction)

	// HANDLER LIST
	// get user handler yang akan dipakai pada api routes
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService, paymentService)

	// API LIST
	// untuk penggunaan api dari gin
	router := gin.Default()
	// static for image api
	router.Static("/images", "./images")
	// api target with version
	api := router.Group("/api/v1")

	// USERS API
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checker", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	// CAMPAIGNS API
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.GET("/campaigns", campaignHandler.FindCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	// TRANSACTIONS API
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)
	// api.DELETE("/campaign-images/:id", authMiddleware(authService, userService), campaignHandler.DeleteCampaignImage)

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
		arrayAuthHeader := strings.Split(authHeader, " ")
		if len(arrayAuthHeader) == 2 {
			tokenString = arrayAuthHeader[1]
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

		userId := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userId)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
