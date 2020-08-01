package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Person struct {
	Id    primitive.ObjectID `json:"id"`
	Name  string             `json:"name" validate:"required"`
	Email string             `json:"email" validate:"required,email"`
	Age   int8               `json:"age" validate:"required"`
}
