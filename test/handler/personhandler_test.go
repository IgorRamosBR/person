package handler

import (
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

func TestDeletePersonByIdSuccess(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	objID, _ := primitive.ObjectIDFromHex(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)

	repo.EXPECT().Delete(gomock.Eq(objID)).Return(int64(1), nil)

	r, _ := http.NewRequest("DELETE", "/person", nil)
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

	r, _ := http.NewRequest("DELETE", "/person", nil)
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

	r, _ := http.NewRequest("DELETE", "/person", nil)
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

	r, _ := http.NewRequest("DELETE", "/person", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(&mapper.PersonMapper{}, repo).Delete(w, r)
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, map[string]string{"message":useful.PersonNotFound}, body)
}