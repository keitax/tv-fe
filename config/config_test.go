package config

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
data-source-name = "user:pass@tcp(address:port)/name"
template-dir = "./templates"
static-dir = "./static"
site-footer = "Copyright &copy; YOUR NAME"
locale = "UTC"
base-url = "http://localhost/"
run-level = "production"
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
		DataSourceName: "user:pass@tcp(address:port)/name",
		TemplateDir:    "./templates",
		StaticDir:      "./static",
		SiteFooter:     "Copyright &copy; YOUR NAME",
		Locale:         "UTC",
		BaseUrl:        "http://localhost/",
		RunLevel:       ProductionRunLevel,
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
