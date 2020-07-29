package handler

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"person/internal/document"
	"person/internal/dto"
	"person/internal/handler"
	"person/internal/mapper"
	"person/internal/useful"
	"person/test/mocks"
	"testing"
)

func TestFindSuccess(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	objID2, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3905")

	docs := []document.Person{
		{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22},
		{Id: objID2, Email: "test@gmail.com", Age: 20},
	}

	repo.EXPECT().Find().Return(docs, nil)

	r, _ := http.NewRequest("GET", "/person", nil)
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Find(w, r)

	var body []dto.Person
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, docs[0].Id, body[0].Id)
	assert.Equal(t, docs[0].Name, body[0].Name)
	assert.Equal(t, docs[0].Email, body[0].Email)
	assert.Equal(t, docs[0].Age, body[0].Age)
	assert.Equal(t, docs[1].Id, body[1].Id)
	assert.Equal(t, docs[1].Name, body[1].Name)
	assert.Equal(t, docs[1].Email, body[1].Email)
	assert.Equal(t, docs[1].Age, body[1].Age)
}

func TestFindReturningErrorFromDatabaseWhenTryFind(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().Find().Return(nil, errors.New("database error"))

	r, _ := http.NewRequest("GET", "/person", nil)
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Find(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, map[string]string{"message":useful.InternalErrorOccurred}, body)
}

func TestFindReturningErrorWhenTryDoMapper(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	mapp := mocks.NewMockMapper(ctrl)

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")

	docs := []document.Person{
		{Id: objID, Name: "Lucas", Email: "lucas@@gmail.com", Age: 22},
	}

	repo.EXPECT().Find().Return(docs, nil)
	mapp.EXPECT().ListDocumentToListDto(docs).Return(nil, errors.New("mapper error"))

	r, _ := http.NewRequest("GET", "/person", nil)
	w := httptest.NewRecorder()

	handler.NewPersonHandler(mapp, repo).Find(w, r)

	var body map[string]string
	_ = json.Unmarshal([]byte(w.Body.String()), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, map[string]string{"message":useful.InternalErrorOccurred}, body)
}

func TestFindSuccessReturningEmptyBody(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().Find().Return(nil, nil)

	r, _ := http.NewRequest("GET", "/person", nil)
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Find(w, r)

	var body []struct{}
	_ = json.Unmarshal([]byte(w.Body.String()), &body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, body, 0)
	assert.Empty(t, body)
}

func TestFindByIdSuccess(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	objID, _ := primitive.ObjectIDFromHex(id)
	doc := document.Person{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().FindById(gomock.Eq(id)).Return(doc, nil)

	r, _ := http.NewRequest("GET", "/person/{id}", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).FindById(w, r)

	var body dto.Person
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, doc.Id, body.Id)
	assert.Equal(t, doc.Name, body.Name)
	assert.Equal(t, doc.Email, body.Email)
	assert.Equal(t, doc.Age, body.Age)
}

func TestFindByIdReturningErrorFromDatabaseWhenTryFind(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().FindById(gomock.Eq(id)).Return(document.Person{}, errors.New("Find error"))

	r, _ := http.NewRequest("GET", "/person/{id}", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).FindById(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, map[string]string{"message":useful.PersonNotFound}, body)
}

func TestFindByIdReturningErrorWhenTryDoMapper(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	objID, _ := primitive.ObjectIDFromHex(id)
	doc := document.Person{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().FindById(gomock.Eq(id)).Return(doc, nil)

	mapp := mocks.NewMockMapper(ctrl)
	mapp.EXPECT().DocumentToDto(gomock.Eq(doc)).Return(dto.Person{}, errors.New("Mapper Error"))

	r, _ := http.NewRequest("GET", "/person/{id}", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(mapp, repo).FindById(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, map[string]string{"message":useful.ParserError}, body)
}

func TestCreateSuccess(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	docWithoutId := document.Person{Name: "Lucas", Email: "lucas@gmail.com", Age: 22}
	doc := document.Person{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	repo.EXPECT().Create(gomock.Eq(docWithoutId)).Return(doc, nil)

	bodySent, _ := json.Marshal(docWithoutId)

	r, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(bodySent))
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Create(w, r)

	var body dto.Person
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, doc.Id, body.Id)
	assert.Equal(t, doc.Name, body.Name)
	assert.Equal(t, doc.Email, body.Email)
	assert.Equal(t, doc.Age, body.Age)
}

func TestCreateBrokenBody(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	bodySent, _ := json.Marshal(struct{Id int; Text string}{123, "Broked"})

	r, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(bodySent))
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Create(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, map[string]string{"message":useful.BrokenBody}, body)
}

