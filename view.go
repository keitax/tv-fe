package textvid

import (
	"html/template"
	"io"
	"path/filepath"
	"time"
)

type View struct {
	urlBuilder   *UrlBuilder
	config       *Config
	templateName string
	context      map[string]interface{}
}

func (v *View) Render(w io.Writer) {
	ts := template.New("root").Funcs(template.FuncMap{
		"RenderMarkdown": ParseMarkdown,
		"ShowTime": func(t time.Time) string {
			return t.Format("Jan. 02, 2006, 3:04 PM")
		},
	})
	ts = template.Must(ts.ParseFiles(
		filepath.Join(v.config.TemplateDir, v.templateName),
		filepath.Join(v.config.TemplateDir, "_header.tmpl"),
		filepath.Join(v.config.TemplateDir, "_footer.tmpl"),
	))
	context_ := map[string]interface{}{
		"SiteTitle":  v.config.SiteTitle,
		"SiteFooter": v.config.SiteFooter,
		"Urls":       v.urlBuilder,
	}
	for key, value := range v.context {
		context_[key] = value
	}
	if err := ts.ExecuteTemplate(w, v.templateName, context_); err != nil {
		panic(err)
	}
}
