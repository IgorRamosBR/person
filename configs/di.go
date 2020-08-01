package configs

import (
	"person/internal/handler"
	"person/internal/mapper"
	"person/internal/repository"
)

var personHandler *handler.PersonHandler

func Di() {
	person()
}

func person() {
	personRepository := repository.PersonRepository{Collection: mongo()}
	personMapper := mapper.PersonMapper{}
	personHandler = handler.NewPersonHandler(&personMapper, personRepository)
}
