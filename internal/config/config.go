package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerPort  string `yaml:"server_port"`
	LogLevel    string `yaml:"log_level"`
	ApiUrl      string `yaml:"api_url"`
	AccessToken string `yaml:"access_token"`
	Model       string `yaml:"model"`
}

var Instance *Config

func Init() {
	config := &Config{}
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	// server port
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort != "" {
		config.ServerPort = serverPort
	}
	// log level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		config.LogLevel = logLevel
	}
	// access_token
	accessToken := os.Getenv("ACCESS_TOKEN")
	if accessToken != "" {
		config.AccessToken = accessToken
	}
	// model
	model := os.Getenv("MODEL")
	if model != "" {
		config.Model = model
	} else {
		config.Model = "text-davinci-002-render-sha"
	}
	Instance = config
}
