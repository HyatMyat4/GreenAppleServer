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
	config.AllowOrigins = []string{"*"} // Replace with your frontend application's domain
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

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
