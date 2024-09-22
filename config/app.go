package config

import (
	_ "embed"
	"github.com/in-rich/lib-go/deploy"
)

//go:embed app.dev.yaml
var appDevFile []byte

//go:embed app.staging.yaml
var appStagingFile []byte

//go:embed app.prod.yaml
var appProdFile []byte

type AppType struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Postgres struct {
		DSN string `yaml:"dsn"`
	} `yaml:"postgres"`
}

var App = deploy.LoadConfig[AppType](
	deploy.DevConfig(appDevFile),
	deploy.StagingConfig(appStagingFile),
	deploy.ProdConfig(appProdFile),
)
