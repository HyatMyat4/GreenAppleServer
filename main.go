package main

import (
	"green-apple-server/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET, POST, PUT, DELETE, OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true // Allow credentials such as cookies
	config.AllowWildcard = true
	router := gin.New()
	router.Use(cors.New(config))

	router.Use(gin.Logger())

	routes.Auth_router(router)
	routes.User_router(router)
	routes.Menu_items_router(router)
	routes.Company_router(router)
	//routes.Auth_router(router)
	router.Run(":" + port)

}

//nodemon --exec go run main.go --signal SIGTERM
