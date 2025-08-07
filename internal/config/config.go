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
	Slack struct {
		Token            string `yaml:"token"`
		BroadcastChannel string `yaml:"broadcast_channel"`
		UserAdminEmail   string `yaml:"user_admin_email"`
	} `yaml:"slack"`
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
