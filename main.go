package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/logeshwarann-dev/bits-bank_auth-service/api"
	"github.com/logeshwarann-dev/bits-bank_auth-service/db"
	"github.com/logeshwarann-dev/bits-bank_auth-service/session"
	"github.com/logeshwarann-dev/bits-bank_auth-service/utils"
)

func init() {
	utils.LoadEnv()
	api.PgDb = db.ConnectToDB()
	session.InitConnToRedis()
}

func main() {
	fmt.Println("Auth service")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.POST("/auth/v1/sign-up", api.SignUp)
	router.POST("/auth/v1/sign-in", api.SignIn)
	router.GET("/auth/v1/get-user", api.GetLoggedInUser)
	router.GET("/auth/v1/logout", api.SignOut)

	router.Run(":8080")
}
