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
		Uri string
		Database string
		Collection string
	}
	Port int
	Log struct {
		Level string
		JsonFormatter bool
	}
}

var properties Properties

func Variables() {
	variables := make(map[string]interface{})
	base, baseErr := ioutil.ReadFile("properties/base.yaml")
	env, envErr := ioutil.ReadFile(fmt.Sprintf("properties/%s.yaml", os.Getenv("env")))

	if baseErr == nil && envErr == nil {
		baseErr = yaml.Unmarshal(base, &variables)
		envErr = yaml.Unmarshal(env, &variables)
		_ = mapstructure.Decode(variables, &properties)
		return
	}

	log.Fatal("Error trying get files to read properties, please input a valid env.")
}