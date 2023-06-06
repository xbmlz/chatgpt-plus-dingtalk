package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	ServerUrl             string `yaml:"SERVER_URL" mapstructure:"SERVER_URL"`
	ServerPort            int    `yaml:"SERVER_PORT" mapstructure:"SERVER_PORT"`
	LogLevel              string `yaml:"LOG_LEVEL" mapstructure:"LOG_LEVEL"`
	ChatgptBaseUrl        string `yaml:"CHATGPT_BASE_URL" mapstructure:"CHATGPT_BASE_URL"`
	ChatgptAccessToken    string `yaml:"CHATGPT_ACCESS_TOKEN" mapstructure:"CHATGPT_ACCESS_TOKEN"`
	ChatgptModel          string `yaml:"CHATGPT_MODEL" mapstructure:"CHATGPT_MODEL"`
	ReplicateBaseUrl      string `yaml:"REPLICATE_BASE_URL" mapstructure:"REPLICATE_BASE_URL"`
	ReplicateApiToken     string `yaml:"REPLICATE_API_TOKEN" mapstructure:"REPLICATE_API_TOKEN"`
	ReplicateModelVersion string `yaml:"REPLICATE_MODEL_VERSION" mapstructure:"REPLICATE_MODEL_VERSION"`
	CleanAllSessionCron   string `yaml:"CLEAN_ALL_SESSION_CRON" mapstructure:"CLEAN_ALL_SESSION_CRON"`
}

var Instance Config

func Initialize() {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")
	// v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err = v.Unmarshal(&Instance); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&Instance); err != nil {
		fmt.Println(err)
	}
}
