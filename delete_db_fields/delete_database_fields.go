package deletedbfields_test

import (
	"context"
	"green-apple-server/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Delete_db_fields() gin.HandlerFunc {
	return func(req *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// Define the filter to match the documents you want to update (empty filter means all documents).
		filter := bson.M{}

		// Define the update operation using $unset to remove the "type" field.
		update := bson.M{
			"$unset": bson.M{
				"companies_id": "",
			},
		}

		updateResult, error := database.UsersCollection.UpdateMany(ctx, filter, update)
		defer cancel()
		if error != nil {
			req.JSON(http.StatusInternalServerError, error.Error())
			return
		}
		req.JSON(http.StatusOK, updateResult)
	}
}
