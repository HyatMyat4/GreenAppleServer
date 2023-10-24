package users

import (
	"context"
	"green-apple-server/database"
	"green-apple-server/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Email_verified() gin.HandlerFunc {
	return func(req *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)

		otp := req.Param("otp")
		userId := req.Param("user_id")
		objID, _ := primitive.ObjectIDFromHex(userId)

		var user models.User

		var update_obj primitive.D

		err := database.UsersCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		defer cancle()

		if err != nil {
			req.JSON(http.StatusNotFound, gin.H{"message": "User Not Found"})
			return
		}

		update_obj = append(update_obj, bson.E{Key: "email_verified", Value: true})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"_id": objID}

		if user.Pin == otp {

			_, _err := database.UsersCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update_obj}}, &opt)

			if _err != nil {
				req.JSON(http.StatusNoContent, gin.H{"message": "Update failed"})
				return
			}

			req.JSON(http.StatusOK, gin.H{"message": "Email Verify Success"})
			return

		}
	}
}
