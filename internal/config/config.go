package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"database"`
	Jwt struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
	Azure struct {
		Active       bool   `yaml:"active"`
		TenantID     string `yaml:"tenantID"`
		ClientID     string `yaml:"clientID"`
		ClientSecret string `yaml:"clientSecret"`
		RedirectURL  string `yaml:"redirectURL"`
	} `yaml:azure`
}

func New() *Config {
	var cfg Config
	readFile(&cfg)
	return &cfg
}

func readFile(cfg *Config) {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
