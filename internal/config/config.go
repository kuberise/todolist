package config

import (
	"github.com/kuberise/todolist/internal/controller"
	"github.com/kuberise/todolist/pkg/postgres"
)

type Config struct {
	HTTP     controller.HTTPConfig `cfg:"http"`
	Postgres postgres.Config       `cfg:"postgres"`
}

func Default() Config {
	return Config{}
}
