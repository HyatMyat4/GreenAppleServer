package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id             primitive.ObjectID `bson:"_id"`
	User_name      string             `json:"user_name" validate:"required,min=2,max=80"`
	Password       string             `json:"password" validate:"required,min=8,max=30"`
	Email          string             `json:"email" validate:"email,required"`
	Phone          string             `json:"phone" validate:"required,min=5,max=30"`
	Role           string             `json:"role" validate:"required,eq=admin|eq=user"`
	Token          string             `json:"token" validate:"required"`
	Email_verified bool               `json:"email_verified" validate:"required"`
	Created_at     time.Time          `json:"created_at"`
	Updated_at     time.Time          `json:"updated_at"`
}

type ResponseCreateUser struct {
	Id             primitive.ObjectID `bson:"_id"`
	User_name      string             `json:"user_name" validate:"required,min=2,max=80"`
	Email          string             `json:"email" validate:"email,required"`
	Phone          string             `json:"phone" validate:"required,min=5,max=30"`
	Role           string             `json:"role" validate:"required,eq=admin|eq=user"`
	Token          string             `json:"token" validate:"required"`
	Email_verified bool               `json:"email_verified" validate:"required"`
	OTP            string             `json:"otp"`
	Created_at     time.Time          `json:"created_at"`
	Updated_at     time.Time          `json:"updated_at"`
}
