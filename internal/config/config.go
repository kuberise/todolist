package config

import (
	"os"

	"github.com/kuberise/todolist/internal/controller"
	"github.com/kuberise/todolist/pkg/postgres"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP     controller.HTTPConfig `yaml:"http"`
	Postgres postgres.Config       `yaml:"postgres"`
}

func New(path string) (*Config, error) {

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(b, &c)

	return &c, err
}
