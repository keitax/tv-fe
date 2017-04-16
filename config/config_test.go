package config

import (
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
site-title: Textvid Blog
data-source-name: user/pass@tcp(address:port)/name
template-dir: ./templates
static-dir: ./static
site-footer: Copyright &copy; YOUR NAME
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
		SiteTitle:      "Textvid Blog",
		DataSourceName: "user/pass@tcp(address:port)/name",
		TemplateDir:    "./templates",
		StaticDir:      "./static",
		SiteFooter:     "Copyright &copy; YOUR NAME",
	}
	if !reflect.DeepEqual(expected, parsed) {
		t.Errorf("Failed to parse: expected: %v, parsed: %v", expected, parsed)
	}
}
