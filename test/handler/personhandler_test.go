package handler

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/assert.v1"
	"net/http"
	"net/http/httptest"
	"person/internal/handler"
	"person/internal/useful"
	"person/test/mocks"
	"testing"
)

func TestDeletePersonByIdSuccess(t *testing.T) {

	id := "5f165e2e4de9b442e60b3904"
	objID, _ := primitive.ObjectIDFromHex(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)
	m.EXPECT().Delete(gomock.Eq(objID)).Return(int64(1), nil)

	r, _ := http.NewRequest("DELETE", "/person", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(m).Delete(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "\"\"", w.Body.String())
}

func TestDeletePersonByIdWithBrokenId(t *testing.T) {

	id := "5f165e2e4de9b442e60b39"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	r, _ := http.NewRequest("DELETE", "/person", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(m).Delete(w, r)
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

	m := mocks.NewMockRepository(ctrl)
	m.EXPECT().Delete(gomock.Eq(objID)).Return(int64(0), errors.New("Error"))

	r, _ := http.NewRequest("DELETE", "/person", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(m).Delete(w, r)
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

	m := mocks.NewMockRepository(ctrl)
	m.EXPECT().Delete(gomock.Eq(objID)).Return(int64(0), nil)

	r, _ := http.NewRequest("DELETE", "/person", nil)
	r = mux.SetURLVars(r, map[string]string{"id":id})
	w := httptest.NewRecorder()

	handler.NewPersonHandler(m).Delete(w, r)
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, map[string]string{"message":useful.PersonNotFound}, body)
}