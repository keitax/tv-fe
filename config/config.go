package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SiteTitle      string `yaml:"site-title"`
	DataSourceName string `yaml:"data-source-name"`
	TemplateDir    string `yaml:"template-dir"`
	StaticDir      string `yaml:"static-dir"`
	SiteFooter     string `yaml:"site-footer"`
	Locale         string `yaml:"locale"`
	BaseUrl        string `yaml:"base-url"`
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
