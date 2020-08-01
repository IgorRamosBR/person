package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"person/internal/document"
)

type Repository interface {
	Find() ([]document.Person, error)
	FindById(id string) (document.Person, error)
	Create(document document.Person) (document.Person, error)
	Update(document document.Person) (int64, error)
	Delete(id primitive.ObjectID) (int64, error)
}
