package configs

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Properties struct {
	Mongo struct {
		Uri        string
		Database   string
		Collection string
	}
	Port int
	Log  struct {
		Level         string
		JsonFormatter bool
	}
}

var properties Properties

func Variables() {
	variables := make(map[string]interface{})

	var basepath string

	if basepath, _ = os.Getwd(); basepath == "/" {
		basepath = "/app"
	}

	base, baseErr := ioutil.ReadFile(fmt.Sprintf("%s/properties/base.yaml", basepath))
	env, envErr := ioutil.ReadFile(fmt.Sprintf("%s/properties/%s.yaml", basepath, os.Getenv("env")))

	if baseErr == nil && envErr == nil {
		baseErr = yaml.Unmarshal(base, &variables)
		envErr = yaml.Unmarshal(env, &variables)
		_ = mapstructure.Decode(variables, &properties)
		return
	}

	log.Fatal("Error trying get files to read properties, please input a valid env.", baseErr, envErr)
}
