package config

import (
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
	ErrorLog struct {
		File  string `yaml:"file"`
		Level string `yaml:"level"`
	} `yaml:"errorLog"`
	Slack struct {
		Token            string `yaml:"token"`
		BroadcastChannel string `yaml:"broadcast_channel"`
		UserAdminEmail   string `yaml:"user_admin_email"`
	} `yaml:"slack"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := readFile(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func readFile(cfg *Config) error {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		return err
	}
	return nil
}
