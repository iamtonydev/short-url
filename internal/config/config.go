package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Grpc struct {
	Port string `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"database"`
	SslMode  string `yaml:"sslmode"`
	Driver   string `yaml:"driver"`
}

type Http struct {
	Port string `yaml:"port"`
}

type Config struct {
	Grpc     Grpc     `yaml:"grpc"`
	Database Database `yaml:"database"`
	Http     Http     `yaml:"http"`
}

func Read(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
