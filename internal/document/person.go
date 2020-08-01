package document

import "go.mongodb.org/mongo-driver/bson/primitive"

type Person struct {
	Id    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
	Age   int8               `bson:"age"`
}
