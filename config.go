package textvid

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
	SiteTitle           string `toml:"site-title"`
	TemplateDir         string `toml:"template-dir"`
	StaticDir           string `toml:"static-dir"`
	SiteFooter          string `toml:"site-footer"`
	Locale              string `toml:"locale"`
	BaseUrl             string `toml:"base-url"`
	RunLevel            string `toml:"run-level"`
	LocalGitRepository  string `toml:"local-git-repository"`
	RemoteGitRepository string `toml:"remote-git-repository"`
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
