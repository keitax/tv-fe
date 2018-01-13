package tvfe

import (
	"os"
)

// Run-level flags.
const (
	ProductionRunLevel  = "production"
	DevelopmentRunLevel = "development"
)

// Config structure of Textvid application.
type Config struct {
	SiteTitle           string
	TemplateDir         string
	StaticDir           string
	SiteFooter          string
	Locale              string
	BaseURL             string
	RunLevel            string
	LocalGitRepository  string
	RemoteGitRepository string
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
