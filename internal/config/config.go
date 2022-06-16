package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var cfg *Config

func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

type Database struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Name       string `yaml:"name"`
	Migrations string `yaml:"migrations"`
}

type Smtp struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Metrics struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

type Jaeger struct {
	Service string `yaml:"service"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

type Kafka struct {
	ConsumerTopic string   `yaml:"consumerTopic"`
	GroupId       string   `yaml:"groupId"`
	Brokers       []string `yaml:"brokers"`
}

type Config struct {
	Database Database `yaml:"database"`
	Kafka    Kafka    `yaml:"kafka"`
	Metrics  Metrics  `yaml:"metrics"`
	Jaeger   Jaeger   `yaml:"jaeger"`
	Smtp     Smtp     `yaml:"smtp"`
	Debug    bool     `yaml:"debug"`
}

func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	return nil
}
