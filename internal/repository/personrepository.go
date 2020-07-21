package repository

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"person/internal/document"
)

type PersonRepository struct {
	Collection *mongo.Collection
}

func (p PersonRepository) Find() ([]document.Person, error) {

	var people []document.Person
	ctx := context.TODO()

	cur, err := p.Collection.Find(ctx, bson.D{})

	if err == nil {
		for cur.Next(ctx) {
			var result document.Person
			err := cur.Decode(&result)
			if err != nil {
				log.Error(err)
			}
			people = append(people, result)
		}
	}

	return people, err
}

func (p PersonRepository) FindById(id string) (document.Person, error) {

	var person document.Person
	ctx := context.TODO()

	objectID, _ := primitive.ObjectIDFromHex(id)
	result := p.Collection.FindOne(ctx, bson.M{"_id": objectID})
	err := result.Decode(&person)

	return person, err
}

func (p PersonRepository) Create(document document.Person) (document.Person, error) {

	ctx := context.TODO()

	document.Id = primitive.NewObjectID()
	_, err := p.Collection.InsertOne(ctx, document)

	if err != nil {
		return document, err
	}

	return document, err
}

func (p PersonRepository) Update(document document.Person) (int64, error) {

	ctx := context.TODO()
	filter := bson.M{"_id": document.Id}
	update := bson.M{"$set": bson.M{
		"name": document.Name,
		"email": document.Email,
		"age": document.Age,
	}}

	result, err := p.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return 0, err
	}

	return result.MatchedCount, err
}

func (p PersonRepository) Delete(id primitive.ObjectID) (int64, error) {

	ctx := context.TODO()
	filter := bson.M{"_id": id}

	result, err := p.Collection.DeleteOne(ctx, filter)

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, err
}