package mapper

import (
	"github.com/mitchellh/mapstructure"
	"person/internal/document"
	"person/internal/dto"
)

type PersonMapper struct {
}

func (p *PersonMapper) DocumentToDto(document document.Person) (dto.Person, error) {
	var person dto.Person
	err := mapstructure.Decode(document, &person)
	return person, err
}

func (p *PersonMapper) ListDocumentToListDto(document []document.Person) ([]dto.Person, error) {
	var people []dto.Person
	err := mapstructure.Decode(document, &people)
	return people, err
}

func (p *PersonMapper) DtoToDocument(dto dto.Person) (document.Person, error) {
	var person document.Person
	err := mapstructure.Decode(dto, &person)
	return person, err
}
