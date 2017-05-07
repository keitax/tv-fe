package view

import (
	"html/template"
	"io"
	"path/filepath"
	"time"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/util"
)

type View struct {
	urlBuilder   *util.UrlBuilder
	config       *config.Config
	templateName string
	context      map[string]interface{}
}

func (v *View) Render(w io.Writer) error {
	ts := template.New("root").Funcs(template.FuncMap{
		"RenderMarkdown": util.ParseMarkdown,
		"ShowTime": func(t time.Time) string {
			return t.Format("Jan. 02, 2006, 3:04 PM")
		},
	})
	ts = template.Must(ts.ParseFiles(
		filepath.Join(v.config.TemplateDir, "layout.tmpl"),
		filepath.Join(v.config.TemplateDir, v.templateName),
	))
	context_ := map[string]interface{}{
		"SiteTitle":  v.config.SiteTitle,
		"SiteFooter": v.config.SiteFooter,
		"Urls":       v.urlBuilder,
	}
	for key, value := range v.context {
		context_[key] = value
	}
	if err := ts.ExecuteTemplate(w, "layout.tmpl", context_); err != nil {
		return err
	}
	return nil
}
