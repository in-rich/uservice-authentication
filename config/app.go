package config

import (
	_ "embed"
	"github.com/in-rich/lib-go/deploy"
)

//go:embed app.yaml
var appFile []byte

type AppType struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Postgres struct {
		DSN string `yaml:"dsn"`
	} `yaml:"postgres"`
}

var App = deploy.LoadConfig[AppType](
	deploy.GlobalConfig(appFile),
)
