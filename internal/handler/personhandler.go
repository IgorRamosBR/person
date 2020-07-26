package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"person/internal/dto"
	"person/internal/mapper"
	"person/internal/repository"
	"person/internal/useful"
)

type PersonHandler struct {
	Mapper mapper.Mapper
	Repository repository.Repository
}

func NewPersonHandler(mapper mapper.Mapper, repo repository.Repository) *PersonHandler {
	return &PersonHandler{Mapper: mapper, Repository: repo}
}

func (p *PersonHandler) Find(w http.ResponseWriter, r *http.Request) {

	log.Infoln(useful.FindAll)

	peopleDocument, err := p.Repository.Find()

	if err != nil {
		log.Errorln(useful.GetDataFromDbError, err)
		useful.BuildError(w, http.StatusInternalServerError, useful.InternalErrorOccurred)
		return
	}

	peopleDTO, err := p.Mapper.ListDocumentToListDto(peopleDocument)

	if err != nil {
		log.Errorln(useful.ParserError, err)
		useful.BuildError(w, http.StatusInternalServerError, useful.InternalErrorOccurred)
		return
	}

	if peopleDTO == nil {
		useful.BuildSuccess(w, http.StatusOK, []string{})
		return
	}

	useful.BuildSuccess(w, http.StatusOK, peopleDTO)
}

func (p *PersonHandler) FindById(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	log.Infoln(useful.FindById, id)

	personDocument, err := p.Repository.FindById(id)

	if err != nil {
		log.Errorln(useful.PersonNotFound, err)
		useful.BuildError(w, http.StatusNotFound, useful.PersonNotFound)
		return
	}

	personDTO, err := p.Mapper.DocumentToDto(personDocument)

	if err != nil {
		log.Errorln(useful.ParserError, err)
		useful.BuildError(w, http.StatusInternalServerError, useful.ParserError)
		return
	}

	useful.BuildSuccess(w, http.StatusOK, personDTO)
}

func (p *PersonHandler) Create(w http.ResponseWriter, r *http.Request) {

	v := validator.New()
	var body dto.Person

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Errorln(useful.ParserError, err)
		useful.BuildError(w, http.StatusUnprocessableEntity, useful.BrokenBody)
		return
	}

	log.Infoln(useful.Create, body)

	if err := v.Struct(body); err != nil {
		log.Errorln(useful.ValidateBodyError, err)
		useful.BuildError(w, http.StatusBadRequest, useful.BrokenBody)
		return
	}

	personDocument, err := p.Mapper.DtoToDocument(body)

	if err != nil {
		log.Errorln(useful.ParserError, err)
		useful.BuildError(w, http.StatusInternalServerError, useful.ParserError)
		return
	}

	personDocument, err = p.Repository.Create(personDocument)

	if err != nil {
		log.Errorln(useful.CreateError, err)
		useful.BuildError(w, http.StatusBadRequest, useful.CreateError)
		return
	}

	personDTO, err := p.Mapper.DocumentToDto(personDocument)

	if err != nil {
		log.Errorln(useful.ParserError, err)
		useful.BuildError(w, http.StatusInternalServerError, useful.ParserError)
		return
	}

	useful.BuildSuccess(w, http.StatusCreated, personDTO)
}

func (p *PersonHandler) Update(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	v := validator.New()
	var body dto.Person

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Errorln(useful.ParserError, err)
		useful.BuildError(w, http.StatusUnprocessableEntity, useful.BrokenBody)
		return
	}

	log.Infoln(useful.Update, id)

	if err := v.Struct(body); err != nil {
		log.Errorln(useful.ValidateBodyError, err)
		useful.BuildError(w, http.StatusBadRequest, useful.BrokenBody)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Errorln(useful.ParserError, err)
		useful.BuildError(w, http.StatusUnprocessableEntity, useful.BrokenId)
		return
	}

	body.Id = objID
	personDocument, err := p.Mapper.DtoToDocument(body)

	if err != nil {
		log.Errorln(useful.ParserError, err)
		useful.BuildError(w, http.StatusInternalServerError, useful.ParserError)
		return
	}

	count, err := p.Repository.Update(personDocument)

	if err != nil {
		log.Errorln(useful.UpdateError, err)
		useful.BuildError(w, http.StatusBadRequest, useful.UpdateError)
		return
	}

	if count == 0 {
		log.Errorln(useful.PersonNotFound)
		useful.BuildError(w, http.StatusNotFound, useful.PersonNotFound)
		return
	}

	personDTO, err := p.Mapper.DocumentToDto(personDocument)

	if err != nil {
		log.Errorln(useful.ParserError, err)
		useful.BuildError(w, http.StatusInternalServerError, useful.ParserError)
		return
	}

	useful.BuildSuccess(w, http.StatusOK, personDTO)
}

func (p *PersonHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Errorln(useful.BrokenId, err)
		useful.BuildError(w, http.StatusBadRequest, useful.BrokenId)
		return
	}

	log.Infoln(useful.Delete, id)

	count, err := p.Repository.Delete(objID)

	if err != nil {
		log.Errorln(useful.DeleteError, err)
		useful.BuildError(w, http.StatusBadRequest, useful.DeleteError)
		return
	}

	if count == 0 {
		log.Errorln(useful.PersonNotFound)
		useful.BuildError(w, http.StatusNotFound, useful.PersonNotFound)
		return
	}

	useful.BuildSuccess(w, http.StatusNoContent, "")
}