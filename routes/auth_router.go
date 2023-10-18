package routes

import (
	controllers "green-apple-server/controllers/auth"

	"github.com/gin-gonic/gin"
)

func Auth_router(req *gin.Engine) {
	req.POST("/signup", controllers.Signup())
	req.POST("/login", controllers.Login())
}
