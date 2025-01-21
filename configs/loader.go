package configs

import (
	_ "embed"
	"log"

	config2 "github.com/romakorinenko/task-manager/internal/config"
	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var cfg []byte

func MustLoadConfig() *config2.Config {
	appCfg := &config2.Config{}

	if err := yaml.Unmarshal(cfg, &appCfg); err != nil {
		log.Fatalln(err)
	}

	return appCfg
}
