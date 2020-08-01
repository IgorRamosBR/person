package mapper

import (
	"person/internal/document"
	"person/internal/dto"
)

type Mapper interface {
	DocumentToDto(document document.Person) (dto.Person, error)
	ListDocumentToListDto(document []document.Person) ([]dto.Person, error)
	DtoToDocument(dto dto.Person) (document.Person, error)
}
