package tvfe

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "config-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(`
site-title = "Textvid Blog"
template-dir = "./templates"
static-dir = "./static"
site-footer = "Copyright &copy; YOUR NAME"
locale = "UTC"
base-url = "http://localhost/"
run-level = "production"
local-git-repository = "/tmp/textvid-blog"
remote-git-repository = "https://host.git/textvid-blog.git"
`); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	parsed, err := Parse(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SiteTitle:           "Textvid Blog",
		TemplateDir:         "./templates",
		StaticDir:           "./static",
		SiteFooter:          "Copyright &copy; YOUR NAME",
		Locale:              "UTC",
		BaseURL:             "http://localhost/",
		RunLevel:            ProductionRunLevel,
		LocalGitRepository:  "/tmp/textvid-blog",
		RemoteGitRepository: "https://host.git/textvid-blog.git",
	}
	if !reflect.DeepEqual(expected, parsed) {
		t.Errorf("Failed to parse: expected: %v, parsed: %v", expected, parsed)
	}
}

func TestParseToFailWhenInvalidRunLevel(t *testing.T) {
	testCases := []struct {
		descr       string
		expectedErr error
		runLevel    string
	}{
		{"valid runlevel", nil, "production"},
		{"valid runlevel", nil, "development"},
		{"invalid runlevel", errors.New(`run-level must be "production" or "development"`), "foobar"},
	}

	for _, tc := range testCases {
		f, err := ioutil.TempFile(os.TempDir(), "config-test-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		if _, err := f.WriteString(fmt.Sprintf(`run-level = "%s"`, tc.runLevel)); err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}

		_, err = Parse(f.Name())
		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Errorf("%s: Parse() = %#v, expected %#v", tc.descr, err, tc.expectedErr)
		}
	}
}

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
