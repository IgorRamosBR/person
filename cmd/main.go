package main

import "person/configs"

// @title Person API
// @version 1.0
// @description This is a crud of people.
// @BasePath /v1
func main() {
	configs.Variables()
	configs.Logrus()
	configs.Di()
	configs.Routes()
}
