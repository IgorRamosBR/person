package configs

import (
	"person/internal/repository"
	"person/internal/handler"
)

var personHandler handler.PersonHandler

func Di() {
	person()
}

func person() {
	personRepository := repository.PersonRepository{Collection: mongo()}
	personHandler = handler.PersonHandler{Repository: personRepository}
}