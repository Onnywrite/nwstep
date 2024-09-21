package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port uint32 `env:"PORT"          env-default:"8080" yaml:"port"`
	Conn string `env:"POSTGRES_CONN" yaml:"conn"`
}

func MustLoad(additionalPath string) *Config {
	var path string

	flag.StringVar(&path, "config", additionalPath, "path to config file")
	flag.Parse()

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
