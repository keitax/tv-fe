package textvid

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type View struct {
	urlBuilder   *UrlBuilder
	config       *Config
	templateName string
	context      map[string]interface{}
}

func (v *View) Render(w io.Writer) {
	context_ := map[string]interface{}{
		"SiteTitle":  v.config.SiteTitle,
		"SiteFooter": v.config.SiteFooter,
		"Urls":       v.urlBuilder,
	}
	for key, value := range v.context {
		context_[key] = value
	}
	if err := v.loadTemplates().ExecuteTemplate(w, v.templateName, context_); err != nil {
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
