package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Db struct {
		Name string `yaml:"name"`
	} `yaml:"db"`
	Server struct {
		Port    string `yaml:"port"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"server"`
	Worker struct {
		Interval int `yaml:"interval"`
	} `yaml:"worker"`
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ParseCLI() (string, error) {
	var path string

	flag.StringVar(&path, "config", "./config.yaml", "path to config file")
	flag.Parse()

	s, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if s.IsDir() {
		return "", err
	}

	return path, nil
}
