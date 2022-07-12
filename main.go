package main

import (
	"bwastartup/handler"
	"bwastartup/user"
	"log"

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

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Riselda Lalusu"
	// userInput.Email = "riselda@gmail.com"
	// userInput.Occupation = "Senior Programmer"
	// userInput.Password = "cobapassword"

	// userService.RegisterUser(userInput)

	userHandler := handler.NewUserHandler(userService)
	// untuk penggunaan api
	router := gin.Default()

	// api versioning
	api := router.Group("/api/v1")

	// api list
	api.POST("/users", userHandler.RegisterUser)

	router.Run()
	// userRepository.Save(user)

	// input dari user
	// handler : mapping input user menjadi sebuah struct inpput
	// service : mapping struct input ke struct User
	// repository
	// db
}
