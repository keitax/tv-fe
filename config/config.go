package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SiteTitle   string `yaml:"site-title"`
	DatabaseDir string `yaml:"database-dir"`
	TemplateDir string `yaml:"template-dir"`
}

func Parse(configFile string) (*Config, error) {
	bs, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(bs, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
