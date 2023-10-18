package company

import (
	"context"
	"green-apple-server/database"
	"green-apple-server/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Get_company() gin.HandlerFunc {
	return func(req *gin.Context) {
		company_Id := req.Param("company_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var company models.Comapny

		err := database.CompanyCollection.FindOne(ctx, bson.M{"company_id": company_Id}).Decode(&company)
		defer cancel()

		if err != nil {
			req.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		req.JSON(http.StatusOK, company)
	}
}
