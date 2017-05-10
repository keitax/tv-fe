package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

const (
	ProductionRunLevel  = "production"
	DevelopmentRunLevel = "development"
)

type Config struct {
	SiteTitle      string `toml:"site-title"`
	DataSourceName string `toml:"data-source-name"`
	TemplateDir    string `toml:"template-dir"`
	StaticDir      string `toml:"static-dir"`
	SiteFooter     string `toml:"site-footer"`
	Locale         string `toml:"locale"`
	BaseUrl        string `toml:"base-url"`
	RunLevel       string `toml:"run-level"`
}

func Parse(configFile string) (*Config, error) {
	bs, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := toml.Unmarshal(bs, &c); err != nil {
		return nil, err
	}
	if !(c.RunLevel == ProductionRunLevel || c.RunLevel == DevelopmentRunLevel) {
		return nil, fmt.Errorf("run-level must be %#v or %#v", ProductionRunLevel, DevelopmentRunLevel)
	}
	return &c, nil
}
