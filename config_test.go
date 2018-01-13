package tvfe

import (
	"os"
	"testing"
)

func TestGetFromEnv(t *testing.T) {
	os.Setenv("TV_SITE_TITLE", "textvid blog")
	os.Setenv("TV_SITE_TITLE", "textvid blog")
	os.Setenv("TV_TEMPLATE_DIR", "templates")
	os.Setenv("TV_STATIC_DIR", "static")
	os.Setenv("TV_SITE_FOOTER", "footer")
	os.Setenv("TV_LOCALE", "UTC")
	os.Setenv("TV_BASE_URL", "/")
	os.Setenv("TV_RUN_LEVEL", "production")
	os.Setenv("TV_LOCAL_GIT_REPOSITORY", "/tmp/textvid")
	os.Setenv("TV_REMOTE_GIT_REPOSITORY", "git@github.com:textvid/blog")
	c := GetFromEnv()

	tests := []struct {
		envName  string
		envValue string
		actual   string
	}{
		{"TV_SITE_TITLE", "textvid blog", c.SiteTitle},
		{"TV_TEMPLATE_DIR", "templates", c.TemplateDir},
		{"TV_STATIC_DIR", "static", c.StaticDir},
		{"TV_SITE_FOOTER", "footer", c.SiteFooter},
		{"TV_LOCALE", "UTC", c.Locale},
		{"TV_BASE_URL", "/", c.BaseURL},
		{"TV_RUN_LEVEL", "production", c.RunLevel},
		{"TV_LOCAL_GIT_REPOSITORY", "/tmp/textvid", c.LocalGitRepository},
		{"TV_REMOTE_GIT_REPOSITORY", "git@github.com:textvid/blog", c.RemoteGitRepository},
	}
	for _, test := range tests {
		if test.envValue != test.actual {
			t.Errorf("%s: expected %#v, actual %#v", test.envName, test.envValue, test.envName)
		}
	}
}