func TestCreateValidatingContentOnFields(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	bodySent, _ := json.Marshal(dto.Person{Name: "Lucas", Email: "lucas@@gmail.com", Age: 22})

	r, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(bodySent))
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Create(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, map[string]string{"message":useful.BrokenBody}, body)
}

func TestCreateReturningErrorWhenTryDoMapperToDocument(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	doc := dto.Person{Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	repo := mocks.NewMockRepository(ctrl)
	mapp := mocks.NewMockMapper(ctrl)
	mapp.EXPECT().DtoToDocument(doc).Return(document.Person{}, errors.New("mapper error"))

	bodySent, _ := json.Marshal(doc)

	r, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(bodySent))
	w := httptest.NewRecorder()

	handler.NewPersonHandler(mapp, repo).Create(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, map[string]string{"message":useful.ParserError}, body)
}

func TestCreateReturningErrorFromDatabaseWhenTryCreate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	doc := document.Person{Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().Create(gomock.Eq(doc)).Return(document.Person{}, errors.New("Create error"))

	bodySent, _ := json.Marshal(doc)

	r, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(bodySent))
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Create(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, map[string]string{"message":useful.CreateError}, body)
}

func TestCreateReturningErrorWhenTryDoMapperToDTO(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pdto := dto.Person{Name: "Lucas", Email: "lucas@gmail.com", Age: 22}
	doc := document.Person{Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	mapp := mocks.NewMockMapper(ctrl)
	mapp.EXPECT().DtoToDocument(pdto).Return(doc, nil)
	mapp.EXPECT().DocumentToDto(doc).Return(pdto, errors.New("Mapper error"))

	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().Create(gomock.Eq(doc)).Return(doc, nil)

	bodySent, _ := json.Marshal(doc)

	r, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(bodySent))
	w := httptest.NewRecorder()

	handler.NewPersonHandler(mapp, repo).Create(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, map[string]string{"message":useful.ParserError}, body)
}

func TestUpdateSuccess(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	docWithoutId := document.Person{Name: "Lucas", Email: "lucas@gmail.com", Age: 22}
	doc := document.Person{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	repo.EXPECT().Update(gomock.Eq(doc)).Return(int64(1), nil)

	bodySent, _ := json.Marshal(docWithoutId)

	r, _ := http.NewRequest("PUT", "/person", bytes.NewBuffer(bodySent))
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Update(w, r)

	var body dto.Person
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, doc.Id, body.Id)
	assert.Equal(t, doc.Name, body.Name)
	assert.Equal(t, doc.Email, body.Email)
	assert.Equal(t, doc.Age, body.Age)
}

func TestUpdateBrokenBody(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	bodySent, _ := json.Marshal(struct{Id int; Text string}{123, "Broked"})

	r, _ := http.NewRequest("PUT", "/person", bytes.NewBuffer(bodySent))
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Update(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, map[string]string{"message":useful.BrokenBody}, body)
}

func TestUpdateValidatingContentOnFields(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	bodySent, _ := json.Marshal(dto.Person{Name: "Lucas", Email: "lucas@@gmail.com", Age: 22})

	r, _ := http.NewRequest("PUT", "/person", bytes.NewBuffer(bodySent))
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Update(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, map[string]string{"message":useful.BrokenBody}, body)
}

