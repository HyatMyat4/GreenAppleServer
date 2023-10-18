package routes

import (
	delete_database_filelds "green-apple-server/delete_db_fields"

	"github.com/gin-gonic/gin"
)

func Delete_fields_router(req *gin.Engine) {
	req.GET("/databasefields/delete/:id", delete_database_filelds.Delete_db_fields())
}
