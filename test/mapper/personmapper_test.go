package mapper

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"person/internal/document"
	dto2 "person/internal/dto"
	"person/internal/mapper"
	"testing"
)

func TestShouldReturnDTOFilled(t *testing.T) {

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	doc := document.Person{
		Id:    objID,
		Name:  "Lucas",
		Email: "lucas@gmail.com",
		Age:   22,
	}

	personMapper := &mapper.PersonMapper{}
	dto, err := personMapper.DocumentToDto(doc)

	assert.Nil(t, err)
	assert.Equal(t, dto.Id, doc.Id)
	assert.Equal(t, dto.Name, doc.Name)
	assert.Equal(t, dto.Email, doc.Email)
	assert.Equal(t, dto.Age, doc.Age)
}

func TestShouldReturnListDTOFilled(t *testing.T) {

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	docs := []document.Person{
		{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22},
		{Id: objID, Email: "test@gmail.com", Age: 20},
	}

	personMapper := &mapper.PersonMapper{}
	dtos, err := personMapper.ListDocumentToListDto(docs)

	assert.Nil(t, err)
	assert.Len(t, dtos, 2)
	assert.Equal(t, dtos[0].Id, docs[0].Id)
	assert.Equal(t, dtos[0].Name, docs[0].Name)
	assert.Equal(t, dtos[0].Email, docs[0].Email)
	assert.Equal(t, dtos[0].Age, docs[0].Age)
	assert.Equal(t, dtos[1].Id, docs[1].Id)
	assert.Equal(t, dtos[1].Name, "")
	assert.Equal(t, dtos[1].Email, docs[1].Email)
	assert.Equal(t, dtos[1].Age, docs[1].Age)
}

func TestShouldReturnDocumentFilled(t *testing.T) {

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	dto := dto2.Person{
		Id:    objID,
		Name:  "Lucas",
		Email: "lucas@gmail.com",
		Age:   22,
	}

	personMapper := &mapper.PersonMapper{}
	doc, err := personMapper.DtoToDocument(dto)

	assert.Nil(t, err)
	assert.Equal(t, doc.Id, dto.Id)
	assert.Equal(t, doc.Name, dto.Name)
	assert.Equal(t, doc.Email, dto.Email)
	assert.Equal(t, doc.Age, dto.Age)
}