func TestUpdateValidatingID(t *testing.T) {

	id := "5f165e2e4de9b442e60b3"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	bodySent, _ := json.Marshal(dto.Person{Name: "Lucas", Email: "lucas@gmail.com", Age: 22})

	r, _ := http.NewRequest("PUT", "/person", bytes.NewBuffer(bodySent))
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Update(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, map[string]string{"message":useful.BrokenId}, body)
}

func TestUpdateReturningErrorWhenTryDoMapperToDocument(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	objID, _ := primitive.ObjectIDFromHex(id)
	doc := dto.Person{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	repo := mocks.NewMockRepository(ctrl)
	mapp := mocks.NewMockMapper(ctrl)
	mapp.EXPECT().DtoToDocument(doc).Return(document.Person{}, errors.New("mapper error"))

	bodySent, _ := json.Marshal(doc)

	r, _ := http.NewRequest("PUT", "/person", bytes.NewBuffer(bodySent))
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(mapp, repo).Update(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, map[string]string{"message":useful.ParserError}, body)
}

func TestUpdateReturningErrorFromDatabaseWhenTryCreate(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	docWithoutId := document.Person{Name: "Lucas", Email: "lucas@gmail.com", Age: 22}
	doc := document.Person{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().Update(gomock.Eq(doc)).Return(int64(0), errors.New("Create error"))

	bodySent, _ := json.Marshal(docWithoutId)

	r, _ := http.NewRequest("PUT", "/person", bytes.NewBuffer(bodySent))
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Update(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, map[string]string{"message":useful.UpdateError}, body)
}

func TestUpdateReturningFromDatabaseZeroDocumentAffected(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	docWithoutId := document.Person{Name: "Lucas", Email: "lucas@gmail.com", Age: 22}
	doc := document.Person{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().Update(gomock.Eq(doc)).Return(int64(0), nil)

	bodySent, _ := json.Marshal(docWithoutId)

	r, _ := http.NewRequest("PUT", "/person", bytes.NewBuffer(bodySent))
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Update(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, map[string]string{"message":useful.PersonNotFound}, body)
}

func TestUpdateReturningErrorWhenTryDoMapperToDTO(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")
	pdto := dto.Person{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22}
	doc := document.Person{Id: objID, Name: "Lucas", Email: "lucas@gmail.com", Age: 22}

	mapp := mocks.NewMockMapper(ctrl)
	mapp.EXPECT().DtoToDocument(pdto).Return(doc, nil)
	mapp.EXPECT().DocumentToDto(doc).Return(pdto, errors.New("Mapper error"))

	repo := mocks.NewMockRepository(ctrl)
	repo.EXPECT().Update(gomock.Eq(doc)).Return(int64(1), nil)

	bodySent, _ := json.Marshal(doc)

	r, _ := http.NewRequest("PUT", "/person", bytes.NewBuffer(bodySent))
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(mapp, repo).Update(w, r)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, map[string]string{"message":useful.ParserError}, body)
}

func TestDeletePersonByIdSuccess(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	objID, _ := primitive.ObjectIDFromHex(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().Delete(gomock.Eq(objID)).Return(int64(1), nil)

	r, _ := http.NewRequest("DELETE", "/person/{id}", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "\"\"", w.Body.String())
}

func TestDeletePersonByIdWithBrokenId(t *testing.T) {

	id := "5f165e2e4de9b442e60b39"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	r, _ := http.NewRequest("DELETE", "/person/{id}", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Delete(w, r)
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, map[string]string{"message":useful.BrokenId}, body)
}

func TestDeletePersonByIdWithErrorReturnedFromRepository(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	objID, _ := primitive.ObjectIDFromHex(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().Delete(gomock.Eq(objID)).Return(int64(0), errors.New("Error"))

	r, _ := http.NewRequest("DELETE", "/person/{id}", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Delete(w, r)
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, map[string]string{"message":useful.DeleteError}, body)
}

func TestDeletePersonByIdReturningZeroAffected(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	objID, _ := primitive.ObjectIDFromHex(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().Delete(gomock.Eq(objID)).Return(int64(0), nil)

	r, _ := http.NewRequest("DELETE", "/person/{id}", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Delete(w, r)
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, map[string]string{"message":useful.PersonNotFound}, body)
}