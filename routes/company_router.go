package routes

import (
	controllers "green-apple-server/controllers/company"

	"github.com/gin-gonic/gin"
)

func Company_router(req *gin.Engine) {
	req.GET("/company/:company_id", controllers.Get_company())
}
