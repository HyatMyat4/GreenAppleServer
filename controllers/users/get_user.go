package users

import (
	"context"
	"fmt"
	"green-apple-server/database"
	models "green-apple-server/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Get_user() gin.HandlerFunc {
	return func(req *gin.Context) {
		userId := req.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var users models.User
		fmt.Printf(userId, "**")
		err := database.UsersCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&users)
		defer cancel()
		if err != nil {
			req.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

}
