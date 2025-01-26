package configs

import (
	_ "embed"
	"log"

	appConfig "github.com/romakorinenko/task-manager/internal/config"
	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var cfg []byte

func MustLoadConfig() *appConfig.Config {
	appCfg := &appConfig.Config{}

	if err := yaml.Unmarshal(cfg, &appCfg); err != nil {
		log.Fatalln(err)
	}

	return appCfg
}
