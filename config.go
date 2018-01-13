package tvfe

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// Run-level flags.
const (
	ProductionRunLevel  = "production"
	DevelopmentRunLevel = "development"
)

// Config structure of Textvid application.
type Config struct {
	SiteTitle           string `toml:"site-title"`
	TemplateDir         string `toml:"template-dir"`
	StaticDir           string `toml:"static-dir"`
	SiteFooter          string `toml:"site-footer"`
	Locale              string `toml:"locale"`
	BaseURL             string `toml:"base-url"`
	RunLevel            string `toml:"run-level"`
	LocalGitRepository  string `toml:"local-git-repository"`
	RemoteGitRepository string `toml:"remote-git-repository"`
}

// Parse reads a config file.
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

// GetFromEnv gets a config from environment vars.
func GetFromEnv() *Config {
	return &Config{
		SiteTitle:           os.Getenv("TV_SITE_TITLE"),
		TemplateDir:         os.Getenv("TV_TEMPLATE_DIR"),
		StaticDir:           os.Getenv("TV_STATIC_DIR"),
		SiteFooter:          os.Getenv("TV_SITE_FOOTER"),
		Locale:              os.Getenv("TV_LOCALE"),
		BaseURL:             os.Getenv("TV_BASE_URL"),
		RunLevel:            os.Getenv("TV_RUN_LEVEL"),
		LocalGitRepository:  os.Getenv("TV_LOCAL_GIT_REPOSITORY"),
		RemoteGitRepository: os.Getenv("TV_REMOTE_GIT_REPOSITORY"),
	}
}
