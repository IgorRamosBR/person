package configs

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Routes() {
	port := fmt.Sprintf(":%d", properties.Port)
	log.Infoln("Applications starting in port", port)

	r := mux.NewRouter()
	r.HandleFunc("/person", personHandler.Find).Methods(http.MethodGet)
	r.HandleFunc("/person/{id}", personHandler.FindById).Methods(http.MethodGet)
	r.HandleFunc("/person", personHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/person/{id}", personHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/person/{id}", personHandler.Delete).Methods(http.MethodDelete)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalln(err)
	}
}