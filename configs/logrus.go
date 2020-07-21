package configs

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func Logrus() {

	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)

	level, err := log.ParseLevel(properties.Log.Level)

	if err == nil {
		log.SetLevel(level)
	}

	if properties.Log.JsonFormatter {
		log.SetFormatter(&log.JSONFormatter{})
	}
}