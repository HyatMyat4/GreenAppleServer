package token

import (
	"fmt"
	"green-apple-server/database"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type SignedDetails struct {
	email   string
	name    string
	user_id string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, name string, user_id string) (access_Token string, err error) {
	access_token := &SignedDetails{
		email:   email,
		name:    name,
		user_id: user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24*30*2)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, access_token).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err

}

func UpdateAllToken(signedToken string, user_id string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true

	filter := bson.M{"Id": user_id}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := database.UsersCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	return
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	//the token is invalid

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	//the token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg
}
