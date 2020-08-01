package configs

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"net/http"
	_ "person/docs"
)

func Routes() {
	port := fmt.Sprintf(":%d", properties.Port)
	log.Infoln("Applications starting in port", port)

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/v1/person", personHandler.Find).Methods(http.MethodGet)
	r.HandleFunc("/v1/person/{id}", personHandler.FindById).Methods(http.MethodGet)
	r.HandleFunc("/v1/person", personHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/v1/person/{id}", personHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/v1/person/{id}", personHandler.Delete).Methods(http.MethodDelete)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalln(err)
	}
}
