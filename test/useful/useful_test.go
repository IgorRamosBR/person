package useful

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"person/internal/dto"
	"person/internal/useful"
	"testing"
)

func TestBuildSuccess200WithBody(t *testing.T) {

	response := httptest.NewRecorder()
	objID, _ := primitive.ObjectIDFromHex("5f165e2e4de9b442e60b3904")

	body := dto.Person{
		Id:    objID,
		Name:  "Lucas",
		Email: "lucas@gmail.com",
		Age:   22,
	}

	useful.BuildSuccess(response, http.StatusOK, body)

	var bodyReceived dto.Person
	_ = json.Unmarshal([]byte(response.Body.String()), &bodyReceived)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, body, bodyReceived)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
}

func TestBuildSuccess200WithEmptyArray(t *testing.T) {

	response := httptest.NewRecorder()

	var body []string
	useful.BuildSuccess(response, http.StatusOK, body)

	var bodyReceived []string
	_ = json.Unmarshal([]byte(response.Body.String()), &bodyReceived)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, body, bodyReceived)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
}

func TestBuildSuccess204WithoutBody(t *testing.T) {

	response := httptest.NewRecorder()

	useful.BuildSuccess(response, http.StatusNoContent, "")

	assert.Equal(t, http.StatusNoContent, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
}

func TestBuildError(t *testing.T) {

	response := httptest.NewRecorder()

	useful.BuildError(response, http.StatusInternalServerError, useful.CreateError)

	var body dto.Error
	_ = json.Unmarshal([]byte(response.Body.String()), &body)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, dto.Error{Message: useful.CreateError}, body)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
}
