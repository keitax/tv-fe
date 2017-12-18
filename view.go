package tvfe

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// View represents the each blog page.
type View struct {
	urlBuilder   *URLBuilder
	config       *Config
	templateName string
	context      map[string]interface{}
}

// Render renders the view content.
func (v *View) Render(w io.Writer) {
	ctx := map[string]interface{}{
		"SiteTitle":  v.config.SiteTitle,
		"SiteFooter": v.config.SiteFooter,
		"Urls":       v.urlBuilder,
	}
	for key, value := range v.context {
		ctx[key] = value
	}
	if err := v.loadTemplates().ExecuteTemplate(w, v.templateName, ctx); err != nil {
		panic(err)
	}
}

func (v *View) loadTemplates() *template.Template {
	fs := []string{}
	if err := filepath.Walk(v.config.TemplateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".tmpl") {
			fs = append(fs, path)
		}
		return nil
	}); err != nil {
		panic(err)
	}
	ts, err := template.New("root").Funcs(template.FuncMap{
		"RenderMarkdown": ParseMarkdown,
		"ShowTime": func(t time.Time) string {
			return t.Format("Jan. 02, 2006, 3:04 PM")
		},
	}).ParseFiles(fs...)
	if err != nil {
		panic(err)
	}
	return ts
}
