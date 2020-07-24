package configs

import (
	"person/internal/handler"
	"person/internal/repository"
)

var personHandler *handler.PersonHandler

func Di() {
	person()
}

func person() {
	personRepository := repository.PersonRepository{Collection: mongo()}
	personHandler = handler.NewPersonHandler(personRepository)
}