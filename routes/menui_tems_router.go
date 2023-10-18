package routes

import (
	controllers "green-apple-server/controllers/menuitems"

	"github.com/gin-gonic/gin"
)

func Menu_items_router(req *gin.Engine) {
	req.GET("/menuitems", controllers.Get_menuitems())
	req.GET("/menuitem/:menu_id", controllers.Get_menuitem())
}
