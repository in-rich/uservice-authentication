package config

import (
	_ "embed"
	"github.com/goccy/go-yaml"
	"os"
)

//go:embed app.dev.yaml
var appDevFile []byte

type AppType struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Postgres struct {
		DSN string `yaml:"dsn"`
	} `yaml:"postgres"`
}

var App AppType

func init() {
	switch os.Getenv("ENV") {
	case "prod":
		panic("not implemented")
	case "staging":
		panic("not implemented")
	default:
		if err := yaml.Unmarshal(appDevFile, &App); err != nil {
			panic(err)
		}
	}
}
