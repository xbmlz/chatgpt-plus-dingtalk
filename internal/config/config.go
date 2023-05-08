package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  int    `yaml:"SERVER_PORT" mapstructure:"SERVER_PORT"`
	LogLevel    string `yaml:"LOG_LEVEL" mapstructure:"LOG_LEVEL"`
	ApiUrl      string `yaml:"API_URL" mapstructure:"API_URL"`
	AccessToken string `yaml:"ACCESS_TOKEN" mapstructure:"ACCESS_TOKEN"`
	Model       string `yaml:"MODEL" mapstructure:"MODEL"`
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
