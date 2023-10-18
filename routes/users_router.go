package routes

import (
	controllers "green-apple-server/controllers/users"

	"github.com/gin-gonic/gin"
)

func User_router(req *gin.Engine) {
	req.GET("/users", controllers.Get_users())
	req.GET("/user/:user_id", controllers.Get_user())
	req.POST("/user/create", controllers.Create_user())
}
