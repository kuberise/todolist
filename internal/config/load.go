package config

import (
	"log"
	"log/slog"
	"strings"
	"sync"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const (
	prefix    = "TODOLIST_"
	delimiter = "."
	separator = "__"
	structTag = "cfg" // the tag on the config struct that we use to map config path to the struct fields.
)

var cfgOnce sync.Once // Singleton config instance.
var Cfg *Config

// TODOLIST_DEBUG -> DEBUG -> debug
// TODOLIST_DATABASE__HOST -> DATABASE__HOST -> database__host -> database.host

func envToPathConverter(source string) string {
	base := strings.ToLower(strings.TrimPrefix(source, prefix))
	return strings.ReplaceAll(base, separator, delimiter)
}

func New(l *slog.Logger) *Config {
	cfgOnce.Do(func() {

		k := koanf.New(".")

		// load default configuration from default function
		if err := k.Load(structs.Provider(Default(), structTag), nil); err != nil {
			log.Fatalf("error loading default: %v", err)
		}

		// load yaml file
		if err := k.Load(file.Provider("/configs.yaml"), yaml.Parser()); err != nil {
			l.Warn("error loading yaml config file: %v", err)
		}

		// load environment variables
		if err := k.Load(env.Provider(prefix, delimiter, envToPathConverter), nil); err != nil {
			log.Printf("error loading environment variables: %v", err)
		}

		Cfg = &Config{}
		if err := k.UnmarshalWithConf("", Cfg, koanf.UnmarshalConf{Tag: structTag}); err != nil {
			log.Fatalf("error unmarshalling config: %s", err)
		}
	})
	return Cfg
}
